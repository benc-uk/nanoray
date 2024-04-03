package main

import (
	"context"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sync"

	pb "nanoray/shared/proto"
	"nanoray/shared/raytrace"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedControllerServer
}

var netRender *raytrace.NetworkRender

var jobIdCounter int32 = 0

func (s *server) GetScene(ctx context.Context, in *pb.SceneRequest) (*pb.SceneResult, error) {
	log.Printf("Returning scene")

	return &pb.SceneResult{}, nil
}

func (s *server) StartRender(ctx context.Context, in *pb.Void) (*pb.Void, error) {
	if workerCount == 0 {
		log.Printf("No workers available to start render")
		return nil, status.Errorf(codes.FailedPrecondition, "No workers available to start render")
	}

	// Create a new render object
	netRender = &raytrace.NetworkRender{
		Render: raytrace.Render{
			Image:  nil,
			Width:  1920,
			Height: 1080,
		},

		Status:       raytrace.READY,
		JobQueue:     sync.Map{},
		JobsTotal:    0,
		JobsComplete: 0,
	}

	imgW := 1920
	imgH := 1080
	jobW := 32
	jobH := 32
	netRender.Image = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{imgW, imgH}})

	jobCount := 0
	for y := 0; y < imgH; y += jobH {
		for x := 0; x < imgW; x += jobW {
			netRender.JobQueue.Store(jobIdCounter, &pb.JobRequest{
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

	netRender.JobsTotal = jobCount

	workers.Range(func(_, worker interface{}) bool {
		w := worker.(WorkerConnection)
		maxJobs := w.Worker.MaxJobs
		netRender.JobQueue.Range(func(jobId, job interface{}) bool {
			if maxJobs == 0 {
				log.Printf("Initial %d jobs dispatched to worker %s", w.Worker.MaxJobs, w.Worker.Id)
				return false
			}

			jobReq := job.(*pb.JobRequest)
			netRender.JobQueue.Delete(jobId)

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
	netRender.Lock.Lock()
	defer netRender.Lock.Unlock()

	job := result.Job

	srcImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(job.Width), int(job.Height)}})
	srcImg.Pix = result.ImageData

	// Update the render image with the job result
	draw.Draw(netRender.Image, image.Rect(int(job.X), int(job.Y), int(job.X+job.Width), int(job.Y+job.Height)), srcImg, image.Point{0, 0}, draw.Src)

	netRender.JobsComplete++

	//log.Printf("Job %d complete, %d jobs remaining", job.Id, jobsToComplete)

	if netRender.JobsComplete == netRender.JobsTotal {
		log.Printf("All jobs completed, saving file!!!")

		f, err := os.Create("render.png")
		if err != nil {
			log.Printf("Failed to create render file\n%s", err.Error())
			return nil, err
		}
		defer f.Close()

		err = png.Encode(f, netRender.Image)
		if err != nil {
			log.Printf("Failed to encode render image\n%s", err.Error())
			return nil, err
		}

		netRender.Status = raytrace.COMPLETE
	} else {
		var nextJob *pb.JobRequest = nil
		netRender.JobQueue.Range(func(_, job interface{}) bool {
			nextJob = job.(*pb.JobRequest)
			return nextJob != nil
		})

		if nextJob != nil {
			//log.Printf("Dispatching job: %d to worker %s", nextJob.Id, result.Worker.Id)
			workerConn, _ := workers.Load(result.Worker.Id)
			worker := workerConn.(WorkerConnection)

			netRender.JobQueue.Delete(nextJob.Id)
			worker.Client.NewJob(context.Background(), nextJob)
		}
	}

	return nil, nil
}

func (s *server) GetProgress(ctx context.Context, in *pb.Void) (*pb.Progress, error) {
	if netRender == nil {
		return &pb.Progress{
			TotalJobs:     0,
			CompletedJobs: 0,
		}, nil
	}

	netRender.Lock.Lock()
	defer netRender.Lock.Unlock()

	return &pb.Progress{
		TotalJobs:     int32(netRender.JobsTotal),
		CompletedJobs: int32(netRender.JobsComplete),
	}, nil
}
