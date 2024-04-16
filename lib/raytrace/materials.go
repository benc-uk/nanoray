package raytrace

import (
	"math"
	"math/rand"
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

// ============================================================
// Dielectric which transmits light e.g. glass, water etc
// ============================================================

type DielectricMaterial struct {
	IOR  float64
	Fuzz float64
	Tint t.RGB
}

// -
// Create a new DielectricMaterial with given index of refraction
// -
func NewDielectricMaterial(ior float64, fuzz float64, tint t.RGB) DielectricMaterial {
	if fuzz < 0 {
		fuzz = 0
	}

	return DielectricMaterial{
		IOR:  ior,
		Fuzz: fuzz,
		Tint: tint,
	}
}

func (m DielectricMaterial) scatter(r Ray, hit Hit) (bool, Ray, t.RGB) {
	attenuation := m.Tint
	ri := m.IOR
	if hit.Front {
		ri = 1.0 / m.IOR
	}

	cosTheta := math.Min(r.Dir.NegateNew().Dot(hit.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	reflect := ri*sinTheta > 1.0

	var scatterDir t.Vec3
	if reflect || reflectance(cosTheta, ri) > rand.Float64() {
		// Reflect the ray
		scatterDir = r.Dir.Reflect(hit.Normal)
	} else {
		// Refract the ray
		scatterDir = r.Dir.Refract(hit.Normal, ri)
		// Allow for frosted glass effect
		fuzz := t.RandVecSphere(false)
		fuzz.MultScalar(m.Fuzz)
		scatterDir.Add(fuzz)
	}

	scatterRay := NewRay(hit.Pos, scatterDir)
	return true, scatterRay, attenuation
}

func (m DielectricMaterial) emitted() t.RGB {
	return t.Black()
}
