package raytrace

import (
	"math"
	t "nanoray/lib/tuples"
)

type Object struct {
	Id       string
	Position t.Vec3
	Material Material
	Colour   t.RGB
}

type Sphere struct {
	Object
	Radius float64
}

func NewSphere(position t.Vec3, radius float64) Sphere {
	return Sphere{
		Object: Object{
			Position: position,
			Id:       "sphere_" + GenerateID("sphere") + position.String(),
		},

		Radius: radius,
	}
}

func (s Sphere) Hit(r Ray, interval Interval) (bool, Hit) {
	oc := r.Origin
	oc.Sub(s.Position)

	a := r.Direction.Dot(r.Direction)
	b := oc.Dot(r.Direction)
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
		hit := Hit{
			T:      t,
			P:      r.GetPoint(t),
			Normal: r.GetPoint(t).SubNew(s.Position).NormalizeNew(),
			O:      s.Object,
		}

		return true, hit
	}

	return false, Hit{}
}
