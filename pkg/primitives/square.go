package primitives

import (
	vao "github.com/akosgarai/opengl_playground/pkg/vertexarrayobject"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Square struct {
	A Vector
	B Vector
	C Vector
	D Vector
}

func NewSquare(v1, v2, v3, v4 Vector) *Square {
	return &Square{v1, v2, v3, v4}
}

func (t *Square) Draw() {
	var points []float32
	points = append(points, float32(t.A.X))
	points = append(points, float32(t.A.Y))
	points = append(points, float32(t.A.Z))
	points = append(points, float32(t.B.X))
	points = append(points, float32(t.B.Y))
	points = append(points, float32(t.B.Z))
	points = append(points, float32(t.C.X))
	points = append(points, float32(t.C.Y))
	points = append(points, float32(t.C.Z))
	points = append(points, float32(t.C.X))
	points = append(points, float32(t.C.Y))
	points = append(points, float32(t.C.Z))
	points = append(points, float32(t.D.X))
	points = append(points, float32(t.D.Y))
	points = append(points, float32(t.D.Z))
	points = append(points, float32(t.A.X))
	points = append(points, float32(t.A.Y))
	points = append(points, float32(t.A.Z))
	vertexArrayObject := vao.New(points)
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
