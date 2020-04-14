package primitives

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Point struct {
	Coordinate mgl32.Vec3
	Color      mgl32.Vec3
}

// SetColor updates the Color of the point.
func (p *Point) SetColor(color mgl32.Vec3) {
	p.Color = color
}
