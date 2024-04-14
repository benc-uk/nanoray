package raytrace

import (
	t "nanoray/lib/tuples"
)

type Material interface {
	scatter(r Ray, hit Hit) (bool, Ray, t.RGB)
}

// ============================================================
// Diffuse Lambertian material
// ============================================================

type DiffuseMaterial struct {
	Albedo t.RGB
}

func NewDiffuseMaterial(c t.RGB) DiffuseMaterial {
	return DiffuseMaterial{
		Albedo: c,
	}
}

func (m DiffuseMaterial) scatter(r Ray, hit Hit) (bool, Ray, t.RGB) {
	scatterDir := hit.Normal.AddNew(t.RandVecSphere(true))
	if scatterDir.IsNearZero() {
		scatterDir = hit.Normal
	}

	scattered := NewRay(hit.Pos, scatterDir)

	return true, scattered, m.Albedo
}

// ============================================================
// Metal like material
// ============================================================

type MetalMaterial struct {
	Albedo t.RGB
	Fuzz   float64
}

func NewMetalMaterial(c t.RGB, f float64) MetalMaterial {
	return MetalMaterial{
		Albedo: c,
		Fuzz:   f,
	}
}

func (m MetalMaterial) scatter(r Ray, hit Hit) (bool, Ray, t.RGB) {
	reflected := r.Dir.Reflect(hit.Normal)
	reflected.Normalize()

	fuzz := t.RandVecSphere(false)
	fuzz.MultScalar(m.Fuzz)
	reflected.Add(fuzz)

	scattered := NewRay(hit.Pos, reflected)

	return scattered.Dir.Dot(hit.Normal) > 0, scattered, m.Albedo
}
