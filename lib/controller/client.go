package controller

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	pb "nanoray/lib/proto"
)

var Client pb.ControllerClient

// Connect to the controller, waiting for the connection to be ready
func Connect(waitTime time.Duration) error {
	controllerAddr := os.Getenv("CONTROLLER_ADDR")
	if controllerAddr == "" {
		controllerAddr = "localhost:5000"
	}

	log.Printf("Connecting to controller at %s", controllerAddr)

	conn, err := grpc.Dial(controllerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), waitTime)
		defer cancel()

		change := conn.WaitForStateChange(ctx, conn.GetState())
		if !change {
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		if conn.GetState() == connectivity.Ready {
			break
		}
	}

	log.Printf("Connected to controller")
	Client = pb.NewControllerClient(conn)
	return nil
}
