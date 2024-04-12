package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "nanoray/lib/proto"

	"google.golang.org/grpc"
)

var (
	portFlag = flag.Int("port", 5000, "The port the controller will listen on")
)

func main() {
	flag.Parse()

	var port int
	if os.Getenv("PORT") == "" {
		port = *portFlag
	} else {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

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
