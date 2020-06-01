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
