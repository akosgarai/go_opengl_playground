package mesh

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestSetScale(t *testing.T) {
	var m Mesh
	scale := mgl32.Vec3{2, 2, 2}
	m.SetScale(scale)
	if m.scale != scale {
		t.Error("Scale mismatch")
	}
}
func TestSetRotationAngle(t *testing.T) {
	var m Mesh
	m.SetRotationAngle(1)
	if m.angle != 1 {
		t.Error("Angle mismatch")
	}
}
func TestGetRotationAngle(t *testing.T) {
	var m Mesh
	if m.angle != 0 {
		t.Error("Angle mismatch")
	}
	m.SetRotationAngle(1)
	if m.GetRotationAngle() != 1 {
		t.Error("Angle mismatch")
	}
}
func TestSetRotationAxis(t *testing.T) {
	var m Mesh
	axis := mgl32.Vec3{0, 1, 0}
	m.SetRotationAxis(axis)
	if m.axis != axis {
		t.Error("Axis mismatch")
	}
}
func TestSetPosition(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	if m.position != pos {
		t.Error("Position mismatch")
	}
}
func TestSetDirection(t *testing.T) {
	var m Mesh
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	if m.direction != dir {
		t.Error("Direction mismatch")
	}
}
func TestSetSpeed(t *testing.T) {
	var m Mesh
	m.SetSpeed(10)
	if m.velocity != 10 {
		t.Error("Speed mismatch")
	}
}
func TestGetPosition(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	if m.GetPosition() != pos {
		t.Error("Position mismatch")
	}
}
func TestGetDirection(t *testing.T) {
	var m Mesh
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	if m.GetDirection() != dir {
		t.Error("Direction mismatch")
	}
}
func TestUpdate(t *testing.T) {
	var m Mesh
	m.SetDirection(mgl32.Vec3{0, 0, 0})
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	m.Update(2)
	if m.GetPosition() != pos {
		t.Error("Invalid position after update")
	}
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	m.SetSpeed(10)
	m.Update(2)
	expectedPosition := mgl32.Vec3{0, 21, 2}
	if m.GetPosition() != expectedPosition {
		t.Error("Invalid position after update")
	}
}
func TestModelTransformation(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 0, 0}
	m.SetPosition(pos)
	axis := mgl32.Vec3{0, 0, 0}
	m.SetRotationAxis(axis)
	m.SetRotationAngle(0)
	scale := mgl32.Vec3{1, 1, 1}
	m.SetScale(scale)
	M := m.ModelTransformation()
	if M != mgl32.Ident4() {
		t.Error("Invalid model matrix")
	}
}
func TestAdd(t *testing.T) {
	t.Skip("Unimplemented")
}
