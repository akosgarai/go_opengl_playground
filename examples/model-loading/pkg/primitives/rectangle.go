package primitives

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Rectangle struct {
	Points   [4]mgl32.Vec3
	Normal   mgl32.Vec3
	Indicies []mgl32.Vec3
}

// NewSquare creates a rectangle with origo as middle point.
// The normal points to +Y.
// The longest side is scaled to one, and the same downscale is done with the other edge.
// - width represents the lenght on the X axis.
// - height represents the lenght on the Z axis.
// ratio = width / length
// ratio == 1 => return NewSquare
// ratio > 1 => width is the longer -> X [-0.5, 0.5], Y [-1/(ratio*2), 1/(ratio*2)]
// ratio < 1 => length is the longer
func New(width, length float32) *Rectangle {
	normal := mgl32.Vec3{0, -1, 0}
	ratio := width / length
	var x0, x1, y0, y1 float32
	if ratio == 1 {
		return NewSquare()
	} else if ratio > 1 {
		x0 = float32(-0.5)
		x1 = float32(0.5)
		y0 = float32(-1 / (ratio * 2))
		y1 = float32(1 / (ratio * 2))
	} else {
		y0 = float32(-0.5)
		y1 = float32(0.5)
		x0 = float32(-ratio / 2)
		x1 = float32(ratio / 2)
	}
	points := [4]mgl32.Vec3{
		mgl32.Vec3{x0, 0, y0},
		mgl32.Vec3{x1, 0, y0},
		mgl32.Vec3{x1, 0, y1},
		mgl32.Vec3{x0, 0, y1},
	}
	return &Rectangle{
		Points: points,
		Normal: normal,
	}
}

// NewSquare creates a rectangle with origo as middle point.
// Each side is 1 unit long, and it's plane is the X-Z plane.
// The normal points to +Y.
func NewSquare() *Rectangle {
	normal := mgl32.Vec3{0, -1, 0}
	points := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	return &Rectangle{
		Points: points,
		Normal: normal,
	}
}
