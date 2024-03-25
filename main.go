package main

import (
	"fmt"
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

type Colour struct {
	R, G, B, A float64
}

type Pixel struct {
	X, Y   int
	Colour Colour
}

// Define the global properties
var spheres = []Point{}
var camPos = Point{X: 0, Y: 0, Z: -50}
var sphereRadius = 5
var sphereColour = Colour{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
var lightPos = Point{X: 1211, Y: -310, Z: -100}

func main() {
	width := 1280
	height := 720
	widthF := 1280.0
	heightF := 720.0

	// Create a new image with the specified width and height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// add 30 spheres to the scene in random positions
	for i := 0; i < 5000; i++ {
		posX := rand.Float64()*70 - 35
		posY := rand.Float64()*40 - 20
		posZ := (rand.Float64()*40 - 20) + 40
		spheres = append(spheres, Point{X: posX, Y: posY, Z: posZ})
	}

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

			hitSphere := -1
			minT := math.MaxFloat64

			for i := 0; i < len(spheres); i++ {
				sphereCenter := spheres[i]
				t := math.MaxFloat64

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

					if t1 > 0 && t2 > 0 {
						t = math.Min(t1, t2)
					}
				}

				if t < minT {
					minT = t
					hitSphere = i
				}
			}

			if hitSphere != -1 {
				sphereCenter := spheres[hitSphere]
				// Calculate the intersection point
				intersection := Point{X: camPos.X + rayDir.X*minT, Y: camPos.Y + rayDir.Y*minT, Z: camPos.Z + rayDir.Z*minT}

				// Calculate the normal at the intersection point
				normal := Point{X: intersection.X - sphereCenter.X, Y: intersection.Y - sphereCenter.Y, Z: intersection.Z - sphereCenter.Z}
				normalLength := math.Sqrt(normal.X*normal.X + normal.Y*normal.Y + normal.Z*normal.Z)
				normal = Point{X: normal.X / normalLength, Y: normal.Y / normalLength, Z: normal.Z / normalLength}

				// Calculate the light direction
				lightDir := Point{X: lightPos.X - intersection.X, Y: lightPos.Y - intersection.Y, Z: lightPos.Z - intersection.Z}
				lightDirLength := math.Sqrt(lightDir.X*lightDir.X + lightDir.Y*lightDir.Y + lightDir.Z*lightDir.Z)
				lightDir = Point{X: lightDir.X / lightDirLength, Y: lightDir.Y / lightDirLength, Z: lightDir.Z / lightDirLength}

				// Calculate the shading intensity
				shade := normal.X*lightDir.X + normal.Y*lightDir.Y + normal.Z*lightDir.Z
				if shade < 0 {
					shade = 0
				}

				// calculate the specular highlight
				reflect := Point{X: 2 * normal.X * shade, Y: 2 * normal.Y * shade, Z: 2 * normal.Z * shade}
				reflect = Point{X: reflect.X - lightDir.X, Y: reflect.Y - lightDir.Y, Z: reflect.Z - lightDir.Z}
				// view vector
				view := Point{X: -rayDir.X, Y: -rayDir.Y, Z: -rayDir.Z}

				specular := reflect.X*view.X + reflect.Y*view.Y + reflect.Z*view.Z
				specular = math.Max(0, specular)
				specular = math.Pow(specular, 10)

				shadeColour := Colour{
					R: (sphereColour.R * shade) + specular,
					G: (sphereColour.G * shade) + specular,
					B: (sphereColour.B * shade) + specular,
					A: 1.0}

				// Clamp the colour values
				shadeColour.clamp()

				img.SetRGBA(x, y, shadeColour.toRGBA())
			}
		}
	}

	// Print the time taken to render the image
	fmt.Printf("Time taken: %dms\n", time.Since(now).Milliseconds())

	// Save the image to a PNG file
	outputFile, err := os.Create("output-go.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Colour) toRGBA() color.RGBA {
	return color.RGBA{R: uint8(c.R * 255), G: uint8(c.G * 255), B: uint8(c.B * 255), A: uint8(c.A * 255)}
}

func (c *Colour) clamp() {
	if c.R > 1.0 {
		c.R = 1.0
	}
	if c.G > 1.0 {
		c.G = 1.0
	}
	if c.B > 1.0 {
		c.B = 1.0
	}
	if c.A > 1.0 {
		c.A = 1.0
	}
}

func getSphereT(i int, rayDir Point) float64 {
	sphereCenter := spheres[i]
	t := math.MaxFloat64

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

		if t1 > 0 && t2 > 0 {
			t = math.Min(t1, t2)
		}
	}

	return t
}
