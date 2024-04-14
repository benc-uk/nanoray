package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"nanoray/lib/proto"
	"os"
	"path/filepath"
	"runtime"
	"time"

	rt "nanoray/lib/raytrace"
)

func main() {
	flag.Usage = func() {
		log.Println("NanoRay - A path based ray tracer")
		flag.PrintDefaults()
	}

	inputFile := flag.String("file", "", "Scene file to render, in YAML format")
	outputFile := flag.String("output", "render.png", "Rendered output PNG file name")
	width := flag.Int("width", 800, "Width of the output image")
	aspectRatio := flag.Float64("aspect", 16.0/9.0, "Aspect ratio of the output image")
	samplesPP := flag.Int("samples", 20, "Samples per pixel, higher values give better quality but slower rendering")
	maxDepth := flag.Int("depth", 5, "Maximum ray recursion depth")

	flag.Parse()

	if *inputFile == "" {
		flag.PrintDefaults()
		log.Fatal("No scene file provided")
	}

	err := os.MkdirAll(filepath.Dir(*outputFile), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	sceneData, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	render := rt.NewRender(*width, *aspectRatio)
	render.SamplesPerPixel = *samplesPP
	render.MaxDepth = *maxDepth
	// cam := rt.NewCamera(render.Width, render.Height, t.Vec3{240, 150, 120}, t.Vec3{0, 0, -100}, 50)

	scene, camera, err := rt.ParseScene(string(sceneData), render.Width, render.Height)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 Rendering started...")

	img := Generate(*camera, *scene, render)

	log.Println("📷 Rendering complete")
	log.Println("🔹 ⌚ Time:", rt.Stats.Time)
	log.Printf("🔹 🔦 Rays: %f Mil", float64(rt.Stats.Rays)/1000000.0)

	log.Println("💾 Writing: " + *outputFile)

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}

func Generate(c rt.Camera, s rt.Scene, render rt.Render) image.Image {
	rt.Stats.Start = time.Now()
	imageOut := render.MakeImage()

	totalJobs := runtime.NumCPU()
	if totalJobs > render.Height {
		totalJobs = render.Height
	}
	// totalJobs := 1
	jobW := render.Width
	jobH := render.Height / totalJobs

	jobCount := 0
	results := make(chan *proto.JobResult)

	for y := 0; y < render.Height; y += jobH {
		for x := 0; x < render.Width; x += jobW {
			go func() {
				j := &proto.JobRequest{
					Id:              int32(jobCount),
					Width:           int32(jobW),
					Height:          int32(jobH),
					X:               int32(x),
					Y:               int32(y),
					SamplesPerPixel: int32(render.SamplesPerPixel),
					MaxDepth:        int32(render.MaxDepth),
					ImageDetails:    render.ImageDetails(),
				}

				jobCount++

				// Work all happens here, with the job + scene + camera
				results <- rt.RenderJob(j, s, c)
			}()
		}
	}

	// Wait for all jobs to complete
	for res := range results {
		log.Printf("Job %d complete", res.Job.Id)
		jobCount--

		jobImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(res.Job.Width), int(res.Job.Height)}})
		jobImg.Pix = res.ImageData

		// Reconstruction of the main image from each job part
		draw.Draw(imageOut, image.Rect(int(res.Job.X), int(res.Job.Y), int(res.Job.X+res.Job.Width), int(res.Job.Y+res.Job.Height)), jobImg, image.Point{0, 0}, draw.Src)

		if jobCount == 0 {
			break
		}
	}

	rt.Stats.End = time.Now()
	rt.Stats.Time = rt.Stats.End.Sub(rt.Stats.Start)

	return imageOut
}
