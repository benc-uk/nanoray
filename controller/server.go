package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"sync"
	"time"

	pb "nanoray/lib/proto"
	rt "nanoray/lib/raytrace"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type server struct {
	pb.UnimplementedControllerServer
}

var netRender *rt.NetworkRender
var img *image.RGBA

var jobIdCounter int32 = 0

func (s *server) StartRender(ctx context.Context, in *pb.RenderRequest) (*pb.Void, error) {
	render := rt.NewRender(int(in.Width), in.AspectRatio)
	render.SamplesPerPixel = int(in.SamplesPerPixel)
	render.MaxDepth = int(in.MaxDepth)

	// Try to parse the scene data, we don't need the scene or camera, just need to know if it's valid
	_, _, err := rt.ParseScene(in.SceneData, render.Width, render.Height)
	if err != nil {
		log.Printf("Failed to parse scene data\n%s", err.Error())
		return nil, status.Errorf(codes.Aborted, "Failed to parse scene data: %s", err.Error())
	}

	if workerCount == 0 {
		log.Printf("No workers available to start render")
		return nil, status.Errorf(codes.FailedPrecondition, "No workers available to start render")
	}

	netRender = &rt.NetworkRender{
		Status:       rt.READY,
		JobQueue:     sync.Map{},
		JobsTotal:    0,
		JobsComplete: 0,
		Start:        time.Now(),
		OutputName:   time.Now().Format("2006-01-02_15:04:05"),
	}

	slices := 12
	jobW := int(render.Width)
	jobH := int(render.Height) / slices

	img = render.MakeImage()

	totalJobs := 0
	for y := 0; y < render.Height; y += jobH {
		for x := 0; x < render.Width; x += jobW {
			netRender.JobQueue.Store(jobIdCounter, &pb.JobRequest{
				Id:              jobIdCounter,
				Width:           int32(jobW),
				Height:          int32(jobH),
				X:               int32(x),
				Y:               int32(y),
				SamplesPerPixel: int32(render.SamplesPerPixel),
				MaxDepth:        int32(render.MaxDepth),
				ImageDetails:    render.ImageDetails(),
			})

			jobIdCounter++
			totalJobs++
		}
	}

	log.Printf("Starting render with %d jobs", totalJobs)

	netRender.JobsTotal = totalJobs

	workers.Range(func(_, worker interface{}) bool {
		w := worker.(WorkerConnection)
		maxJobs := w.Worker.MaxJobs

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		w.Client.PrepareRender(timeoutCtx, &pb.PrepRenderRequest{
			SceneData:    in.SceneData,
			ImageDetails: render.ImageDetails(),
		})

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
	draw.Draw(img, image.Rect(int(job.X), int(job.Y), int(job.X+job.Width), int(job.Y+job.Height)), srcImg, image.Point{0, 0}, draw.Src)

	netRender.JobsComplete++

	log.Printf("Job %d complete, %d jobs remaining", job.Id, netRender.JobsTotal-netRender.JobsComplete)

	if netRender.JobsComplete == netRender.JobsTotal {
		log.Printf("Time to complete: %s", time.Since(netRender.Start))
		log.Printf("All jobs completed, saving file!!!")

		os.Mkdir("output", os.ModePerm)
		f, err := os.Create(fmt.Sprintf("output/%s.png", netRender.OutputName))
		if err != nil {
			log.Printf("Failed to create render file\n%s", err.Error())
			return nil, err
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			log.Printf("Failed to encode render image\n%s", err.Error())
			return nil, err
		}

		netRender.Status = rt.COMPLETE
	} else {
		var nextJob *pb.JobRequest = nil
		netRender.JobQueue.Range(func(_, job interface{}) bool {
			nextJob = job.(*pb.JobRequest)
			return nextJob != nil
		})

		if nextJob != nil {
			log.Printf("Dispatching job: %d to worker %s", nextJob.Id, result.Worker.Id)
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
			OutputName:    netRender.OutputName,
		}, nil
	}

	netRender.Lock.Lock()
	defer netRender.Lock.Unlock()

	return &pb.Progress{
		TotalJobs:     int32(netRender.JobsTotal),
		CompletedJobs: int32(netRender.JobsComplete),
		OutputName:    netRender.OutputName,
	}, nil
}

func (s *server) ListRenderedImages(ctx context.Context, in *pb.Void) (*pb.ImageList, error) {
	files, err := os.ReadDir("output")
	if err != nil {
		return nil, err
	}

	var images []string
	for _, file := range files {
		images = append(images, file.Name())
	}
	for i, j := 0, len(images)-1; i < j; i, j = i+1, j-1 {
		images[i], images[j] = images[j], images[i]
	}

	return &pb.ImageList{
		Images: images,
	}, nil
}

func (s *server) GetRenderedImage(ctx context.Context, in *wrapperspb.StringValue) (*wrapperspb.BytesValue, error) {
	if in.Value == "" {
		return nil, status.Errorf(codes.InvalidArgument, "No image name provided")
	}

	if in.Value == "latest" {
		files, err := os.ReadDir("output")
		if err != nil {
			return nil, err
		}

		if len(files) == 0 {
			return nil, status.Errorf(codes.NotFound, "No images found")
		}

		in.Value = files[len(files)-1].Name()
	}

	if _, err := os.Stat(fmt.Sprintf("output/%s", in.Value)); os.IsNotExist(err) {
		return nil, status.Errorf(codes.NotFound, "Image not found")
	}

	f, err := os.Open(fmt.Sprintf("output/%s", in.Value))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rawBytes, err := os.ReadFile(fmt.Sprintf("output/%s", in.Value))
	if err != nil {
		return nil, err
	}

	return &wrapperspb.BytesValue{Value: rawBytes}, nil
}
