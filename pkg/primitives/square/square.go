package square

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	P "github.com/akosgarai/opengl_playground/pkg/primitives/point"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

type Square struct {
	A P.Point
	B P.Point
	C P.Point
	D P.Point
}

func NewSquare(v1, v2, v3, v4 vec.Vector) *Square {
	return &Square{P.Point{v1, vec.Vector{1, 1, 1}}, P.Point{v2, vec.Vector{1, 1, 1}}, P.Point{v3, vec.Vector{1, 1, 1}}, P.Point{v4, vec.Vector{1, 1, 1}}}
}

func (s *Square) SetColor(color vec.Vector) {
	s.A.SetColor(color)
	s.B.SetColor(color)
	s.C.SetColor(color)
	s.D.SetColor(color)
}

func (s *Square) buildVaoWithoutColor() []float32 {
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

	points = append(points, float32(s.A.Coordinate.X))
	points = append(points, float32(s.A.Coordinate.Y))
	points = append(points, float32(s.A.Coordinate.Z))

	points = append(points, float32(s.C.Coordinate.X))
	points = append(points, float32(s.C.Coordinate.Y))
	points = append(points, float32(s.C.Coordinate.Z))

	points = append(points, float32(s.D.Coordinate.X))
	points = append(points, float32(s.D.Coordinate.Y))
	points = append(points, float32(s.D.Coordinate.Z))

	return points
}
func (s *Square) appendPointToVao(currentVao []float32, p Point) []float32 {
	currentVao = append(currentVao, float32(p.Coordinate.X))
	currentVao = append(currentVao, float32(p.Coordinate.Y))
	currentVao = append(currentVao, float32(p.Coordinate.Z))
	currentVao = append(currentVao, float32(p.Color.X))
	currentVao = append(currentVao, float32(p.Color.Y))
	currentVao = append(currentVao, float32(p.Color.Z))
	return currentVao
}
func (s *Square) setupVao() []float32 {
	var points []float32

	points = s.appendPointToVao(points, s.A)
	points = s.appendPointToVao(points, s.B)
	points = s.appendPointToVao(points, s.C)
	points = s.appendPointToVao(points, s.A)
	points = s.appendPointToVao(points, s.C)
	points = s.appendPointToVao(points, s.D)

	return points
}
func (s *Square) buildAndSetupVao() uint32 {
	points := s.setupVao()

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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	return vertexArrayObject
}

func (s *Square) Draw() {
	vertexArrayObject := s.buildAndSetupVao()
	gl.BindVertexArray(vertexArrayObject)
	// The square is represented by 2 triangle, so we have 2 * 3 points here.
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (s *Square) SetupVaoPoligonMode() {
	vao := s.buildVaoWithoutColor()
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vao), gl.Ptr(vao), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.BindVertexArray(vertexArrayObject)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}
