package bug

import (
	"testing"

	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultPosition = mgl32.Vec3{0, 0, 0}
	Material_1      = material.Chrome
	Material_2      = material.Jade
	Material_3      = material.Gold
	DefaultSize     = float32(1)
)

type testShader struct {
}

func (t testShader) Use() {
}
func (t testShader) SetUniformMat4(s string, m mgl32.Mat4) {
}
func (t testShader) DrawTriangles(i int32) {
}
func (t testShader) Close(i int) {
}
func (t testShader) VertexAttribPointer(i uint32, c int32, s int32, o int) {
}
func (t testShader) BindVertexArray() {
}
func (t testShader) BindBufferData(d []float32) {
}
func (t testShader) SetUniform3f(s string, f1, f2, f3 float32) {
}
func (t testShader) SetUniform1f(s string, f1 float32) {
}

var shader testShader

func TestFirefly(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)

	if ff.position != DefaultPosition || ff.GetCenterPoint() != DefaultPosition {
		t.Error("Invalid position")
	}
	nullVector := mgl32.Vec3{0, 0, 0}
	if ff.direction != nullVector || ff.GetDirection() != nullVector {
		t.Error("Invalid direction")
	}
	if ff.axis != nullVector || ff.GetAxis() != nullVector {
		t.Error("Invalid axis")
	}
	if ff.speed != 0 {
		t.Error("Invalid speed")
	}
	if ff.angle != 0 || ff.GetAngle() != 0 {
		t.Error("Invalid angle")
	}
	if ff.size != DefaultSize {
		t.Error("Invalid size")
	}
}
func TestSetCenter(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	newCenter := mgl32.Vec3{1, 2, 3}
	ff.SetCenter(newCenter)
	if ff.GetCenterPoint() != newCenter {
		t.Error("Invalid position")
	}
}
func TestSetDirection(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	newDirection := mgl32.Vec3{1, 0, 0}
	ff.SetDirection(newDirection)
	if ff.GetDirection() != newDirection {
		t.Error("Invalid direction")
	}
}
func TestSetSpeed(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	newSpeed := float32(5.0)
	ff.SetSpeed(newSpeed)
	if ff.speed != newSpeed {
		t.Error("Invalid speed")
	}
}
func TestSetAngle(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	newAngle := float32(5.0)
	ff.SetAngle(newAngle)
	if ff.GetAngle() != newAngle {
		t.Error("Invalid angle")
	}
}
func TestSetAxis(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	newAxis := mgl32.Vec3{1, 0, 0}
	ff.SetAxis(newAxis)
	if ff.GetAxis() != newAxis {
		t.Error("Invalid axis")
	}
}
func TestDraw(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	ff.Draw()
}
func TestDrawWithUniforms(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	ff.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
}
func TestUpdate(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	ff.Update(10)
}
func TestDrawMode(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	ff.DrawMode(0)
}
func TestLog(t *testing.T) {
	ff := Firefly(DefaultPosition, DefaultSize, [3]*material.Material{Material_1, Material_2, Material_3}, shader)
	log := ff.Log()
	if len(log) < 5 {
		t.Errorf("Log too short. only '%d' chars.", len(log))
	}
}
