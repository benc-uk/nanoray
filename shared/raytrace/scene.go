package raytrace

import "image"

type Scene struct {
	Id            string
	Objects       []Hitable
	Lights        []Light
	ImageTextures []*image.RGBA
}

func NewScene() *Scene {
	return &Scene{
		Objects:       []Hitable{},
		Lights:        []Light{},
		ImageTextures: []*image.RGBA{},
	}
}

func (s *Scene) AddObject(o Hitable) {
	s.Objects = append(s.Objects, o)
}
