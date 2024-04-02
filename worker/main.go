package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"raynet/shared/controller"
	pb "raynet/shared/proto"

	"google.golang.org/grpc"
)

var (
	portFlag     = flag.Int("port", 4400, "The port worker will listen on")
	hostnameFlag = flag.String("hostname", "", "Override the hostname to use for the worker")
	maxJobsFlag  = flag.Int("maxjobs", runtime.NumCPU(), "The maximum number of jobs to run concurrently")
	workerInfo   pb.WorkerInfo
)

func main() {
	flag.Parse()

	var port int
	if os.Getenv("PORT") == "" {
		port = *portFlag
	} else {
		port, _ = strconv.Atoi(os.Getenv("PORT"))
	}

	var maxJobs int
	if os.Getenv("MAX_JOBS") == "" {
		maxJobs = *maxJobsFlag
	} else {
		maxJobs, _ = strconv.Atoi(os.Getenv("MAX_JOBS"))
	}

	hostname, _ := os.Hostname()
	if os.Getenv("HOSTNAME") != "" {
		hostname = os.Getenv("HOSTNAME")
	} else if *hostnameFlag != "" {
		hostname = *hostnameFlag
	}

	workerInfo = pb.WorkerInfo{
		Id:      generateID(fmt.Sprintf("%s:%d", hostname, port), 8),
		Host:    hostname,
		Port:    int32(port),
		MaxJobs: int32(maxJobs),
	}

	log.Printf("Starting worker, will handle max jobs: %d", maxJobs)

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

	// Start gRPC server in a goroutine, use a channel to block until it's done
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
