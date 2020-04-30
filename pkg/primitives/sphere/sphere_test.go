package sphere

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

type testShader struct {
}

func (t testShader) Use() {
}
func (t testShader) SetUniformMat4(s string, m mgl32.Mat4) {
}

var shader testShader

func TestNew(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	sphere := New(center, color, radius, shader)
	if sphere.center != center {
		t.Error("Center mismatch")
	}
	if sphere.color != color {
		t.Error("Color mismatch")
	}
	if sphere.radius != radius {
		t.Error("Radius mismatch")
	}
	if sphere.speed != 0.0 {
		t.Error("Speed should be 0")
	}
	if sphere.direction.X() != 0.0 || sphere.direction.Y() != 0.0 || sphere.direction.Z() != 0.0 {
		t.Error("Direction vector is not 0")
	}
}
func TestLog(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	sphere := New(center, color, radius, shader)
	log := sphere.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetPrecision(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	prec := 20
	sphere := New(center, color, radius, shader)
	sphere.SetPrecision(prec)
	if sphere.precision != prec {
		t.Error("Precision mismatch")
	}

}
func TestSetCenter(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newCenter := mgl32.Vec3{0, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetCenter(newCenter)
	if sphere.center != newCenter {
		t.Error("Center mismatch")
	}
}
func TestGetCenter(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newCenter := mgl32.Vec3{0, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetCenter(newCenter)
	if sphere.GetCenter() != newCenter {
		t.Error("Center mismatch")
	}
}
func TestSetColor(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newColor := mgl32.Vec3{1, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetColor(newColor)
	if sphere.color != newColor {
		t.Error("Color mismatch")
	}
}
func TestGetColor(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newColor := mgl32.Vec3{1, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetColor(newColor)
	if sphere.GetColor() != newColor {
		t.Error("Color mismatch")
	}
}
func TestSetRadius(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newRadius := float32(4.0)
	sphere := New(center, color, radius, shader)
	sphere.SetRadius(newRadius)
	if sphere.radius != newRadius {
		t.Error("Radius mismatch")
	}
}
func TestGetRadius(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	newRadius := float32(4.0)
	sphere := New(center, color, radius, shader)
	sphere.SetRadius(newRadius)
	if sphere.GetRadius() != newRadius {
		t.Error("Radius mismatch")
	}
}
func TestSetupVao(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	sphere := New(center, color, radius, shader)
	if len(sphere.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	sphere.setupVao()
	if len(sphere.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVao(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestDrawWithUniforms(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestDraw(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestSetDirection(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	direction := mgl32.Vec3{1, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetDirection(direction)
	if sphere.direction != direction {
		t.Error("Direction mismatch")
	}
}
func TestSetIndexDirection(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	direction := mgl32.Vec3{1, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetIndexDirection(0, direction[0])
	if sphere.direction != direction {
		t.Error("Direction mismatch")
	}
}
func TestSetSpeed(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	speed := float32(5.0)
	sphere := New(center, color, radius, shader)
	sphere.SetSpeed(speed)
	if sphere.speed != speed {
		t.Error("Speed mismatch")
	}
}
func TestUpdate(t *testing.T) {
	center := mgl32.Vec3{3, 3, 5}
	color := mgl32.Vec3{0, 0, 1}
	radius := float32(2.0)
	speed := float32(5.0)
	direction := mgl32.Vec3{1, 0, 0}
	sphere := New(center, color, radius, shader)
	sphere.SetSpeed(speed)
	sphere.SetDirection(direction)
	sphere.Update(10)
	expectedNewCenter := mgl32.Vec3{53, 3, 5}
	if sphere.center != expectedNewCenter {
		t.Error("Center mismatch after update")
	}
}
