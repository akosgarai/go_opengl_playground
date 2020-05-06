package sphere

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultRadius    = float32(2.0)
	DefaultColor     = mgl32.Vec3{0, 0, 1}
	DefaultCenter    = mgl32.Vec3{3, 3, 5}
	DefaultSpeed     = float32(0.0)
	DefaultDirection = mgl32.Vec3{0, 0, 0}
	DefaultAngle     = float32(0.0)
	DefaultAxis      = mgl32.Vec3{0, 0, 0}
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

var shader testShader

func TestNew(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	if sphere.center != DefaultCenter {
		t.Error("Center mismatch")
	}
	if sphere.color != DefaultColor {
		t.Error("Color mismatch")
	}
	if sphere.radius != DefaultRadius {
		t.Error("Radius mismatch")
	}
	if sphere.speed != DefaultSpeed {
		t.Error("Speed should be 0")
	}
	if sphere.direction != DefaultDirection {
		t.Error("Direction vector is not 0")
	}
	if sphere.GetDirection() != DefaultDirection {
		t.Error("Direction vector is not 0")
	}
	if sphere.axis != DefaultAxis {
		t.Error("Axis vector is not 0")
	}
	if sphere.angle != DefaultAngle {
		t.Error("Angle is not null")
	}
}
func TestLog(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	log := sphere.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetPrecision(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	prec := 20
	sphere.SetPrecision(prec)
	if sphere.precision != prec {
		t.Errorf("Precision mismatch. instead of '%d', we have '%d'", prec, sphere.precision)
	}

}
func TestSetCenter(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newCenter := mgl32.Vec3{0, 0, 0}
	sphere.SetCenter(newCenter)
	if sphere.center != newCenter {
		t.Error("Center mismatch")
	}
}
func TestGetCenter(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newCenter := mgl32.Vec3{0, 0, 0}
	sphere.SetCenter(newCenter)
	if sphere.GetCenter() != newCenter {
		t.Error("Center mismatch")
	}
}
func TestSetColor(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newColor := mgl32.Vec3{1, 0, 0}
	sphere.SetColor(newColor)
	if sphere.color != newColor {
		t.Error("Color mismatch")
	}
}
func TestGetColor(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newColor := mgl32.Vec3{1, 0, 0}
	sphere.SetColor(newColor)
	if sphere.GetColor() != newColor {
		t.Error("Color mismatch")
	}
}
func TestSetRadius(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newRadius := float32(4.0)
	sphere.SetRadius(newRadius)
	if sphere.radius != newRadius {
		t.Error("Radius mismatch")
	}
}
func TestGetRadius(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	newRadius := float32(4.0)
	sphere.SetRadius(newRadius)
	if sphere.GetRadius() != newRadius {
		t.Error("Radius mismatch")
	}
}
func TestSetupVao(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	if len(sphere.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	sphere.setupVao()
	if len(sphere.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVao(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	if len(sphere.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	sphere.buildVao()
	if len(sphere.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawWithUniforms(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	if len(sphere.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	sphere.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(sphere.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDraw(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	if len(sphere.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	sphere.Draw()
	if len(sphere.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestSetDirection(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	direction := mgl32.Vec3{1, 0, 0}
	sphere.SetDirection(direction)
	if sphere.direction != direction {
		t.Error("Direction mismatch")
	}
	if sphere.GetDirection() != direction {
		t.Error("Direction mismatch")
	}
}
func TestSetIndexDirection(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	direction := mgl32.Vec3{1, 0, 0}
	sphere.SetIndexDirection(0, direction[0])
	if sphere.direction != direction {
		t.Error("Direction mismatch")
	}
	if sphere.GetDirection() != direction {
		t.Error("Direction mismatch")
	}
}
func TestSetSpeed(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	speed := float32(5.0)
	sphere.SetSpeed(speed)
	if sphere.speed != speed {
		t.Error("Speed mismatch")
	}
}
func TestUpdate(t *testing.T) {
	sphere := New(DefaultCenter, DefaultColor, DefaultRadius, shader)
	direction := mgl32.Vec3{1, 0, 0}
	speed := float32(5.0)
	sphere.SetSpeed(speed)
	sphere.SetDirection(direction)
	sphere.Update(10)
	expectedNewCenter := mgl32.Vec3{53, 3, 5}
	if sphere.center != expectedNewCenter {
		t.Error("Center mismatch after update")
	}
}
