package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	pb "nanoray/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedControllerServer
}

type WorkerConnection struct {
	Worker *pb.WorkerInfo
	Client pb.WorkerClient
}

var workers map[string]WorkerConnection

func (s *server) RegisterWorker(ctx context.Context, worker *pb.WorkerInfo) (*pb.RegistrationResult, error) {
	log.Printf("Worker: %s connecting from %s:%d", worker.Id, worker.GetAddress(), worker.GetPort())

	// Connect back to the worker
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", worker.GetAddress(), worker.GetPort()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &pb.RegistrationResult{
			StatusCode: pb.StatusCode_ERROR,
			Message:    "Failed to connect to worker: " + err.Error(),
		}, err
	}

	newClient := pb.NewWorkerClient(conn)
	workers[worker.Id] = WorkerConnection{
		Worker: worker,
		Client: newClient,
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
				delete(workers, worker.Id)
				log.Printf("Workers online: %d", len(workers))
				return
			}
		}
	}()

	log.Printf("Workers online: %d", len(workers))

	return &pb.RegistrationResult{
		StatusCode: pb.StatusCode_OK,
		Message:    "Success",
	}, nil
}

func (s *server) GetWorkers(ctx context.Context, in *emptypb.Empty) (*pb.WorkerList, error) {
	log.Printf("Returning workers")

	workerList := make([]*pb.WorkerInfo, 0, len(workers))
	for _, worker := range workers {
		workerList = append(workerList, worker.Worker)
	}

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
