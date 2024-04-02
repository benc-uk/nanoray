package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "raynet/shared/proto"

	"google.golang.org/grpc"
)

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
		log.Fatalf("Failed to bind to port\n%s", err.Error())
	}

	grpcSrv := grpc.NewServer()

	pb.RegisterControllerServer(grpcSrv, &server{})

	log.Printf("Controller started on port %d", port)
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve\n%s", err.Error())
	}
}
