package raytrace

import (
	"math"
	t "nanoray/shared/tuples"
)

type Camera struct {
	Position  t.Vec3
	Direction t.Vec3
	FOV       float64
	Theta     float64
}

func NewCamera(position t.Vec3, direction t.Vec3, fov float64) Camera {
	return Camera{
		Position:  position,
		Direction: direction,
		FOV:       fov,
		Theta:     fov * (math.Pi / 180.0),
	}
}

func (c Camera) Ray(x, y int, r Render) Ray {
	width := float64(r.Width)
	height := float64(r.Height)

	halfHeight := math.Tan(c.Theta / 2.0)
	halfWidth := r.AspectRatio * halfHeight

	dir := t.Vec3{
		X: halfWidth*2*float64(x)/width - halfWidth,
		Y: -(halfHeight*2*float64(y)/height - halfHeight),
		Z: -1,
	}

	return Ray{
		Origin:    t.Zero(),
		Direction: dir,
	}
}
