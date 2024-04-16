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

	pixelDeltaU  t.Vec3
	pixelDeltaV  t.Vec3
	defocusDiskU t.Vec3
	defocusDiskV t.Vec3
	focusDist    float64

	pixel00 t.Vec3
}

func NewCamera(imgW, imgH int, position t.Vec3, lookAt t.Vec3, fov float64, focusDist float64, defocusAngle float64) Camera {
	newFov := math.Max(1, math.Min(179.9, fov))

	if focusDist == 0 {
		focusDist = position.SubNew(lookAt).Length()
	}
	log.Printf("Focus distance: %.1f", focusDist)

	c := Camera{
		Position:  position,
		LookAt:    lookAt,
		FOV:       newFov,
		focusDist: focusDist,
	}

	log.Printf("Creating camera at %v looking at %v with FOV %.1f", position, lookAt, newFov)

	// 99% of this function is copied from the 'Ray Tracing In One Weekend' book
	// https://raytracing.github.io/books/RayTracingInOneWeekend.html#positionablecamera

	// Viewport details
	//focalLength := position.SubNew(lookAt).Length()
	theta := fov * math.Pi / 180.0
	h := 2 * math.Tan(theta/2.0)

	// Calculate the width and height of the viewport
	viewHeight := 2 * h * focusDist
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
	focalTimesW := w.MultNew(focusDist)

	upperLeft := position.SubNew(focalTimesW).SubNew(viewUHalf).SubNew(viewVHalf)
	c.pixel00 = upperLeft

	// Calculate the camera defocus disk basis vectors.
	defocusRadius := focusDist * math.Tan(math.Pi*defocusAngle/360.0)
	c.defocusDiskU = u.MultNew(defocusRadius)
	c.defocusDiskV = v.MultNew(defocusRadius)

	return c
}

func (c Camera) MakeRay(pixelX, pixelY int) Ray {
	// Randomly sample within the pixel
	offsetX := rand.Float64() - 0.5
	offsetY := rand.Float64() - 0.5
	pX := float64(pixelX) + offsetX
	pY := float64(pixelY) + offsetY

	pixelSample := c.pixel00.AddNew(c.pixelDeltaU.MultNew(pX)).AddNew(c.pixelDeltaV.MultNew(pY))

	origin := c.Position
	if c.focusDist > 0 {
		diskRandom := t.RandVecDisk(true)
		diskOffset := c.defocusDiskU.MultNew(diskRandom.X).AddNew(c.defocusDiskV.MultNew(diskRandom.Y))
		origin = origin.AddNew(diskOffset)
	}

	return Ray{
		Origin: origin,
		Dir:    pixelSample.SubNew(origin).NormalizeNew(),
	}
}
