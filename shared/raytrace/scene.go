package raytrace

import "image"

type Scene struct {
	Id            string
	Objects       []Hitable
	ImageTextures []*image.RGBA
	MaxDepth      int
}

func NewScene() *Scene {
	return &Scene{
		Objects:       []Hitable{},
		ImageTextures: []*image.RGBA{},
	}
}

func (s *Scene) AddObject(o Hitable) {
	s.Objects = append(s.Objects, o)
}
