package raytrace

import (
	"log"
	t "nanoray/lib/tuples"

	"gopkg.in/yaml.v3"
)

type Scene struct {
	Name       string
	Background t.RGB
	Objects    []Hitable
}

type File struct {
	Name       string       `yaml:"name"`
	Background t.RGB        `yaml:"background"`
	Camera     FileCamera   `yaml:"camera"`
	Objects    []FileObject `yaml:"objects"`
}

type FileObject struct {
	Type     string         `yaml:"type"`
	Position t.Vec3         `yaml:"position"`
	Radius   float64        `yaml:"radius"`
	Material map[string]any `yaml:"material"`
}

type FileMaterial struct {
	Dielectric FileDielectricMat `yaml:"dielectric"`
	Diffuse    FileDiffuseMat    `yaml:"diffuse"`
	Metal      FileMetalMat      `yaml:"metal"`
}

type FileDiffuseMat struct {
	Albedo t.RGB `yaml:"albedo"`
}

type FileMetalMat struct {
	Albedo t.RGB   `yaml:"albedo"`
	Fuzz   float64 `yaml:"fuzz"`
}

type FileDielectricMat struct {
	Tint t.RGB   `yaml:"tint"`
	Fuzz float64 `yaml:"fuzz"`
	IOR  float64 `yaml:"ior"`
}

type FileCamera struct {
	Position  t.Vec3  `yaml:"position"`
	LookAt    t.Vec3  `yaml:"lookAt"`
	Fov       float64 `yaml:"fov"`
	FocalDist float64 `yaml:"focalDist"`
	Aperture  float64 `yaml:"aperture"`
}

// -
// Parse a scene & camera from a YAML string
// -
func ParseScene(sceneData string, imgW, imgH int) (*Scene, *Camera, error) {
	log.Printf("Parsing scene data: %d bytes", len(sceneData))

	var File File
	err := yaml.Unmarshal([]byte(sceneData), &File)
	if err != nil {
		return nil, nil, err
	}

	if File.Camera.Fov == 0 {
		log.Printf("No FOV specified, defaulting to 50")
		File.Camera.Fov = 50
	}

	if File.Camera.Position.Equals(File.Camera.LookAt) {
		log.Printf("Camera position and lookAt are the same, this is probably not what you want")
	}

	camera := NewCamera(imgW, imgH, File.Camera.Position,
		File.Camera.LookAt, File.Camera.Fov, File.Camera.FocalDist, File.Camera.Aperture)

	scene := &Scene{
		Name:       File.Name,
		Objects:    []Hitable{},
		Background: File.Background,
	}

	for _, obj := range File.Objects {
		switch obj.Type {
		case "sphere":
			worldObj, err := NewSphere(obj.Position, obj.Radius)
			if err != nil {
				log.Printf("Failed to create sphere: %s", err.Error())
				continue
			}

			m := parseMaterial(obj.Material)
			if m != nil {
				log.Printf("Added sphere at %v with radius %.1f, material type: %s", obj.Position, obj.Radius, m.Type())
				worldObj.Material = m
				scene.AddObject(worldObj)
			}

		default:
			log.Printf("Unknown object type: %s", obj.Type)
		}
	}

	return scene, &camera, nil
}

func parseMaterial(material map[string]any) Material {
	if material == nil {
		return nil
	}

	if props, ok := material["dielectric"]; ok {
		if props == nil {
			log.Printf("Warning, dielectric material was empty")
			return DielectricMaterial{}
		}

		propMap := props.(map[string]any)

		m := DielectricMaterial{}
		m.Tint = t.RGB{0.95, 0.95, 0.95}
		if propMap["tint"] != nil {
			var err error
			m.Tint, err = t.ParseRGB(propMap["tint"])
			if err != nil {
				log.Printf("Failed to parse tint: %s", err.Error())
			}
		}

		m.Fuzz = parseFloatOrInt(propMap["fuzz"])
		m.IOR = parseFloatOrInt(propMap["ior"])
		return &m
	}

	if props, ok := material["diffuse"]; ok {
		if props == nil {
			log.Printf("Warning, diffuse material was empty")
			return DiffuseMaterial{}
		}

		propMap := props.(map[string]any)

		m := DiffuseMaterial{}

		var err error
		m.Albedo, err = t.ParseRGB(propMap["albedo"])
		if err != nil {
			log.Printf("Failed to parse albedo: %s", err.Error())
		}

		return &m
	}

	if props, ok := material["metal"]; ok {
		if props == nil {
			log.Printf("Warning, metal material was empty")
			return MetalMaterial{}
		}

		propMap := props.(map[string]any)

		m := MetalMaterial{}
		var err error
		m.Albedo, err = t.ParseRGB(propMap["albedo"])
		if err != nil {
			log.Printf("Failed to parse albedo: %s", err.Error())
		}
		m.Fuzz = parseFloatOrInt(propMap["fuzz"])
		return &m
	}

	if props, ok := material["light"]; ok {
		if props == nil {
			log.Printf("Warning, light material was empty")
			return LightMaterial{}
		}

		propMap := props.(map[string]any)

		m := LightMaterial{}
		var err error
		m.Emission, err = t.ParseRGB(propMap["emission"])
		if err != nil {
			log.Printf("Failed to parse emission: %s", err.Error())
		}
		return &m
	}

	log.Printf("Unknown material type: %v", material)
	return nil
}

func parseFloatOrInt(data any) float64 {
	if data == nil {
		return 0
	}

	switch v := data.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		log.Printf("Failed to convert %v to float64", data)
		return 0
	}
}

// -
// Add an object to the scene
// -
func (s *Scene) AddObject(o Hitable) {
	s.Objects = append(s.Objects, o)
}
