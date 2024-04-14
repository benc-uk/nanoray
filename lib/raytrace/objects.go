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

type Hit struct {
	T      float64
	Pos    t.Vec3
	Normal t.Vec3
	Obj    Object
	Front  bool
}

// All objects should implement this interface
type Hitable interface {
	Hit(r Ray, i Interval) (bool, Hit)
}

func (h Hit) String() string {
	return "Hit{" + h.Pos.String() + ", " + h.Normal.String() + "}"
}
