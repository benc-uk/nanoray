package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "nanoray/pkg/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedControllerServer
}

var (
	portFlag = flag.String("port", "5000", "The port to listen on")
)

func main() {
	portNum := os.Getenv("PORT")
	flag.Parse()

	if portNum == "" {
		portNum = *portFlag
	}
	port, _ := strconv.Atoi(portNum)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	workers = []WorkerConnection{}

	s := grpc.NewServer()
	pb.RegisterControllerServer(s, &server{})

	log.Printf("Server started on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
