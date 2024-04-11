package main

import (
	"context"
	"log"

	"nanoray/shared/controller"
	pb "nanoray/shared/proto"
	"nanoray/shared/raytrace"
	t "nanoray/shared/tuples"
)

type server struct {
	pb.UnimplementedWorkerServer
}

var s raytrace.Scene
var c raytrace.Camera

func init() {
	s = raytrace.Scene{}
	s.MaxDepth = 5
	sphere := raytrace.NewSphere(t.Vec3{-120, 0, -80}, 50)
	sphere.Colour = t.RGB{1, 0, 0}
	s.AddObject(sphere)

	big := raytrace.NewSphere(t.Vec3{0, -9050, -120}, 9000)
	big.Colour = t.RGB{1, 1, 1}
	s.AddObject(big)

	sphere3 := raytrace.NewSphere(t.Vec3{0, 0, -100}, 50)
	sphere3.Colour = t.RGB{1, 1, 0}
	s.AddObject(sphere3)

	sphere4 := raytrace.NewSphere(t.Vec3{120, 0, -80}, 50)
	sphere4.Colour = t.RGB{0, 0.9, 0}
	s.AddObject(sphere4)

	c = raytrace.NewCamera(t.Vec3{0, -5, 170}, t.Zero(), 60)
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
	// Get the scene from the server NOT USED YET
	// scene, err := controller.Client.GetScene(context.Background(), &pb.SceneRequest{Id: job.SceneId})
	// log.Printf("Got scene: %s", scene.Data)

	// 🔥🔥🔥 TEMP SCENE

	res := raytrace.RenderJob(job, s, c)
	res.Worker = &workerInfo

	_, err := controller.Client.JobComplete(context.Background(), res)
	if err != nil {
		log.Printf("Failed to send completed job result: %s", err.Error())
		return err
	}

	return nil
}
