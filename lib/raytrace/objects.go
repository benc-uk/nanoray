package raytrace

import (
	t "nanoray/lib/tuples"
)

// All objects should embed this struct
type Object struct {
	ID       string
	Position t.Vec3
	Material Material
}

// All objects must implement this interface
type Hitable interface {
	Hit(r Ray, i Interval) (bool, Hit)
}

// Hit represents a ray hit against an object
type Hit struct {
	T      float64 // Distance along ray
	Pos    t.Vec3  // Position of hit in world space
	Normal t.Vec3  // Normal at hit point
	Obj    Object  // Ref to object that was hit
	Front  bool    // Is the hit on the front/outside of the object
}

func (h Hit) String() string {
	return "Hit{" + h.Pos.String() + ", " + h.Normal.String() + "}"
}
