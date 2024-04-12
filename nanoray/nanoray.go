package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"log"
	"nanoray/lib/proto"
	rt "nanoray/lib/raytrace"
	t "nanoray/lib/tuples"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	flag.Usage = func() {
		log.Println("NanoRay - A path based ray tracer")
		flag.PrintDefaults()
	}

	inputScene := flag.String("file", "", "Scene file to render, in YAML format")
	outputFile := flag.String("output", "render.png", "Rendered output PNG file name")
	width := flag.Int("width", 800, "Width of the output image")
	aspectRatio := flag.Float64("aspect", 16.0/9.0, "Aspect ratio of the output image")
	samplesPP := flag.Int("samples", 20, "Samples per pixel, higher values give better quality but slower rendering")
	maxDepth := flag.Int("depth", 5, "Maximum ray recursion depth")
	chunky := flag.Int("chunk", 1, "Speed up rendering with chunky pixels")

	flag.Parse()

	if *inputScene == "" {
		flag.PrintDefaults()
		log.Fatal("No scene file provided")
	}

	err := os.MkdirAll(filepath.Dir(*outputFile), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cam := rt.NewCamera(t.Vec3{0, 0, 0}, t.Zero(), 60)

	sceneData, err := os.ReadFile("scenes/test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	scene, err := rt.ParseScene(string(sceneData))
	if err != nil {
		log.Fatal(err)
	}

	render := rt.NewRender(*width, *aspectRatio)
	render.SamplesPerPixel = *samplesPP
	render.MaxDepth = *maxDepth
	render.PixelChunk = *chunky

	log.Println("🚀 Rendering started...")

	img := Generate(cam, *scene, render)

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

	totalJobs := runtime.NumCPU()
	jobW := render.Width
	jobH := render.Height / totalJobs

	jobCount := 0
	results := make(chan *proto.JobResult)
	imageOut := render.MakeImage()

	for y := 0; y < render.Height; y += jobH {
		for x := 0; x < render.Width; x += jobW {
			go func() {
				j := &proto.JobRequest{
					Id:              int32(jobCount),
					SceneId:         "", // Not used in CLI version of the renderer
					Width:           int32(jobW),
					Height:          int32(jobH),
					X:               int32(x),
					Y:               int32(y),
					SamplesPerPixel: int32(render.SamplesPerPixel),
					ChunkSize:       int32(render.PixelChunk),
					MaxDepth:        int32(render.MaxDepth),
					ImageDetails:    render.ToProto(),
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

		srcImg := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(res.Job.Width), int(res.Job.Height)}})
		srcImg.Pix = res.ImageData

		// Reconstruction of the main image from each job part
		draw.Draw(imageOut, image.Rect(int(res.Job.X), int(res.Job.Y), int(res.Job.X+res.Job.Width), int(res.Job.Y+res.Job.Height)), srcImg, image.Point{0, 0}, draw.Src)

		if jobCount == 0 {
			break
		}
	}

	rt.Stats.End = time.Now()
	rt.Stats.Time = rt.Stats.End.Sub(rt.Stats.Start)

	return imageOut
}
