package raytrace

import (
	"math"
	t "nanoray/lib/tuples"
)

type Sphere struct {
	Object
	Radius float64
}

func NewSphere(position t.Vec3, radius float64) (*Sphere, error) {
	if radius <= 0 {
		return nil, ErrInvalidRadius
	}

	return &Sphere{
		Object: Object{
			Position: position,
			ID:       "sphere_" + GenerateID("sphere") + position.String(),
		},

		Radius: radius,
	}, nil
}

func (s Sphere) Hit(r Ray, interval Interval) (bool, Hit) {
	oc := r.Origin
	oc.Sub(s.Position)

	a := r.Dir.Dot(r.Dir)
	b := oc.Dot(r.Dir)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := b*b - a*c
	if discriminant < 0 {
		return false, Hit{}
	}

	sqrtDisc := math.Sqrt(discriminant)

	t1 := (-b - sqrtDisc) / a
	t2 := (-b + sqrtDisc) / a
	t := t1
	if t1 < interval.Min {
		t = t2
	}

	if t > interval.Min && t < interval.Max {
		normal := r.GetPoint(t).SubNew(s.Position).NormalizeNew()
		hit := r.MakeHit(t, normal, s.Object)

		return true, hit
	}

	return false, Hit{}
}
