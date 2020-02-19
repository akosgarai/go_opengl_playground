package point

import (
	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

type Point struct {
	Coordinate V.Vector
	Color      V.Vector
}

// SetColor updates the Color of the point
func (p *Point) SetColor(color V.Vector) {
	p.Color = color
}
