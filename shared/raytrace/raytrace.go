package raytrace

import (
	"image"
	"sync"
)

// Enum for render status
const (
	READY = iota
	STARTED
	COMPLETE
	FAILED
)

type Render struct {
	Width, Height int
	AspectRatio   float64
	Scene         Scene
	Image         *image.RGBA
}

type NetworkRender struct {
	Render
	Lock         sync.Mutex
	Status       int
	JobQueue     sync.Map
	JobsTotal    int
	JobsComplete int
}

func NewRender(width int, aspect float64) *Render {
	r := &Render{
		Width:       width,
		Height:      int(float64(width) / aspect),
		AspectRatio: aspect,
	}

	r.Image = image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))

	return r
}

func (r *Render) Generate(c Camera, s Scene) {
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			ray := c.Ray(x, y, *r)
			pixel := ray.Shade(s)

			r.Image.Set(x, y, pixel.ToRGBA())
		}
	}
}
