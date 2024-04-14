package raytrace

import (
	"log"
	"math"
	"math/rand"
	t "nanoray/lib/tuples"
)

type Camera struct {
	Position t.Vec3
	LookAt   t.Vec3
	FOV      float64

	pixelDeltaU t.Vec3
	pixelDeltaV t.Vec3
	pixel00     t.Vec3
}

func NewCamera(imgW, imgH int, position t.Vec3, lookAt t.Vec3, fov float64) Camera {
	newFov := math.Max(1, math.Min(179.9, fov))

	c := Camera{
		Position: position,
		LookAt:   lookAt,
		FOV:      newFov,
	}

	log.Printf("Creating camera at %v looking at %v with FOV %.1f", position, lookAt, newFov)

	// 99% of this function is copied from the 'Ray Tracing In One Weekend' book
	// https://raytracing.github.io/books/RayTracingInOneWeekend.html#positionablecamera

	// Viewport details
	focalLength := position.SubNew(lookAt).Length()
	theta := fov * math.Pi / 180.0
	h := 2 * math.Tan(theta/2.0)

	// Calculate the width and height of the viewport
	viewHeight := 2 * h * focalLength
	viewWidth := viewHeight * (float64(imgW) / float64(imgH))

	// Calculate the basis vectors for the camera
	w := (position.SubNew(lookAt)).NormalizeNew()
	upVector := t.Vec3{0, 1, 0}
	u := upVector.Cross(w).NormalizeNew()
	v := w.Cross(u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges
	viewU := u.MultNew(viewWidth)
	viewV := v.MultNew(-viewHeight)

	// This is what it's all for - vectors used to calculate the rays from camera to pixels
	c.pixelDeltaU = viewU.DivNew(float64(imgW))
	c.pixelDeltaV = viewV.DivNew(float64(imgH))

	viewUHalf := viewU.DivNew(2)
	viewVHalf := viewV.DivNew(2)
	focalLTimesW := w.MultNew(focalLength)

	upperLeft := position.SubNew(focalLTimesW).SubNew(viewUHalf).SubNew(viewVHalf)
	c.pixel00 = upperLeft

	return c
}

func (c Camera) MakeRay(pixelX, pixelY int) Ray {
	// Randomly sample within the pixel
	offsetX := rand.Float64() - 0.5
	offsetY := rand.Float64() - 0.5
	pX := float64(pixelX) + offsetX
	pY := float64(pixelY) + offsetY

	pixelSample := c.pixel00.AddNew(c.pixelDeltaU.MultNew(pX)).AddNew(c.pixelDeltaV.MultNew(pY))

	return Ray{
		Origin: c.Position,
		Dir:    pixelSample.SubNew(c.Position).NormalizeNew(),
	}
}
