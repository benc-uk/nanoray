package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sort"
	"sync"

	pb "raynet/shared/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedControllerServer
}

type WorkerConnection struct {
	Worker *pb.WorkerInfo
	Client pb.WorkerClient
	Conn   *grpc.ClientConn
}

type Render struct {
	SceneId string
	Status  int
}

var workers sync.Map
var workerCount = 0
var workerCountLock sync.Mutex

var outputImage *image.RGBA
var pendingJobs sync.Map

var jobsTotal = 0
var jobsOutstanding = 0
var jobLock sync.Mutex

var jobIdCounter int32 = 0

func init() {
	workers = sync.Map{}
}

func (s *server) RegisterWorker(ctx context.Context, worker *pb.WorkerInfo) (*pb.Void, error) {
	log.Printf("Worker %s connecting from %s:%d", worker.Id, worker.GetHost(), worker.GetPort())

	// Connect back to the worker
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", worker.GetHost(), worker.GetPort()),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	newClient := pb.NewWorkerClient(conn)

	// A handshake ping to ensure the worker is connectable via the hostname provided
	_, err = newClient.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("Unable to RPC ping worker, handshake failed\n%s", err.Error())
		return nil, err
	}

	go func() {
		// This traps the worker connection state changes
		for {
			change := conn.WaitForStateChange(context.Background(), conn.GetState())
			if !change {
				return
			}

			currentState := conn.GetState()
			if currentState == connectivity.Idle {
				log.Printf("Worker %s disconnected", worker.Id)
				w, ok := workers.Load(worker.Id)
				if !ok {
					return
				}
				w.(WorkerConnection).Conn.Close()
				workers.Delete(worker.Id)

				workerCountLock.Lock()
				defer workerCountLock.Unlock()
				workerCount--

				log.Printf("Workers online: %d", workerCount)
				return
			}
		}
	}()

	workers.Store(worker.Id, WorkerConnection{
		Worker: worker,
		Client: newClient,
		Conn:   conn,
	})

	workerCountLock.Lock()
	workerCount++
	workerCountLock.Unlock()

	log.Printf("Workers online: %d", workerCount)

	return nil, nil
}

func (s *server) GetWorkers(ctx context.Context, in *pb.Void) (*pb.WorkerList, error) {
	workerList := []*pb.WorkerInfo{}

	workers.Range(func(key, value interface{}) bool {
		workerList = append(workerList, value.(WorkerConnection).Worker)
		return true
	})

	sort.Slice(workerList, func(i, j int) bool {
		return workerList[i].Id < workerList[j].Id
	})

	return &pb.WorkerList{
		Workers: workerList,
	}, nil
}

func (s *server) GetScene(ctx context.Context, in *pb.SceneRequest) (*pb.SceneResult, error) {
	log.Printf("Returning scene")

	return &pb.SceneResult{}, nil
}

func (s *server) StartRender(ctx context.Context, in *pb.Void) (*pb.Void, error) {
	if workerCount == 0 {
		log.Printf("No workers available to start render")
		return nil, status.Errorf(codes.FailedPrecondition, "No workers available to start render")
	}

	imgW := 1920
	imgH := 1080
	jobW := 32
	jobH := 32
	outputImage = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{imgW, imgH}})

	pendingJobs = sync.Map{}
	jobCount := 0
	for y := 0; y < imgH; y += jobH {
		for x := 0; x < imgW; x += jobW {
			pendingJobs.Store(jobIdCounter, &pb.JobRequest{
				Id:      jobIdCounter,
				SceneId: "test",
				Width:   int32(jobW),
				Height:  int32(jobH),
				X:       int32(x),
				Y:       int32(y),
			})

			jobIdCounter++
			jobCount++
		}
	}

	log.Printf("Starting render with %d jobs", jobCount)

	jobsOutstanding = jobCount
	jobsTotal = jobCount

	workers.Range(func(_, worker interface{}) bool {
		w := worker.(WorkerConnection)
		maxJobs := w.Worker.MaxJobs
		pendingJobs.Range(func(jobId, job interface{}) bool {
			if maxJobs == 0 {
				log.Printf("Initial %d jobs dispatched to worker %s", w.Worker.MaxJobs, w.Worker.Id)
				return false
			}

			jobReq := job.(*pb.JobRequest)
			pendingJobs.Delete(jobId)

			_, err := w.Client.NewJob(context.Background(), jobReq)
			if err != nil {
				log.Printf("Failed to send job %d to worker %s\n%s", jobId, w.Worker.Id, err.Error())
				return false
			}

			maxJobs--
			return true
		})

		return true
	})

	return nil, nil
}

func (s *server) JobComplete(ctx context.Context, result *pb.JobResult) (*pb.Void, error) {
	jobLock.Lock()
	defer jobLock.Unlock()

	job := result.Job

	srcImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(job.Width), int(job.Height)}})
	srcImg.Pix = result.ImageData

	// Update the render image with the job result
	draw.Draw(outputImage, image.Rect(int(job.X), int(job.Y), int(job.X+job.Width), int(job.Y+job.Height)), srcImg, image.Point{0, 0}, draw.Src)

	jobsOutstanding--

	//log.Printf("Job %d complete, %d jobs remaining", job.Id, jobsToComplete)

	if jobsOutstanding == 0 {
		log.Printf("All jobs completed, saving file!!!")

		f, err := os.Create("render.png")
		if err != nil {
			log.Printf("Failed to create render file\n%s", err.Error())
			return nil, err
		}
		defer f.Close()

		err = png.Encode(f, outputImage)
		if err != nil {
			log.Printf("Failed to encode render image\n%s", err.Error())
			return nil, err
		}

	} else if jobsOutstanding > 0 {
		var nextJob *pb.JobRequest = nil
		pendingJobs.Range(func(_, job interface{}) bool {
			nextJob = job.(*pb.JobRequest)
			return nextJob != nil
		})

		if nextJob != nil {
			//log.Printf("Dispatching job: %d to worker %s", nextJob.Id, result.Worker.Id)
			workerConn, _ := workers.Load(result.Worker.Id)
			worker := workerConn.(WorkerConnection)

			pendingJobs.Delete(nextJob.Id)
			worker.Client.NewJob(context.Background(), nextJob)
		}
	}

	return nil, nil
}

func (s *server) GetProgress(ctx context.Context, in *pb.Void) (*pb.Progress, error) {
	jobLock.Lock()
	defer jobLock.Unlock()

	return &pb.Progress{
		TotalJobs:     int32(jobsTotal),
		CompletedJobs: int32(jobsTotal - jobsOutstanding),
	}, nil
}
