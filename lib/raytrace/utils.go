package raytrace

import (
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
)

// -
// Generate a unique ID from a string
// -
func GenerateID(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	id := fmt.Sprintf("%x", hash.Sum(nil))
	return id[:6]
}

// ============================================================
// Simple interval between two numbers
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

func (i Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}

	return x
}

func (i Interval) String() string {
	minStr := strconv.FormatFloat(i.Min, 'f', -1, 64)
	maxStr := strconv.FormatFloat(i.Max, 'f', -1, 64)
	return "Interval{" + minStr + ", " + maxStr + "}"
}

// -
// Schlick's approximation for reflectance
// -
func reflectance(cosine float64, ior float64) float64 {
	r0 := (1 - ior) / (1 + ior)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
