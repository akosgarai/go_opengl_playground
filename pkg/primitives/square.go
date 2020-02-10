package primitives

import (
	//vao "github.com/akosgarai/opengl_playground/pkg/vertexarrayobject"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Square struct {
	A     Vector
	B     Vector
	C     Vector
	D     Vector
	Color Vector
}

func NewSquare(v1, v2, v3, v4 Vector) *Square {
	return &Square{v1, v2, v3, v4, Vector{1, 1, 1}}
}

func (s *Square) SetColor(color Vector) {
	s.Color = color
}

func (s *Square) buildVao() uint32 {
	var points []float32
	// Coordinates
	points = append(points, float32(s.A.X))
	points = append(points, float32(s.A.Y))
	points = append(points, float32(s.A.Z))

	points = append(points, float32(s.B.X))
	points = append(points, float32(s.B.Y))
	points = append(points, float32(s.B.Z))

	points = append(points, float32(s.C.X))
	points = append(points, float32(s.C.Y))
	points = append(points, float32(s.C.Z))

	points = append(points, float32(s.C.X))
	points = append(points, float32(s.C.Y))
	points = append(points, float32(s.C.Z))

	points = append(points, float32(s.D.X))
	points = append(points, float32(s.D.Y))
	points = append(points, float32(s.D.Z))

	points = append(points, float32(s.A.X))
	points = append(points, float32(s.A.Y))
	points = append(points, float32(s.A.Z))

	// Colors (red)
	for _, _ = range []int{0, 1, 2, 3, 4, 5} {
		points = append(points, float32(s.Color.X))
		points = append(points, float32(s.Color.Y))
		points = append(points, float32(s.Color.Z))
	}

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
