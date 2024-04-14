package raytrace

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

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
