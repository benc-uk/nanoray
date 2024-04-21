package tuples

import (
	"fmt"
	"image/color"
	"math"
)

type RGB struct {
	R, G, B float64
}

func (c *RGB) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data []float64
	err := unmarshal(&data)
	if err != nil {
		return err
	}

	if len(data) != 3 {
		return fmt.Errorf("cannot unmarshal RGB from %v", data)
	}

	c.R = data[0]
	c.G = data[1]
	c.B = data[2]

	return nil
}

func FromHexString(hex string) RGB {
	var r, g, b uint8
	fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return RGB{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

func From8Bit(r, g, b uint8) RGB {
	return RGB{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

func (c RGB) ToHexString() string {
	return fmt.Sprintf("#%02x%02x%02x", uint8(c.R*255), uint8(c.G*255), uint8(c.B*255))
}

func (c *RGB) Clamp() {
	c.R = math.Min(1, math.Max(0, c.R))
	c.G = math.Min(1, math.Max(0, c.G))
	c.B = math.Min(1, math.Max(0, c.B))
}

func (c *RGB) Add(c2 RGB) {
	c.R += c2.R
	c.G += c2.G
	c.B += c2.B
}

func (c *RGB) AddSome(c2 RGB, t float64) {
	c.R += c2.R * t
	c.G += c2.G * t
	c.B += c2.B * t
}

func (c RGB) AddNew(c2 RGB) RGB {
	return RGB{c.R + c2.R, c.G + c2.G, c.B + c2.B}
}

func (c RGB) SubNew(c2 RGB) RGB {
	return RGB{c.R - c2.R, c.G - c2.G, c.B - c2.B}
}

func (c RGB) Blend(c2 RGB, t float64) RGB {
	return RGB{
		c.R*(1-t) + c2.R*t,
		c.G*(1-t) + c2.G*t,
		c.B*(1-t) + c2.B*t,
	}
}

func (c *RGB) MultScalar(s float64) {
	c.R *= s
	c.G *= s
	c.B *= s
}

func (c RGB) MultScalarNew(s float64) RGB {
	return RGB{c.R * s, c.G * s, c.B * s}
}

func (c *RGB) Mult(c2 RGB) {
	c.R *= c2.R
	c.G *= c2.G
	c.B *= c2.B
}

func (c RGB) MultNew(c2 RGB) RGB {
	return RGB{c.R * c2.R, c.G * c2.G, c.B * c2.B}
}

// -
// Convert to a color.RGBA object with gamma correction
// -
func (c RGB) ToRGBA(gamma float64) color.RGBA {
	var r, g, b float64

	// Gamma correction
	r = math.Pow(c.R, 1/gamma)
	g = math.Pow(c.G, 1/gamma)
	b = math.Pow(c.B, 1/gamma)

	r = math.Min(0.999, math.Max(0, r))
	g = math.Min(0.999, math.Max(0, g))
	b = math.Min(0.999, math.Max(0, b))

	return color.RGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		255,
	}
}

func (c RGB) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", c.R, c.G, c.B)
}

func ParseRGB(data any) (RGB, error) {
	if data == nil {
		return Black(), fmt.Errorf("ParseRGB: Data is nil")
	}

	tuple, ok := data.([]any)
	if !ok {
		return Black(), fmt.Errorf("ParseRGB: Failed to convert data to array: %+v", data)
	}

	// check if the array has 3 elements
	if len(tuple) != 3 {
		return Black(), fmt.Errorf("ParseRGB: Must have 3 elements, but got: %d", len(tuple))
	}

	// convert array elements to float64
	for i := 0; i < 3; i++ {
		switch v := tuple[i].(type) {
		case int:
			tuple[i] = float64(v)
		case float64:
			// do nothing
		default:
			return Black(), fmt.Errorf("ParseRGB: Invalid data type: %T %v", tuple[i], tuple[i])
		}
	}

	return RGB{tuple[0].(float64), tuple[1].(float64), tuple[2].(float64)}, nil
}

func (c RGB) Equals(c2 RGB) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B
}

// ============================================================================
// Predefined colors
// ============================================================================

func Black() RGB {
	return RGB{0, 0, 0}
}

func Red() RGB {
	return RGB{1, 0, 0}
}

func Green() RGB {
	return RGB{0, 1, 0}
}

func Blue() RGB {
	return RGB{0, 0, 1}
}

func White() RGB {
	return RGB{1, 1, 1}
}
