package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Point struct {
	X, Y, Z float64
}

func main() {
	width := 1280
	height := 1720
	widthF := 1280.0
	heightF := 1720.0

	// Create a new image with the specified width and height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Define the sphere properties
	spheres := []Point{}

	// add 30 spheres to the scene in random positions
	for i := 0; i < 1000; i++ {
		posX := rand.Float64()*70 - 35
		posY := rand.Float64()*40 - 20
		spheres = append(spheres, Point{X: posX, Y: posY, Z: 40})
	}

	// sphereCenter := Point{X: 0, Y: 0, Z: 25}
	sphereRadius := 5
	sphereColor := color.RGBA{R: 255, G: 70, B: 0, A: 255}
	camPos := Point{X: 0, Y: 0, Z: -20}

	lightDir := Point{X: 1, Y: -2, Z: -1}
	lightDirLength := math.Sqrt(lightDir.X*lightDir.X + lightDir.Y*lightDir.Y + lightDir.Z*lightDir.Z)
	lightDir = Point{X: lightDir.X / lightDirLength, Y: lightDir.Y / lightDirLength, Z: lightDir.Z / lightDirLength}

	now := time.Now()

	// Iterate over each pixel in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			ix := float64(x) - widthF/2.0
			iy := float64(y) - heightF/2.0

			// Calculate the ray direction
			rayDir := Point{X: ix / heightF, Y: iy / heightF, Z: 1}
			rayDirLength := math.Sqrt(rayDir.X*rayDir.X + rayDir.Y*rayDir.Y + rayDir.Z*rayDir.Z)
			rayDir = Point{X: rayDir.X / rayDirLength, Y: rayDir.Y / rayDirLength, Z: rayDir.Z / rayDirLength}

			hit := false
			for i := 0; i < len(spheres); i++ {
				sphereCenter := spheres[i]

				// Calculate the intersection of the ray with the sphere
				oc := Point{X: camPos.X - sphereCenter.X, Y: camPos.Y - sphereCenter.Y, Z: camPos.Z - sphereCenter.Z}
				a := rayDir.X*rayDir.X + rayDir.Y*rayDir.Y + rayDir.Z*rayDir.Z
				b := 2 * (oc.X*rayDir.X + oc.Y*rayDir.Y + oc.Z*rayDir.Z)
				c := oc.X*oc.X + oc.Y*oc.Y + oc.Z*oc.Z - float64(sphereRadius*sphereRadius)
				discriminant := b*b - 4*a*c

				// If the discriminant is positive, there are two intersection points
				if discriminant > 0 {
					t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
					t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

					if t1 < 0 && t2 < 0 {
						img.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
						continue
					}

					// // Choose the closest intersection point
					// t := t1
					// if t2 < t1 {
					// 	t = t2
					// }

					// // Calculate the intersection point
					// intersection := Point{X: camPos.X + rayDir.X*t, Y: camPos.Y + rayDir.Y*t, Z: camPos.Z + rayDir.Z*t}

					// // Calculate the normal at the intersection point
					// normal := Point{X: intersection.X - sphereCenter.X, Y: intersection.Y - sphereCenter.Y, Z: intersection.Z - sphereCenter.Z}
					// normalLength := math.Sqrt(normal.X*normal.X + normal.Y*normal.Y + normal.Z*normal.Z)
					// normal = Point{X: normal.X / normalLength, Y: normal.Y / normalLength, Z: normal.Z / normalLength}

					// // Calculate the shading intensity
					// shade := normal.X*lightDir.X + normal.Y*lightDir.Y + normal.Z*lightDir.Z
					// if shade < 0 {
					// 	shade = 0
					// }
					shade := 1.0

					img.SetRGBA(x, y, color.RGBA{R: uint8(float64(sphereColor.R) * shade), G: uint8(float64(sphereColor.G) * shade), B: uint8(float64(sphereColor.B) * shade), A: 255})
					hit = true
					//break
				}
			}
			if !hit {
				img.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	// Print the time taken to render the image
	log.Printf("Time taken: %v", time.Since(now))

	// Save the image to a PNG file
	outputFile, err := os.Create("output2.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
