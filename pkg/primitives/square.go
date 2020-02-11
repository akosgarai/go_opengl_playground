package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Square struct {
	A Point
	B Point
	C Point
	D Point
}

func NewSquare(v1, v2, v3, v4 Vector) *Square {
	return &Square{Point{v1, Vector{1, 1, 1}}, Point{v2, Vector{1, 1, 1}}, Point{v3, Vector{1, 1, 1}}, Point{v4, Vector{1, 1, 1}}}
}

func (s *Square) SetColor(color Vector) {
	s.A.SetColor(color)
	s.B.SetColor(color)
	s.C.SetColor(color)
	s.D.SetColor(color)
}

func (s *Square) buildVao() uint32 {
	var points []float32
	// Coordinates
	points = append(points, float32(s.A.Coordinate.X))
	points = append(points, float32(s.A.Coordinate.Y))
	points = append(points, float32(s.A.Coordinate.Z))

	points = append(points, float32(s.B.Coordinate.X))
	points = append(points, float32(s.B.Coordinate.Y))
	points = append(points, float32(s.B.Coordinate.Z))

	points = append(points, float32(s.C.Coordinate.X))
	points = append(points, float32(s.C.Coordinate.Y))
	points = append(points, float32(s.C.Coordinate.Z))

	points = append(points, float32(s.C.Coordinate.X))
	points = append(points, float32(s.C.Coordinate.Y))
	points = append(points, float32(s.C.Coordinate.Z))

	points = append(points, float32(s.D.Coordinate.X))
	points = append(points, float32(s.D.Coordinate.Y))
	points = append(points, float32(s.D.Coordinate.Z))

	points = append(points, float32(s.A.Coordinate.X))
	points = append(points, float32(s.A.Coordinate.Y))
	points = append(points, float32(s.A.Coordinate.Z))

	// Colors
	points = append(points, float32(s.A.Color.X))
	points = append(points, float32(s.A.Color.Y))
	points = append(points, float32(s.A.Color.Z))

	points = append(points, float32(s.B.Color.X))
	points = append(points, float32(s.B.Color.Y))
	points = append(points, float32(s.B.Color.Z))

	points = append(points, float32(s.C.Color.X))
	points = append(points, float32(s.C.Color.Y))
	points = append(points, float32(s.C.Color.Z))

	points = append(points, float32(s.C.Color.X))
	points = append(points, float32(s.C.Color.Y))
	points = append(points, float32(s.C.Color.Z))

	points = append(points, float32(s.D.Color.X))
	points = append(points, float32(s.D.Color.Y))
	points = append(points, float32(s.D.Color.Z))

	points = append(points, float32(s.A.Color.X))
	points = append(points, float32(s.A.Color.Y))
	points = append(points, float32(s.A.Color.Z))

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(4*3*6))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	return vertexArrayObject
}

func (s *Square) Draw() {
	vertexArrayObject := s.buildVao()
	gl.BindVertexArray(vertexArrayObject)
	// The square is represented by 2 triangle, so we have 2 * 3 points here.
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
