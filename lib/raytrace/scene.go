package raytrace

import (
	"log"
	t "nanoray/lib/tuples"

	"gopkg.in/yaml.v3"
)

type Scene struct {
	Name    string
	Objects []Hitable
}

type SceneFile struct {
	Name    string            `yaml:"name"`
	Objects []SceneFileObject `yaml:"objects"`
}

type SceneFileObject struct {
	Type     string  `yaml:"type"`
	Position t.Vec3  `yaml:"position"`
	Radius   float64 `yaml:"radius"`
	Colour   t.RGB   `yaml:"colour"`
}

func ParseScene(sceneData string) (*Scene, error) {
	log.Printf("Parsing scene data: %d bytes", len(sceneData))

	var sceneFile SceneFile
	err := yaml.Unmarshal([]byte(sceneData), &sceneFile)
	if err != nil {
		return nil, err
	}

	scene := &Scene{
		Name:    sceneFile.Name,
		Objects: []Hitable{},
	}

	for _, obj := range sceneFile.Objects {
		switch obj.Type {
		case "sphere":
			sphere := NewSphere(obj.Position, obj.Radius)
			sphere.Colour = obj.Colour
			scene.AddObject(sphere)
			log.Printf("Added sphere at %v with radius %f", obj.Position, obj.Radius)
		}
	}

	return scene, nil
}

func (s *Scene) AddObject(o Hitable) {
	s.Objects = append(s.Objects, o)
}
