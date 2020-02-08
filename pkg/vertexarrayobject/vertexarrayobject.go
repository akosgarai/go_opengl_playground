package vertexarrayobject

import (
	"github.com/akosgarai/opengl_playground/pkg/primitives"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// New initializes and returns a vertex array from the points provided.
func New(vectors []primitives.Vector) uint32 {
	var points []float32
	for _, vector := range vectors {
		points = append(points, float32(vector.X))
		points = append(points, float32(vector.Y))
		points = append(points, float32(vector.Z))
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
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vertexArrayObject
}
