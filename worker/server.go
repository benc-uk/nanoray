package main

import (
	"context"
	"log"
	"time"

	"nanoray/lib/controller"
	pb "nanoray/lib/proto"
	"nanoray/lib/raytrace"
	t "nanoray/lib/tuples"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedWorkerServer
}

var scene *raytrace.Scene
var c raytrace.Camera

func init() {
	c = raytrace.NewCamera(t.Vec3{0, 0, 0}, t.Zero(), 60)
}

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

	res := raytrace.RenderJob(job, *scene, c)
	res.Worker = &workerInfo

	_, err := controller.Client.JobComplete(context.Background(), res)
	if err != nil {
		log.Printf("Failed to send completed job result: %s", err.Error())
		return err
	}

	return nil
}

func (s *server) LoadScene(ctx context.Context, in *pb.SceneRaw) (*pb.Void, error) {
	log.Printf("Received new scene data, will parse it")

	sceneData := in.Data
	sceneNew, err := raytrace.ParseScene(sceneData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to parse scene data: %s", err.Error())
	}

	scene = sceneNew

	time.Sleep(1 * time.Second)

	return nil, nil
}
