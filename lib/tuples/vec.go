package tuples

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v *Vec3) UnmarshalYAML(unmarshal func(in any) error) error {
	var data []float64
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	if len(data) != 3 {
		return fmt.Errorf("cannot unmarshal Vec3 from %v", data)
	}

	v.X = data[0]
	v.Y = data[1]
	v.Z = data[2]

	return nil
}

func Zero() Vec3 {
	return Vec3{0, 0, 0}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3) Add(v2 Vec3) {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
}

func (v Vec3) AddNew(v2 Vec3) Vec3 {
	return Vec3{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v *Vec3) Sub(v2 Vec3) {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
}

func (v Vec3) SubNew(v2 Vec3) Vec3 {
	return Vec3{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v *Vec3) Mult(s float64) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

func (v Vec3) MultNew(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v *Vec3) Div(s float64) {
	v.X /= s
	v.Y /= s
	v.Z /= s
}

func (v Vec3) DivNew(s float64) Vec3 {
	return Vec3{v.X / s, v.Y / s, v.Z / s}
}

func (v *Vec3) MultScalar(s float64) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

func (v Vec3) MultScalarNew(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) Dot(v2 Vec3) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		v.Y*v2.Z - v.Z*v2.Y,
		v.Z*v2.X - v.X*v2.Z,
		v.X*v2.Y - v.Y*v2.X,
	}
}

func (v *Vec3) Normalize() {
	v.Div(v.Length())
}

func (v Vec3) NormalizeNew() Vec3 {
	len := v.Length()
	return Vec3{v.X / len, v.Y / len, v.Z / len}
}

func (v Vec3) Sqrt() Vec3 {
	return Vec3{math.Sqrt(v.X), math.Sqrt(v.Y), math.Sqrt(v.Z)}
}

func (v *Vec3) Negate() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

func (v Vec3) NegateNew() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Reflect(normal Vec3) Vec3 {
	return v.SubNew(normal.MultScalarNew(2 * v.Dot(normal)))
}

func (v Vec3) Refract(normal Vec3, etaiOverEtat float64) Vec3 {
	// Code from https://raytracing.github.io/books/RayTracingInOneWeekend.html#dielectrics

	cosTheta := math.Min(v.NegateNew().Dot(normal), 1.0)

	// vec3 r_out_perp =  etai_over_etat * (uv + cos_theta*n);
	rOutPerp := normal.MultScalarNew(cosTheta).AddNew(v).MultScalarNew(etaiOverEtat)

	// vec3 r_out_parallel = -sqrt(fabs(1.0 - r_out_perp.length_squared())) * n;
	rOutPar := normal.MultScalarNew(-math.Sqrt(math.Abs(1.0 - rOutPerp.SquaredLength())))

	// return r_out_perp + r_out_parallel;
	return rOutPerp.AddNew(rOutPar)
}

func (v Vec3) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

func (v Vec3) IsNearZero() bool {
	const s = 1e-8
	return math.Abs(v.X) < s && math.Abs(v.Y) < s && math.Abs(v.Z) < s
}

func (v Vec3) Equals(v2 Vec3) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

func (v Vec3) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", v.X, v.Y, v.Z)
}

func (v Vec3) Clone() Vec3 {
	return Vec3{v.X, v.Y, v.Z}
}

func ParseVec3(data any) (Vec3, error) {
	tuple, ok := data.([]any)
	if !ok {
		return Zero(), fmt.Errorf("ParseVec3: Failed to convert data to array: %+v", data)
	}

	return Vec3{tuple[0].(float64), tuple[1].(float64), tuple[2].(float64)}, nil
}

// ============================================================
// Random Vec3 functions for path tracing
// ============================================================

// -
// Generate a random vector inside the unit cube
// -
func RandVecCube() Vec3 {
	return Vec3{
		rand.Float64()*2 - 1,
		rand.Float64()*2 - 1,
		rand.Float64()*2 - 1,
	}
}

// -
// Generate a random vector inside the unit sphere
// -
func RandVecSphere(normalise bool) Vec3 {
	var v Vec3
	// Keep generating random vectors until we get one inside the unit sphere
	for {
		v = RandVecCube()
		if v.SquaredLength() < 1 {
			break
		}
	}

	if normalise {
		v.Normalize()
	}

	return v
}

// -
// Generate a random vector inside the unit disk
// -
func RandVecDisk(normalise bool) Vec3 {
	var v Vec3

	for {
		v = Vec3{rand.Float64()*2 - 1, rand.Float64()*2 - 1, 0}
		if v.SquaredLength() < 1 {
			break
		}
	}

	if normalise {
		v.Normalize()
	}

	return v
}

func RandVecSphereHemisphere(normal Vec3) Vec3 {
	v := RandVecSphere(true)
	if v.Dot(normal) > 0 {
		return v
	} else {
		return v.NegateNew()
	}
}
