package raytrace

import (
	"math"
	t "nanoray/lib/tuples"
)

type Ray struct {
	Origin t.Vec3
	Dir    t.Vec3
}

// -
// Create a new ray, use this to keep track of rays cast
// -
func NewRay(origin, direction t.Vec3) Ray {
	Stats.Rays++
	return Ray{origin, direction}
}

// -
// Cast a ray into the scene and return the colour it hits
// This is the core of the entire raytracing algorithm and is recursive
// -
func (r Ray) Shade(scene Scene, depth int, maxDepth int) t.RGB {
	if depth > maxDepth {
		return t.Black()
	}

	interval := Interval{0.001, math.MaxFloat64}
	var hit *Hit = nil

	// Main ray collision loop against all objects
	for _, obj := range scene.Objects {
		// Find the closest hit
		didHit, objHit := obj.Hit(r, interval)
		if didHit {
			interval.Max = objHit.T
			hit = &objHit
		}
	}

	if hit != nil {
		// Hit something, scatter a new ray from surface based on material
		scattered, scatterRay, attenColour := hit.Obj.Material.scatter(r, *hit)
		if scattered {
			// Recurse and shade the scattered ray
			scatterColour := scatterRay.Shade(scene, depth+1, maxDepth)
			// Magic to blend the scattered colour with the attenuation colour
			scatterColour.Mult(attenColour)
			return scatterColour
		}

		return t.Black()
	}

	// On miss return a fake sky gradient thing
	unitDirection := r.Dir.NormalizeNew()
	a := 0.5 * (-unitDirection.Y + 1.0)
	return t.White().Blend(t.RGB{0.5, 0.7, 1.0}, a)
}

// -
// Get a point along the ray at a distance t
// -
func (r Ray) GetPoint(t float64) t.Vec3 {
	p := r.Dir.MultScalarNew(t)
	p.Add(r.Origin)
	return p
}

// -
// Helper to create a hit object from a ray and a distance
// Does some checks for inside/outside hits
// -
func (r *Ray) MakeHit(t float64, normal t.Vec3, obj Object) Hit {
	h := Hit{
		T:      t,
		Pos:    r.GetPoint(t),
		Obj:    obj,
		Front:  true,
		Normal: normal,
	}

	// Check if hit was inside or outside
	if r.Dir.Dot(normal) > 0 {
		h.Front = false
		h.Normal = normal.NegateNew()
	}

	return h
}
