package main

import (
	"context"
	"image"
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"nanoray/shared/controller"
	pb "nanoray/shared/proto"

	"google.golang.org/protobuf/types/known/durationpb"
)

type server struct {
	pb.UnimplementedWorkerServer
}

func (s *server) NewJob(ctx context.Context, job *pb.JobRequest) (*pb.Void, error) {
	log.Printf("Received job: %d", job.Id)

	go runJob(job)

	return nil, nil
}

func (s *server) Ping(ctx context.Context, in *pb.Void) (*pb.Void, error) {
	return nil, nil
}

func runJob(job *pb.JobRequest) error {
	now := time.Now()

	// Simulate some work calculate prime numbers
	for i := 2; i < int(10000); i++ {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				break
			}
		}
	}

	imgW := int(job.Width)
	imgH := int(job.Height)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{imgW, imgH}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	val := uint8(rand.IntN(128)) + 128
	randColour := color.RGBA{val, 0, 0, 255}

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.Set(x, y, randColour)
		}
	}

	res := &pb.JobResult{
		TimeTaken: durationpb.New(time.Since(now)),
		ImageData: img.Pix,
		Worker:    &workerInfo,
		Job:       job,
	}

	_, err := controller.Client.JobComplete(context.Background(), res)
	if err != nil {
		log.Printf("Failed to send completed job result: %s", err.Error())
		return err
	}

	return nil
}
