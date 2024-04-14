package raytrace

import (
	"math"
	v "nanoray/lib/tuples"
	"strconv"
)

// ============================================================
// Rays
// ============================================================

type Ray struct {
	Origin v.Vec3
	Dir    v.Vec3
}

func NewRay(origin, direction v.Vec3) Ray {
	Stats.Rays++
	return Ray{origin, direction}
}

func (r Ray) Shade(s Scene, depth int, maxDepth int) v.RGB {
	if depth > maxDepth {
		return v.Black()
	}

	interval := Interval{0.001, math.MaxFloat64}
	var hit *Hit = nil

	for _, o := range s.Objects {
		didHit, objHit := o.Hit(r, interval)
		if didHit {
			interval.Max = objHit.T
			hit = &objHit
		}
	}

	if hit != nil {
		normal := hit.Normal
		randDir := normal.AddNew(v.RandVecSphere(true))
		randRay := NewRay(hit.P, randDir)

		hitColour := randRay.Shade(s, depth+1, maxDepth)

		// HACK: This is a placeholder
		objColour := v.White().Blend(hit.O.Colour, 0.7)
		hitColour.Mult(objColour)

		return hitColour
	}

	unitDirection := r.Dir.NormalizeNew()
	a := 0.5 * (-unitDirection.Y + 1.0)
	return v.White().Blend(v.RGB{0.5, 0.7, 1.0}, a)
}

func (r Ray) GetPoint(t float64) v.Vec3 {
	p := r.Dir.MultScalarNew(t)
	p.Add(r.Origin)
	return p
}

// ============================================================
// Object hits
// ============================================================

type Hit struct {
	T      float64
	P      v.Vec3
	Normal v.Vec3
	O      Object
}

type Hitable interface {
	Hit(r Ray, i Interval) (bool, Hit)
}

func (h Hit) String() string {
	return "Hit{" + h.P.String() + ", " + h.Normal.String() + "}"
}

// ============================================================
// Simple interval between two numbers
// ============================================================

type Interval struct {
	Min float64
	Max float64
}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}

	return x
}

func (i Interval) String() string {
	minStr := strconv.FormatFloat(i.Min, 'f', -1, 64)
	maxStr := strconv.FormatFloat(i.Max, 'f', -1, 64)
	return "Interval{" + minStr + ", " + maxStr + "}"
}

// ============================================================
// Materials (placeholder for now)
// ============================================================

type DiffuseMaterial struct {
}

type Material interface {
	scatter(r Ray, hit Hit) (bool, Ray)
}
