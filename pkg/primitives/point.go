package primitives

import (
	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

type Point struct {
	Coordinate Vector
	Color      Vector
}

func (p *Point) SetColor(color Vector) {
	p.Color = color
}

type PointBB struct {
	Coordinate V.Vector
	Color      V.Vector
}

// SetColor updates the Color of the point
func (p *PointBB) SetColor(color V.Vector) {
	p.Color = color
}
