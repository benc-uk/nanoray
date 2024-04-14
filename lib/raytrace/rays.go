package raytrace

import (
	"math"
	t "nanoray/lib/tuples"
)

// ============================================================
// Rays
// ============================================================

type Ray struct {
	Origin t.Vec3
	Dir    t.Vec3
}

func NewRay(origin, direction t.Vec3) Ray {
	Stats.Rays++
	return Ray{origin, direction}
}

func (r Ray) Shade(s Scene, depth int, maxDepth int) t.RGB {
	if depth > maxDepth {
		return t.Black()
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
		// Hit something, scatter a new ray from surface material
		didScatter, scatterRay, attenColour := hit.Obj.Material.scatter(r, *hit)
		if didScatter {
			// Recurse and shade the scattered ray
			scatterColour := scatterRay.Shade(s, depth+1, maxDepth)
			scatterColour.Mult(attenColour)
			return scatterColour
		}

		return t.Black()
	}

	unitDirection := r.Dir.NormalizeNew()
	a := 0.5 * (-unitDirection.Y + 1.0)
	return t.White().Blend(t.RGB{0.5, 0.7, 1.0}, a)
}

func (r Ray) GetPoint(t float64) t.Vec3 {
	p := r.Dir.MultScalarNew(t)
	p.Add(r.Origin)
	return p
}

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
