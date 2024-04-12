package raytrace

import (
	"image"
	"log"
	"math/rand"
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
}

type Render struct {
	Width           int     `yaml:"width"`
	Height          int     `yaml:"height"`
	AspectRatio     float64 `yaml:"aspectRatio"`
	SamplesPerPixel int     `yaml:"samplesPerPixel"`
	MaxDepth        int     `yaml:"maxDepth"`
	PixelChunk      int     `yaml:"pixelChunks"`
}

var (
	Stats Statistics = Statistics{}
)

func NewRender(width int, aspectRatio float64) Render {
	return Render{
		Width:           width,
		Height:          int(float64(width) / aspectRatio),
		AspectRatio:     aspectRatio,
		SamplesPerPixel: 10,
		MaxDepth:        5,
		PixelChunk:      1,
	}
}

func (r Render) MakeImage() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
}

func (r Render) ToProto() *proto.ImageDetails {
	return &proto.ImageDetails{
		Width:       int32(r.Width),
		Height:      int32(r.Height),
		AspectRatio: r.AspectRatio,
	}
}

func RenderJob(job *proto.JobRequest, s Scene, c Camera) *proto.JobResult {
	log.Printf("Rendering job %d, w:%d, h:%d (%d,%d) samp:%d", job.Id, job.Width, job.Height, job.X, job.Y, job.SamplesPerPixel)

	samplesPerPixel := int(job.SamplesPerPixel)
	sampleScale := 1.0 / float64(samplesPerPixel)
	chunk := int(job.ChunkSize)

	jobImg := image.NewRGBA(image.Rect(0, 0, int(job.Width), int(job.Height)))

	for y := 0; y < int(job.Height); y += chunk {
		for x := 0; x < int(job.Width); x += chunk {
			// Note that x and y are relative to the job, not the image
			xf := float64(x) + float64(job.X)
			yf := float64(y) + float64(job.Y)

			pixel := t.Black()

			// Always sample in the pixel center
			ray := c.MakeRay(xf, yf, job.ImageDetails)
			sample := ray.Shade(s, 0, int(job.MaxDepth))
			pixel.AddSome(sample, sampleScale)

			// Path tracing with multiple samples
			for i := 0; i < samplesPerPixel-1; i++ {
				ray := c.MakeRay(xf+rand.Float64()-0.5, yf+rand.Float64()-0.5, job.ImageDetails)
				sample := ray.Shade(s, 0, int(job.MaxDepth))
				pixel.AddSome(sample, sampleScale)
			}

			jobImg.Set(x, y, pixel.ToRGBA())

			// For speedy chunky mode, fill in the rest of the chunk
			if chunk > 1 {
				for y2 := y; y2 < y+chunk; y2++ {
					for x2 := x; x2 < x+chunk; x2++ {
						jobImg.Set(x2, y2, pixel.ToRGBA())
					}
				}
			}
		}
	}

	jobRes := proto.JobResult{
		ImageData: jobImg.Pix,
		Job:       job,
	}

	return &jobRes
}
