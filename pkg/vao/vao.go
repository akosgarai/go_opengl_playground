package vao

import (
	"github.com/go-gl/mathgl/mgl32"
)

type VAO struct {
	vao []float32
}

func NewVAO() *VAO {
	return &VAO{
		vao: []float32{},
	}
}

func (v *VAO) appendVector(v1 mgl32.Vec3) {
	v.vao = append(v.vao, v1.X())
	v.vao = append(v.vao, v1.Y())
	v.vao = append(v.vao, v1.Z())
}

// AppendVectors gets two vec3 input and appends them to the vao.
// In other words it appends 6 float32. It can be used for coordinate & color
// or coordinate & normal vector.
func (v *VAO) AppendVectors(v1, v2 mgl32.Vec3) {
	v.appendVector(v1)
	v.appendVector(v2)
}

// AppendPoint gets two vec3 input and a float and appends them to the vao.
// In other words it appends 7 float32. It can be used for coordinate & color, & size for points
func (v *VAO) AppendPoint(v1, v2 mgl32.Vec3, size float32) {
	v.appendVector(v1)
	v.appendVector(v2)
	v.vao = append(v.vao, size)
}

// Get returns the vao as []float32
func (v *VAO) Get() []float32 {
	return v.vao
}

// Clear makes the vao empty.
func (v *VAO) Clear() {
	v.vao = []float32{}
}
