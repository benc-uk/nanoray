package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"nanoray/lib/proto"
	"os"
	"path/filepath"
	"strings"
	"time"

	rt "nanoray/lib/raytrace"
)

func main() {
	flag.Usage = func() {
		fmt.Println("NanoRay - A path based parallel ray tracer")
		flag.PrintDefaults()
	}

	inputFile := flag.String("file", "", "Scene file to render, in YAML format")
	outputFile := flag.String("output", "render.png", "Rendered output PNG file name")
	width := flag.Int("width", 800, "Width of the output image")
	aspectRatio := flag.Float64("aspect", 16.0/9.0, "Aspect ratio of the output image")
	samplesPP := flag.Int("samples", 10, "Samples per pixel, higher values give better quality but slower rendering")
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

func Generate(cam rt.Camera, scene rt.Scene, render rt.Render) image.Image {
	rt.Stats.Start = time.Now()
	imageOut := render.MakeImage()

	totalJobs := 16 //runtime.NumCPU()
	if totalJobs > render.Height {
		totalJobs = render.Height
	}
	jobW := render.Width
	jobH := render.Height / totalJobs

	jobCount := 0
	// Create a channel to receive results from the jobs and synchronize them
	results := make(chan *proto.JobResult)

	for y := 0; y < render.Height; y += jobH {
		for x := 0; x < render.Width; x += jobW {
			job := &proto.JobRequest{
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

			// Use goroutines to parallelize the rendering
			go func() {
				// Work all happens here, with the job + scene + camera
				results <- rt.RenderJob(job, scene, cam)
			}()
		}
	}

	// Wait for all jobs to complete
	for res := range results {
		jobCount--

		// Shonky progress bar
		spaces := int((float64(jobCount) / float64(totalJobs) * 30.0))
		blocks := 30 - spaces
		if blocks > 0 {
			fmt.Printf("\033[2K\r🎥 Rendering Progress: [%s%s]", strings.Repeat("█", blocks), strings.Repeat(" ", spaces))
		}

		jobImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(res.Job.Width), int(res.Job.Height)}})
		jobImg.Pix = res.ImageData

		// Reconstruction of the main image from each job part
		draw.Draw(imageOut, image.Rect(int(res.Job.X), int(res.Job.Y), int(res.Job.X+res.Job.Width), int(res.Job.Y+res.Job.Height)), jobImg, image.Point{0, 0}, draw.Src)

		if jobCount == 0 {
			close(results)
			fmt.Println()
			break
		}
	}

	rt.Stats.End = time.Now()
	rt.Stats.Time = rt.Stats.End.Sub(rt.Stats.Start)

	return imageOut
}
