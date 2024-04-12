package raytrace

import (
	"math"
	"nanoray/lib/proto"
	t "nanoray/lib/tuples"
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

func (c Camera) MakeRay(x, y float64, imgDetail *proto.ImageDetails) Ray {
	width := float64(imgDetail.Width)
	height := float64(imgDetail.Height)

	halfHeight := math.Tan(c.Theta / 2.0)
	halfWidth := imgDetail.AspectRatio * halfHeight

	dir := t.Vec3{
		X: halfWidth*2*float64(x)/width - halfWidth,
		Y: -(halfHeight*2*float64(y)/height - halfHeight), // NOTE: Invert Y so +ve is up
		Z: -1,
	}

	return NewRay(c.Position, dir)
}
