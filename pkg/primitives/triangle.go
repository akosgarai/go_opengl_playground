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

func (t *Triangle) buildVao() uint32 {
	var points []float32

	// Coordinates
	points = append(points, float32(t.A.Coordinate.X))
	points = append(points, float32(t.A.Coordinate.Y))
	points = append(points, float32(t.A.Coordinate.Z))

	points = append(points, float32(t.B.Coordinate.X))
	points = append(points, float32(t.B.Coordinate.Y))
	points = append(points, float32(t.B.Coordinate.Z))

	points = append(points, float32(t.C.Coordinate.X))
	points = append(points, float32(t.C.Coordinate.Y))
	points = append(points, float32(t.C.Coordinate.Z))

	// Colors
	points = append(points, float32(t.A.Color.X))
	points = append(points, float32(t.A.Color.Y))
	points = append(points, float32(t.A.Color.Z))

	points = append(points, float32(t.B.Color.X))
	points = append(points, float32(t.B.Color.Y))
	points = append(points, float32(t.B.Color.Z))

	points = append(points, float32(t.C.Color.X))
	points = append(points, float32(t.C.Color.Y))
	points = append(points, float32(t.C.Color.Z))

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
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(4*3*3))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)

	return vertexArrayObject
}
func (t *Triangle) Draw() {
	vertexArrayObject := t.buildVao()
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
