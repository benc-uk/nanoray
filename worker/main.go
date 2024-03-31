package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	pb "nanoray/pkg/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	portFlag     = flag.String("port", "0", "The offset of the port to listen on")
	portBaseFlag = flag.String("portBase", "4400", "The port number base to listen on")
	port         int32
	ctrlConn     *grpc.ClientConn
	ctrlClient   pb.ControllerClient
	workerInfo   pb.WorkerInfo
)

const controllerAddr = "localhost:5000"

func main() {
	// Lots of guff for port handling, env vars and flags
	// We do things a little differently with port offsets, to make it easier to run multiple workers on the same machine
	portNum := os.Getenv("PORT")
	portBase := os.Getenv("PORT_BASE")
	flag.Parse()

	if portNum == "" {
		portNum = *portFlag
	}
	if portBase == "" {
		portBase = *portBaseFlag
	}
	portNumInt, _ := strconv.Atoi(portNum)
	portBaseInt, _ := strconv.Atoi(portBase)

	// Final port number
	port = int32(portBaseInt + portNumInt)

	workerInfo = pb.WorkerInfo{
		Id:      strings.Split(uuid.NewString(), "-")[0],
		Address: "localhost",
		Port:    port,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to bind to port\n%s", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &server{})

	waitForController()
	defer ctrlConn.Close()

	res, err := ctrlClient.RegisterWorker(context.Background(), &workerInfo)
	if err != nil {
		log.Fatalf("Registration failed\n%s", err.Error())
	}

	log.Printf("Registered with controller: %s", res.GetMessage())

	log.Printf("Worker started on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve\n%s", err.Error())
	}
}

func waitForController() {
	log.Printf("Connecting to controller at %s", controllerAddr)

	for {
		var err error
		ctrlConn, err = grpc.Dial(controllerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Unable to dial controller\n%s", err.Error())
		}

		ctrlClient = pb.NewControllerClient(ctrlConn)

		ctrlConn.WaitForStateChange(context.Background(), ctrlConn.GetState())
		if ctrlConn.GetState() == connectivity.Ready {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
