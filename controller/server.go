package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "nanoray/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerConnection struct {
	Worker *pb.WorkerInfo
	Client pb.WorkerClient
}

var workers []WorkerConnection

func (s *server) RegisterWorker(ctx context.Context, worker *pb.WorkerInfo) (*pb.RegistrationResult, error) {
	log.Printf("Received new worker: %s %d", worker.GetAddress(), worker.GetPort())

	// Connect to the worker
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", worker.GetAddress(), worker.GetPort()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		log.Printf("Connected to worker %s:%d", worker.GetAddress(), worker.GetPort())
	}

	newClient := pb.NewWorkerClient(conn)
	workers = append(workers, WorkerConnection{
		Worker: worker,
		Client: newClient,
	})

	go func() {
		time.Sleep(1 * time.Second)
		log.Printf("Sending job to worker")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		jobRes, err := newClient.NewJob(ctx, &pb.JobRequest{
			Id: 1,
		})

		if err != nil {
			log.Fatalf("failed to send worker job: %v", err)
		}

		log.Printf("Job result: %v+", jobRes)
	}()

	return &pb.RegistrationResult{
		StatusCode: 0,
		Message:    "Success",
	}, nil

}

func (s *server) GetScene(ctx context.Context, in *pb.SceneRequest) (*pb.SceneResult, error) {
	log.Printf("Returning scene")

	return &pb.SceneResult{}, nil
}
