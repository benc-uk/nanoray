package tuples

import (
	"fmt"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func Zero() Vec3 {
	return Vec3{0, 0, 0}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
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

func (v Vec3) MultScalar(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) MultScalarNew(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v1 Vec3) Dot(v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
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

func (v Vec3) String() string {
	return fmt.Sprintf("[%f, %f, %f]", v.X, v.Y, v.Z)
}
