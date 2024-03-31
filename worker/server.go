package main

import (
	"context"
	"log"
	"time"

	pb "nanoray/pkg/proto"

	"google.golang.org/protobuf/types/known/durationpb"
)

type server struct {
	pb.UnimplementedWorkerServer
}

func (s *server) NewJob(ctx context.Context, in *pb.JobRequest) (*pb.JobResult, error) {
	log.Printf("Received job: %d", in.GetId())

	now := time.Now()
	time.Sleep(5 * time.Second)

	return &pb.JobResult{
		StatusCode: 0,
		Message:    "Success",
		TimeTaken:  durationpb.New(time.Since(now)),
	}, nil
}
