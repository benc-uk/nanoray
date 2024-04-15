package raytrace

import (
	"math"
	t "nanoray/lib/tuples"
)

// Material interface for all materials
type Material interface {
	// Used when shading a hit point on this material
	scatter(r Ray, hit Hit) (didScatter bool, scattedRay Ray, attenuation t.RGB)

	// Used when calculating emitted light from this material
	emitted() t.RGB
}

// ============================================================
// DiffuseMaterial represents a perfect Lambertian or ideal diffuse material
// ============================================================

type DiffuseMaterial struct {
	Albedo t.RGB
}

// -
// Create a new DiffuseMaterial with given albedo color
// -
func NewDiffuseMaterial(albedo t.RGB) DiffuseMaterial {
	return DiffuseMaterial{
		Albedo: albedo,
	}
}

// -
// Scattering function for a diffuse material
// -
func (m DiffuseMaterial) scatter(r Ray, hit Hit) (bool, Ray, t.RGB) {
	// Scatter in a random direction in unit sphere around the normal
	scatterDir := hit.Normal.AddNew(t.RandVecSphere(true))
	if scatterDir.IsNearZero() {
		scatterDir = hit.Normal
	}

	scatterRay := NewRay(hit.Pos, scatterDir)

	return true, scatterRay, m.Albedo
}

func (m DiffuseMaterial) emitted() t.RGB {
	return t.Black()
}

// ============================================================
// Metal like material
// ============================================================

type MetalMaterial struct {
	Albedo t.RGB
	Fuzz   float64 // Brushed or fuzzy look of the metal
}

// -
// Create a new MetalMaterial with given albedo color and fuzziness
// -
func NewMetalMaterial(albedo t.RGB, fuzz float64) MetalMaterial {
	return MetalMaterial{
		Albedo: albedo,
		Fuzz:   math.Max(0, math.Min(fuzz, 1)),
	}
}

func (m MetalMaterial) scatter(r Ray, hit Hit) (bool, Ray, t.RGB) {
	// Metal is reflective
	scatterDir := r.Dir.Reflect(hit.Normal)
	scatterDir.Normalize()

	// Add some randomness to the reflected ray
	fuzz := t.RandVecSphere(false)
	fuzz.MultScalar(m.Fuzz)
	scatterDir.Add(fuzz)

	scatterRay := NewRay(hit.Pos, scatterDir)

	didScatter := scatterRay.Dir.Dot(hit.Normal) > 0
	return didScatter, scatterRay, m.Albedo
}

func (m MetalMaterial) emitted() t.RGB {
	return t.Black()
}
