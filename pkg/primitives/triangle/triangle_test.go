package triangle

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

func TestNewTriangle(t *testing.T) {
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)

	if triangle.speed != 0.0 {
		t.Error("Speed should be 0")
	}
	if triangle.direction.X() != 0.0 || triangle.direction.Y() != 0.0 || triangle.direction.Z() != 0.0 {
		t.Error("Direction vector is not 0")
	}

	for i := 0; i < 3; i++ {
		if triangle.points[i] != points[i] {
			t.Error("Mismatch in the coordinates")
		}
		if triangle.colors[i] != colors[i] {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestSetColor(t *testing.T) {
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	triangle.SetColor(newColor)

	for i := 0; i < 3; i++ {
		if triangle.colors[i] != newColor {
			t.Error("Mismatch in the colors")
		}
	}

}
func TestSetIndexColor(t *testing.T) {
	origColor := mgl32.Vec3{1, 0, 0}
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{origColor, origColor, origColor}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	triangle.SetIndexColor(0, newColor)

	if triangle.colors[0] != newColor {
		t.Error("Mismatch in the new color")
	}
	for i := 1; i < 3; i++ {
		if triangle.colors[i] != origColor {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestLog(t *testing.T) {
	origColor := mgl32.Vec3{1, 0, 0}
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{origColor, origColor, origColor}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	log := triangle.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetupVao(t *testing.T) {
	origColor := mgl32.Vec3{1, 0, 0}
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{origColor, origColor, origColor}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	if len(triangle.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	triangle.setupVao()
	if len(triangle.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVao(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestDraw(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestDrawWithUniforms(t *testing.T) {
	t.Skip("It needs opengl init.")
}
func TestUpdate(t *testing.T) {
	origColor := mgl32.Vec3{1, 0, 0}
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{origColor, origColor, origColor}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	triangle.SetDirection(mgl32.Vec3{1, 0, 0})
	triangle.SetSpeed(1)
	dt := 10.0
	triangle.Update(dt)
	if triangle.points[0].X() != 10.0 || triangle.points[0].Y() != 0.0 || triangle.points[0].Z() != 0.0 {
		t.Error("Invalid position for p0")
	}
	if triangle.points[1].X() != 11.0 || triangle.points[1].Y() != 0.0 || triangle.points[1].Z() != 0.0 {
		t.Error("Invalid position for p1")
	}
	if triangle.points[2].X() != 10.5 || triangle.points[2].Y() != 1.0 || triangle.points[2].Z() != 0.0 {
		t.Error("Invalid position for p2")
	}
}
func TestSetDirection(t *testing.T) {
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	newDirection := mgl32.Vec3{1, 1, 0}
	triangle.SetDirection(newDirection)

	if triangle.direction != newDirection {
		t.Error("Mismatch in the direction")
	}
}
func TestSetIndexDirection(t *testing.T) {
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	triangle.SetIndexDirection(0, 1)

	if triangle.direction.X() != 1.0 || triangle.direction.Y() != 0.0 || triangle.direction.Z() != 0.0 {
		t.Error("Mismatch in the direction")
	}
}
func TestSetSpeed(t *testing.T) {
	points := [3]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{0.5, 1, 0},
	}
	colors := [3]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	triangle := NewTriangle(points, colors, shader)
	triangle.SetSpeed(10)

	if triangle.speed != 10.0 {
		t.Error("Mismatch in the speed")
	}
}
