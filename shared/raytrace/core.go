package raytrace

import (
	"math"
	v "nanoray/shared/tuples"
	"strconv"
)

type Ray struct {
	Origin    v.Vec3
	Direction v.Vec3
}

func (r Ray) Shade(s Scene) v.RGB {
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
		objColor := v.RGB{normal.X + 1, normal.Y + 1, normal.Z + 1}.MultScalarNew(0.5)
		return objColor
	}

	t := 0.5 * (-r.Direction.Y + 1.0)
	out := v.White().Blend(v.RGB{0.2, 0.4, 1.0}, t)
	return out
}

func (r Ray) GetPoint(t float64) v.Vec3 {
	p := r.Direction.MultScalarNew(t)
	p.Add(r.Origin)
	return p
}

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

func (i Interval) String() string {
	minStr := strconv.FormatFloat(i.Min, 'f', -1, 64)
	maxStr := strconv.FormatFloat(i.Max, 'f', -1, 64)
	return "Interval{" + minStr + ", " + maxStr + "}"
}

// ============================================================

type DiffuseMaterial struct {
}

type Material interface {
	scatter(r Ray, hit Hit) (bool, Ray)
}
