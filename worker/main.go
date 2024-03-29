package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "nanoray/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	pb.UnimplementedWorkerServer
}

var (
	portFlag     = flag.String("port", "0", "The offset of the port to listen on")
	portBaseFlag = flag.String("portBase", "4400", "The port number base to listen on")
	port         int32
	ctrlClient   pb.ControllerClient
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

	// client setup
	conn, err := grpc.Dial(controllerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ctrlClient = pb.NewControllerClient(conn)

	// Final port number
	port = int32(portBaseInt + portNumInt)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &server{})

	res, err := ctrlClient.RegisterWorker(context.Background(), &pb.WorkerInfo{
		Address: "localhost",
		Port:    port,
	})
	if err != nil {
		log.Fatalf("failed to register worker: %v", err)
	}

	log.Printf("Registered worker: %s", res.GetMessage())

	log.Printf("Server started on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
