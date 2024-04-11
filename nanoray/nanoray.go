package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"nanoray/shared/proto"
	"nanoray/shared/raytrace"
	t "nanoray/shared/tuples"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: nanoray <output.png>")
	}

	outputFile := os.Args[1]

	err := os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	c := raytrace.NewCamera(t.Vec3{0, -5, 170}, t.Zero(), 60)

	s := raytrace.Scene{}
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

	log.Println("🚀 Rendering started...")

	img := Generate(c, s)

	log.Println("📷 Rendering complete")
	log.Println("🔹 ⌚ Time:", raytrace.Stat.Time)
	log.Println("🔹 🔦 Rays:", raytrace.Stat.Rays)

	log.Println("💾 Writing: " + outputFile)

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}

func Generate(c raytrace.Camera, s raytrace.Scene) image.Image {
	raytrace.Stat = raytrace.Stats{}
	raytrace.Stat.Start = time.Now()

	totalJobs := 60
	imgAspect := 4.0 / 3.0
	imgW := 1800
	imgH := int(float64(imgW) / imgAspect)
	jobW := imgW
	jobH := imgH / totalJobs

	jobCount := 0
	results := make(chan *proto.JobResult)
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{imgW, imgH}})

	for y := 0; y < imgH; y += jobH {
		for x := 0; x < imgW; x += jobW {
			go func() {
				j := &proto.JobRequest{
					Id:              int32(jobCount),
					SceneId:         "__local__", // Not used in CLI version of the renderer
					Width:           int32(jobW),
					Height:          int32(jobH),
					X:               int32(x),
					Y:               int32(y),
					SamplesPerPixel: int32(30),
					ChunkSize:       int32(1),
					ImageDetails: &proto.ImageDetails{
						Width:       int32(imgW),
						Height:      int32(imgH),
						AspectRatio: imgAspect,
					},
				}

				jobCount++

				// Work all happens here, with the job + scene + camera
				results <- raytrace.RenderJob(j, s, c)
			}()
		}
	}

	// Wait for all jobs to complete
	for res := range results {
		log.Printf("Job %d complete", res.Job.Id)
		jobCount--

		srcImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(res.Job.Width), int(res.Job.Height)}})
		srcImg.Pix = res.ImageData

		// Reconstruction of the main image from each job part
		draw.Draw(img, image.Rect(int(res.Job.X), int(res.Job.Y), int(res.Job.X+res.Job.Width), int(res.Job.Y+res.Job.Height)), srcImg, image.Point{0, 0}, draw.Src)

		if jobCount == 0 {
			break
		}
	}

	raytrace.Stat.End = time.Now()
	raytrace.Stat.Time = raytrace.Stat.End.Sub(raytrace.Stat.Start)

	return img
}
