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
func (v *VAO) appendVec2(v1 mgl32.Vec2) {
	v.vao = append(v.vao, v1.X())
	v.vao = append(v.vao, v1.Y())
}

// AppendVectors gets two vec3 input and appends them to the vao.
// In other words it appends 6 float32. It can be used for coordinate & color
// or coordinate & normal vector.
func (v *VAO) AppendVectors(v1, v2 mgl32.Vec3) {
	v.appendVector(v1)
	v.appendVector(v2)
}

// AppendTextureVectors gets two vec3 and a vec2 input and appends them to the vao.
// In other words it appends 8 float32. It can be used for coordinate & color & texture corrdinates.
func (v *VAO) AppendTextureVectors(v1, v2 mgl32.Vec3, tex mgl32.Vec2) {
	v.appendVector(v1)
	v.appendVector(v2)
	v.appendVec2(tex)
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
