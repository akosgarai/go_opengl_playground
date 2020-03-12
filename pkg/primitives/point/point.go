package point

import (
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

type Point struct {
	Coordinate vec.Vector
	Color      vec.Vector
}

// SetColor updates the Color of the point.
func (p *Point) SetColor(color vec.Vector) {
	p.Color = color
}
