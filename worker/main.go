package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"raynet/pkg/controller"
	pb "raynet/pkg/proto"

	"google.golang.org/grpc"
)

var (
	portFlag     = flag.String("port", "0", "The offset of the port to listen on")
	portBaseFlag = flag.String("portBase", "4400", "The port number base to listen on")
	hostnameFlag = flag.String("hostname", "", "Override the hostname to use for the worker")
	workerInfo   pb.WorkerInfo
)

func main() {
	// Lots of guff for port handling, env vars and flags
	// We do things a little differently using an offset, easier to run multiple workers on the same machine
	portNum := os.Getenv("PORT")
	portBase := os.Getenv("PORT_BASE")
	hostnameEnv := os.Getenv("HOSTNAME")

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
	port := int(portBaseInt + portNumInt)

	hostname := "localhost"
	if hostnameEnv != "" {
		hostname = hostnameEnv
	} else if *hostnameFlag != "" {
		hostname = *hostnameFlag
	} else {

		var err error
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "localhost"
		}
	}

	workerInfo = pb.WorkerInfo{
		Id:   generateID(fmt.Sprintf("%s:%d", hostname, port), 8),
		Host: hostname,
		Port: int32(port),
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to bind to port\n%s", err.Error())
	}

	s := grpc.NewServer()
	pb.RegisterWorkerServer(s, &server{})

	err = controller.Connect(time.Second * 20)
	if err != nil {
		log.Fatalf("Failed to connect to controller: %s", err.Error())
	}

	log.Printf("Worker started on port %d", port)

	// Serve the gRPC server but continue to register with the controller
	// Use a channel to block until the server is done
	ch := make(chan struct{})
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve gRPC server\n%s", err.Error())
		}

		close(ch)
	}()

	log.Printf("Will register worker with hostname: %s", hostname)
	_, err = controller.Client.RegisterWorker(context.Background(), &workerInfo)
	if err != nil {
		log.Fatalf("Registration failed\n%s", err.Error())
	}

	log.Printf("Registered with controller, we are ready to work!")

	// Block until the server is done
	<-ch
}

func generateID(input string, len int) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	sum := hash.Sum(nil)

	id := fmt.Sprintf("%x", sum)
	return id[:len]
}
