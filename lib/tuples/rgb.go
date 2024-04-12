package tuples

import (
	"image/color"
	"math"
)

type RGB struct {
	R, G, B float64
}

func (v *RGB) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp [3]float64

	err := unmarshal(&tmp)
	if err != nil {
		return err
	}

	v.R = tmp[0]
	v.G = tmp[1]
	v.B = tmp[2]

	return nil
}

func (c *RGB) Clamp() {
	c.R = math.Min(1, math.Max(0, c.R))
	c.G = math.Min(1, math.Max(0, c.G))
	c.B = math.Min(1, math.Max(0, c.B))
}

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

func (c RGB) ToRGBA() color.RGBA {
	return color.RGBA{
		uint8(math.Min(255, math.Max(0, c.R*255))),
		uint8(math.Min(255, math.Max(0, c.G*255))),
		uint8(math.Min(255, math.Max(0, c.B*255))),
		255,
	}
}
