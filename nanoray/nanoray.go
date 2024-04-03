package main

import (
	"image/png"
	"log"
	"nanoray/shared/raytrace"
	t "nanoray/shared/tuples"
	"os"
)

func main() {
	aspect := 4.0 / 3.0
	r := raytrace.NewRender(1280, aspect)
	c := raytrace.NewCamera(t.Zero(), t.Zero(), 60)
	s := raytrace.Scene{}

	sphere := raytrace.NewSphere(t.Vec3{0, 0, -90}, 30)
	s.AddObject(sphere)

	sphere2 := raytrace.NewSphere(t.Vec3{0, -528, -90}, 500)
	s.AddObject(sphere2)

	// s.Lights = append(s.Lights, raytrace.Light{
	// 	Position: v.Vec3{170, -100, 3},
	// 	Color:    v.White(),
	// })

	r.Generate(c, s)

	img := r.Image

	f, err := os.Create("output.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
