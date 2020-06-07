package mesh

import (
	"testing"

	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/primitives/boundingobject"
	"github.com/akosgarai/opengl_playground/pkg/testhelper"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	wrapperMock testhelper.GLWrapperMock
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
func TestSetParent(t *testing.T) {
	var m Mesh
	var parent interfaces.Mesh
	m.SetParent(parent)
	if m.parentSet != true {
		t.Error("After setting the parent, the flag supposed to be true")
	}
	if m.parent != parent {
		t.Error("The parent supposed to be the same")
	}
}
func TestIsParentMesh(t *testing.T) {
	var m Mesh
	var parent interfaces.Mesh
	if m.IsParentMesh() != true {
		t.Error("Before setting the parent, it should return true")
	}
	m.SetParent(parent)
	if m.IsParentMesh() != false {
		t.Error("After setting the parent, it should return false")
	}
}
func TestTransformationGettersWithParent(t *testing.T) {
	var m Mesh
	parent := NewPointMesh(wrapperMock)
	parent.position = mgl32.Vec3{1.0, 0.0, 0.0}
	m.SetParent(parent)
	modelTr := m.ModelTransformation()
	if modelTr == mgl32.Ident4() {
		t.Error("Model tr shouldn't be ident, if the parent transformation is set.")
	}
}
func TestRotateDirection(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	rotationAxis := mgl32.Vec3{0.0, 1.0, 0.0}
	m.rotateDirection(rotationAngle, rotationAxis)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.rotateDirection(rotationAngle, rotationAxis)
	if !m.direction.ApproxEqualThreshold(upDir, 0.0001) {
		t.Log(m.direction)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	m.direction = leftDir
	m.rotateDirection(rotationAngle, rotationAxis)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating different dir and axis should change the direction.")
	}
}
func TestRotatePosition(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	rotationAxis := mgl32.Vec3{0.0, 1.0, 0.0}
	m.RotatePosition(rotationAngle, rotationAxis)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upPos := mgl32.Vec3{0.0, 1.0, 0.0}
	leftPos := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontPos := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.position != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.position = upPos
	m.RotatePosition(rotationAngle, rotationAxis)
	if !m.position.ApproxEqualThreshold(upPos, 0.0001) {
		t.Log(m.position)
		t.Error("Rotating the same pos as axis shouldn't change the position.")
	}
	m.position = leftPos
	m.RotatePosition(rotationAngle, rotationAxis)
	if !m.position.ApproxEqualThreshold(frontPos, 0.001) {
		t.Log(m.position)
		t.Log(frontPos)
		t.Error("Rotating different pos and axis should change the position.")
	}
}
func TestRotateY(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateY(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.yaw != rotationAngle {
		t.Error("RotateY should update the yaw")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateY(rotationAngle)
	if !m.direction.ApproxEqualThreshold(upDir, 0.0001) {
		t.Log(m.direction)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.yaw != rotationAngle*2 {
		t.Error("RotateY should update the yaw")
	}
	m.direction = leftDir
	m.RotateY(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.yaw != rotationAngle*3 {
		t.Error("RotateY should update the yaw")
	}
}
func TestRotateX(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateX(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.pitch != rotationAngle {
		t.Error("RotateX should update the pitch")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateX(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.pitch != rotationAngle*2 {
		t.Error("RotateX should update the pitch")
	}
	m.direction = leftDir
	m.RotateX(rotationAngle)
	if !m.direction.ApproxEqualThreshold(leftDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.pitch != rotationAngle*3 {
		t.Error("RotateX should update the pitch")
	}
}
func TestRotateZ(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateZ(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.roll != rotationAngle {
		t.Error("RotateZ should update the roll")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateZ(rotationAngle)
	if !m.direction.ApproxEqualThreshold(leftDir, 0.001) {
		t.Log(m.direction)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.roll != rotationAngle*2 {
		t.Error("RotateZ should update the roll")
	}
	m.direction = frontDir
	m.RotateZ(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.roll != rotationAngle*3 {
		t.Error("RotateZ should update the roll")
	}
}
func TestIsBoundingObjectParamsSet(t *testing.T) {
	var m Mesh
	boParams := make(map[string]float32)
	if m.IsBoundingObjectSet() != false {
		t.Error("Before setting the bo, it should return false")
	}
	m.SetBoundingObject(boundingobject.New("AABB", boParams))
	if m.IsBoundingObjectSet() != true {
		t.Error("After setting the bo, it should return true")
	}
}
