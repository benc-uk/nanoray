package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	pb "raynet/pkg/proto"

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
	Conn   *grpc.ClientConn
}

var workers sync.Map
var workerCount int

func init() {
	workers = sync.Map{}
}

func (s *server) RegisterWorker(ctx context.Context, worker *pb.WorkerInfo) (*emptypb.Empty, error) {
	log.Printf("Worker %s connecting from %s:%d", worker.Id, worker.GetHost(), worker.GetPort())

	// Connect back to the worker
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", worker.GetHost(), worker.GetPort()),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	newClient := pb.NewWorkerClient(conn)

	// A handshake ping to ensure the worker is connectable via the hostname provided
	_, err = newClient.Ping(context.Background(), &emptypb.Empty{})
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
	workerCount++

	log.Printf("Workers online: %d", workerCount)

	return &emptypb.Empty{}, nil
}

func (s *server) GetWorkers(ctx context.Context, in *emptypb.Empty) (*pb.WorkerList, error) {
	workerList := make([]*pb.WorkerInfo, 0)

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
