package point

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

func getPoint() *Point {
	return &Point{
		coordinate: mgl32.Vec3{0, 0, 0},
		color:      mgl32.Vec3{0, 0, 1},
		size:       1.0,

		direction: mgl32.Vec3{0, 0, 0},
		speed:     0.0,
	}
}
func TestNew(t *testing.T) {
	points := New(shader)
	if len(points.vao.Get()) != 0.0 {
		t.Error("Vao should be empty")
	}
	if len(points.points) != 0.0 {
		t.Error("Points should be empty")
	}
}
func TestSetColor(t *testing.T) {
	point := getPoint()
	color := mgl32.Vec3{1, 1, 1}
	point.SetColor(color)

	if point.color != color {
		t.Error("Color should be updated")
	}
}
func TestSetSpeed(t *testing.T) {
	point := getPoint()
	speed := float32(5.0)
	point.SetSpeed(speed)
	if point.speed != speed {
		t.Error("Speed should be updated")
	}
}
func TestSetDirection(t *testing.T) {
	point := getPoint()
	dir := mgl32.Vec3{0, 1, 0}
	point.SetDirection(dir)
	if point.direction != dir {
		t.Error("Direction should be updated")
	}
}
func TestSetIndexDirection(t *testing.T) {
	point := getPoint()
	expectedDir := mgl32.Vec3{0, 1, 0}
	point.SetIndexDirection(1, 1)
	if point.direction != expectedDir {
		t.Error("Direction should be updated")
	}
}
func TestAdd(t *testing.T) {
	points := New(shader)
	coords := mgl32.Vec3{1, 1, 1}
	col := mgl32.Vec3{1, 0, 0}
	size := float32(3.0)

	insertedPoint := points.Add(coords, col, size)

	if insertedPoint.coordinate != coords {
		t.Error("Coordinate mismatch")
	}
	if insertedPoint.color != col {
		t.Error("Color mismatch")
	}
	if insertedPoint.size != size {
		t.Error("Size mismatch")
	}
}
func TestUpdate(t *testing.T) {
	points := New(shader)
	coords := mgl32.Vec3{0, 0, 0}
	col := mgl32.Vec3{1, 0, 0}
	size := float32(3.0)
	points.Update(10.0)
	point := points.Add(coords, col, size)

	if point.coordinate.Y() != 0.0 || point.coordinate.X() != 0.0 || point.coordinate.Z() != 0.0 {
		t.Error("Invalid coordinates after Update - not moving")
	}

	point.SetSpeed(2.0)
	point.SetDirection(mgl32.Vec3{0, 1, 0})

	points.Update(10.0)

	if point.coordinate.Y() != float32(20.0) || point.coordinate.X() != 0.0 || point.coordinate.Z() != 0.0 {
		t.Error("Invalid coordinates after Update - moving")
		t.Log(point)
		t.Log(points)
	}
}
func TestLog(t *testing.T) {
	points := New(shader)
	log := points.Log()
	if len(log) < 6 {
		t.Error("Log too short")
	}
	coords := mgl32.Vec3{0, 0, 0}
	col := mgl32.Vec3{1, 0, 0}
	size := float32(3.0)
	points.Add(coords, col, size)
	log = points.Log()
	if len(log) < 16 {
		t.Error("Log too short")
	}
}
func TestSetupVao(t *testing.T) {
	points := New(shader)
	points.setupVao()
	if len(points.vao.Get()) != 0 {
		t.Error("Invalid vao length")
	}
	coords := mgl32.Vec3{0, 0, 0}
	col := mgl32.Vec3{1, 0, 0}
	size := float32(3.0)
	points.Add(coords, col, size)
	points.setupVao()
	if len(points.vao.Get()) != 7 {
		t.Error("Invalid vao length")
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
