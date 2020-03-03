package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Triangle struct {
	A Point
	B Point
	C Point
}

func NewTriangle(v1, v2, v3 Vector) *Triangle {
	return &Triangle{Point{v1, Vector{1, 1, 1}}, Point{v2, Vector{1, 1, 1}}, Point{v3, Vector{1, 1, 1}}}
}

func (t *Triangle) SetColor(color Vector) {
	t.A.SetColor(color)
	t.B.SetColor(color)
	t.C.SetColor(color)
}
func (t *Triangle) appendPointToVao(currentVao []float32, p Point) []float32 {
	currentVao = append(currentVao, float32(p.Coordinate.X))
	currentVao = append(currentVao, float32(p.Coordinate.Y))
	currentVao = append(currentVao, float32(p.Coordinate.Z))
	currentVao = append(currentVao, float32(p.Color.X))
	currentVao = append(currentVao, float32(p.Color.Y))
	currentVao = append(currentVao, float32(p.Color.Z))
	return currentVao
}
func (t *Triangle) setupVao() []float32 {
	var points []float32

	points = t.appendPointToVao(points, t.A)
	points = t.appendPointToVao(points, t.B)
	points = t.appendPointToVao(points, t.C)

	return points
}

func (t *Triangle) buildVao() uint32 {
	points := t.setupVao()

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
func (t *Triangle) Draw() {
	vertexArrayObject := t.buildVao()
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
