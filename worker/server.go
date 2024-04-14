package main

import (
	"context"
	"log"

	"nanoray/lib/controller"
	pb "nanoray/lib/proto"
	"nanoray/lib/raytrace"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedWorkerServer
}

var scene *raytrace.Scene
var camera *raytrace.Camera

func (s *server) NewJob(ctx context.Context, job *pb.JobRequest) (*pb.Void, error) {
	log.Printf("Received job: %d", job.Id)

	go runJob(job)

	return nil, nil
}

func (s *server) Ping(ctx context.Context, in *pb.Void) (*pb.Void, error) {
	return nil, nil
}

func runJob(job *pb.JobRequest) error {
	if scene == nil {
		return status.Errorf(codes.FailedPrecondition, "No scene loaded")
	}

	res := raytrace.RenderJob(job, *scene, *camera)
	res.Worker = &workerInfo

	_, err := controller.Client.JobComplete(context.Background(), res)
	if err != nil {
		log.Printf("Failed to send completed job result: %s", err.Error())
		return err
	}

	return nil
}

func (s *server) PrepareRender(ctx context.Context, in *pb.PrepRenderRequest) (*pb.Void, error) {
	log.Printf("Preparing render with new scene & camera data")

	sceneData := in.SceneData
	sceneNew, cameraNew, err := raytrace.ParseScene(sceneData, int(in.ImageDetails.Width), int(in.ImageDetails.Height))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse scene data: %s", err.Error())
	}

	scene = sceneNew
	camera = cameraNew

	return nil, nil
}
