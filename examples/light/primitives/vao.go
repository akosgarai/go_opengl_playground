package primitives

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

func (v *VAO) AppendPoint(p Point) {
	v.vao = append(v.vao, p.Coordinate.X())
	v.vao = append(v.vao, p.Coordinate.Y())
	v.vao = append(v.vao, p.Coordinate.Z())
	v.vao = append(v.vao, p.Color.X())
	v.vao = append(v.vao, p.Color.Y())
	v.vao = append(v.vao, p.Color.Z())
}
func (v *VAO) AppendTrianglePoints(pa, pb, pc Point) {
	v.AppendPoint(pa)
	v.AppendPoint(pb)
	v.AppendPoint(pc)
}
func (v *VAO) appendVector(v1 mgl32.Vec3) {
	v.vao = append(v.vao, v1.X())
	v.vao = append(v.vao, v1.Y())
	v.vao = append(v.vao, v1.Z())
}
func (v *VAO) AppendVectors(v1, v2 mgl32.Vec3) {
	v.appendVector(v1)
	v.appendVector(v2)
}
func (v *VAO) Get() []float32 {
	return v.vao
}
func (v *VAO) Clear() {
	v.vao = []float32{}
}
