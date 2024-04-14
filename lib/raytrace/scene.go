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
	Camera  SceneFileCamera   `yaml:"camera"`
	Objects []SceneFileObject `yaml:"objects"`
}

type SceneFileObject struct {
	Type     string  `yaml:"type"`
	Position t.Vec3  `yaml:"position"`
	Radius   float64 `yaml:"radius"`
	Colour   t.RGB   `yaml:"colour"`
}

type SceneFileCamera struct {
	Position t.Vec3  `yaml:"position"`
	LookAt   t.Vec3  `yaml:"lookAt"`
	Fov      float64 `yaml:"fov"`
}

func ParseScene(sceneData string, imgW, imgH int) (*Scene, *Camera, error) {
	log.Printf("Parsing scene data: %d bytes", len(sceneData))

	var sceneFile SceneFile
	err := yaml.Unmarshal([]byte(sceneData), &sceneFile)
	if err != nil {
		return nil, nil, err
	}

	if sceneFile.Camera.Fov == 0 {
		log.Printf("No FOV specified, defaulting to 50")
		sceneFile.Camera.Fov = 50
	}

	if sceneFile.Camera.Position.Equals(sceneFile.Camera.LookAt) {
		log.Printf("Camera position and lookAt are the same, this is probably not what you want")
	}

	camera := NewCamera(imgW, imgH, sceneFile.Camera.Position, sceneFile.Camera.LookAt, sceneFile.Camera.Fov)

	scene := &Scene{
		Name:    sceneFile.Name,
		Objects: []Hitable{},
	}

	for _, obj := range sceneFile.Objects {
		switch obj.Type {
		case "sphere":
			sphere, err := NewSphere(obj.Position, obj.Radius)
			if err != nil {
				log.Printf("Failed to create sphere: %s", err.Error())
				continue
			}
			sphere.Colour = obj.Colour
			scene.AddObject(sphere)

			log.Printf("Added sphere at %v with radius %.1f", obj.Position, obj.Radius)
			continue
		}

		log.Printf("Skipping unknown object type: %s", obj.Type)
	}

	return scene, &camera, nil
}

func (s *Scene) AddObject(o Hitable) {
	s.Objects = append(s.Objects, o)
}
