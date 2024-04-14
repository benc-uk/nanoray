package raytrace

import (
	"image"
	"log"
	"nanoray/lib/proto"
	t "nanoray/lib/tuples"
	"sync"
	"time"
)

// Enum for render status
const (
	READY = iota
	STARTED
	COMPLETE
	FAILED
)

type Statistics struct {
	Start time.Time
	End   time.Time
	Time  time.Duration
	Rays  int
}

type NetworkRender struct {
	Lock         sync.Mutex
	Status       int
	JobQueue     sync.Map
	JobsTotal    int
	JobsComplete int
	Start        time.Time
	OutputName   string
}

type Render struct {
	Width           int     `yaml:"width"`
	Height          int     `yaml:"height"`
	AspectRatio     float64 `yaml:"aspectRatio"`
	SamplesPerPixel int     `yaml:"samplesPerPixel"`
	MaxDepth        int     `yaml:"maxDepth"`
}

var (
	Stats Statistics = Statistics{}
)

func NewRender(width int, aspectRatio float64) Render {
	// Create a default render object
	return Render{
		Width:           width,
		Height:          int(float64(width) / aspectRatio),
		AspectRatio:     aspectRatio,
		SamplesPerPixel: 10,
		MaxDepth:        5,
	}
}

func (r Render) MakeImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
}

func (r Render) ImageDetails() *proto.ImageDetails {
	return &proto.ImageDetails{
		Width:       int32(r.Width),
		Height:      int32(r.Height),
		AspectRatio: r.AspectRatio,
	}
}

func RenderJob(job *proto.JobRequest, s Scene, c Camera) *proto.JobResult {
	log.Printf("Rendering job %4d: slice:%4d/%4d samp:%d", job.Id, job.Height, job.Y, job.SamplesPerPixel)

	samples := int(job.SamplesPerPixel)
	sampleScale := 1.0 / float64(samples)

	jobImg := image.NewRGBA(image.Rect(0, 0, int(job.Width), int(job.Height)))

	for y := 0; y < int(job.Height); y += 1 {
		for x := 0; x < int(job.Width); x += 1 {
			// Note that x and y are relative to the job, NOT the image
			pixelX := int(job.X) + x
			pixelY := int(job.Y) + y

			pixel := t.Black()

			// Path tracing uses many samples
			for i := 0; i < samples; i++ {
				ray := c.MakeRay(pixelX, pixelY)
				sample := ray.Shade(s, 0, int(job.MaxDepth))
				pixel.AddSome(sample, sampleScale)
			}

			// TODO: Remove hard-coded gamma
			jobImg.Set(x, y, pixel.ToRGBA(0.8))
		}
	}

	jobRes := proto.JobResult{
		ImageData: jobImg.Pix,
		Job:       job,
	}

	return &jobRes
}
