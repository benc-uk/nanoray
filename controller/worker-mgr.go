package main

import (
	"context"
	"fmt"
	"log"
	pb "nanoray/shared/proto"
	"sort"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerConnection struct {
	Worker *pb.WorkerInfo
	Client pb.WorkerClient
	Conn   *grpc.ClientConn
}

var workers sync.Map
var workerCount = 0
var workerCountLock sync.Mutex

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
		addressI := fmt.Sprintf("%s:%d", workerList[i].Host, workerList[i].Port)
		addressJ := fmt.Sprintf("%s:%d", workerList[j].Host, workerList[j].Port)
		return addressI < addressJ
	})

	return &pb.WorkerList{
		Workers: workerList,
	}, nil
}
