package camera

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewCamera(t *testing.T) {
	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)

	if cam.pitch != pitch {
		t.Errorf("Invalid pitch instead of '%f', we have '%f'", pitch, cam.pitch)
	}
	if cam.yaw != yaw {
		t.Errorf("Invalid yaw instead of '%f', we have '%f'", yaw, cam.yaw)
	}
	if cam.cameraPosition != position {
		t.Errorf("Invalid position")
	}
	if cam.worldUp != worldUp {
		t.Errorf("Invalid worldUp")
	}
	var front, up, right mgl32.Vec3
	front = mgl32.Vec3{1, 0, 0}
	up = mgl32.Vec3{0, -1, 0}
	right = mgl32.Vec3{0, 0, -1}
	if cam.cameraFrontDirection != front {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != up {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != right {
		t.Error("Invalid right direction")
	}
}
func TestLog(t *testing.T) {
	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)
	log := cam.Log()

	if len(log) < 10 {
		t.Errorf("Log too short: '%s'", log)
	}
}
func TestWalk(t *testing.T) {
	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	amountToMove := float32(2)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.Walk(amountToMove)
	// the front is +X, amount is 2 -> position {2,0,0}
	expectedPos := mgl32.Vec3{2, 0, 0}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	var front, up, right mgl32.Vec3
	front = mgl32.Vec3{1, 0, 0}
	up = mgl32.Vec3{0, -1, 0}
	right = mgl32.Vec3{0, 0, -1}
	if cam.cameraFrontDirection != front {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != up {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != right {
		t.Error("Invalid right direction")
	}
}
func TestStrafe(t *testing.T) {
	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	amountToMove := float32(2)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.Strafe(amountToMove)
	// the right is -Z, amount is 2 -> position {0,0,-2}
	expectedPos := mgl32.Vec3{0, 0, -2}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	var front, up, right mgl32.Vec3
	front = mgl32.Vec3{1, 0, 0}
	up = mgl32.Vec3{0, -1, 0}
	right = mgl32.Vec3{0, 0, -1}
	if cam.cameraFrontDirection != front {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != up {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != right {
		t.Error("Invalid right direction")
	}
}
func TestLift(t *testing.T) {
	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	amountToMove := float32(2)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.Lift(amountToMove)
	// the up is -Y, amount is 2 -> position {0,-2,0}
	expectedPos := mgl32.Vec3{0, -2, 0}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	var front, up, right mgl32.Vec3
	front = mgl32.Vec3{1, 0, 0}
	up = mgl32.Vec3{0, -1, 0}
	right = mgl32.Vec3{0, 0, -1}
	if cam.cameraFrontDirection != front {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != up {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != right {
		t.Error("Invalid right direction")
	}
}
func TestSetupProjection(t *testing.T) {
	fov := float32(45)
	aspRatio := float32(1.0)
	near := float32(0.1)
	far := float32(100)

	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.SetupProjection(fov, aspRatio, near, far)
	if cam.projectionOptions.fov != fov {
		t.Errorf("Invalid fov instead of '%f', we have '%f'", fov, cam.projectionOptions.fov)
	}
	if cam.projectionOptions.aspectRatio != aspRatio {
		t.Errorf("Invalid aspectRatio instead of '%f', we have '%f'", aspRatio, cam.projectionOptions.aspectRatio)
	}
	if cam.projectionOptions.near != near {
		t.Errorf("Invalid near instead of '%f', we have '%f'", near, cam.projectionOptions.near)
	}
	if cam.projectionOptions.far != far {
		t.Errorf("Invalid far instead of '%f', we have '%f'", far, cam.projectionOptions.far)
	}
}
func TestGetProjectionMatrix(t *testing.T) {
	fov := float32(45)
	aspRatio := float32(1.0)
	near := float32(0.1)
	far := float32(100)

	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.SetupProjection(fov, aspRatio, near, far)
	cam.GetProjectionMatrix()
}
func TestGetViewMatrix(t *testing.T) {
	fov := float32(45)
	aspRatio := float32(1.0)
	near := float32(0.1)
	far := float32(100)

	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.SetupProjection(fov, aspRatio, near, far)
	cam.GetViewMatrix()
}
func TestUpdateDirection(t *testing.T) {
	fov := float32(45)
	aspRatio := float32(1.0)
	near := float32(0.1)
	far := float32(100)

	position := mgl32.Vec3{0, 0, 0}
	worldUp := mgl32.Vec3{0, 1, 0}
	yaw := float32(0)
	pitch := float32(0)

	cam := NewCamera(position, worldUp, yaw, pitch)
	cam.SetupProjection(fov, aspRatio, near, far)
	cam.UpdateDirection(0, 0)
	var front, up, right mgl32.Vec3
	front = mgl32.Vec3{1, 0, 0}
	up = mgl32.Vec3{0, -1, 0}
	right = mgl32.Vec3{0, 0, -1}
	if cam.cameraFrontDirection != front {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != up {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != right {
		t.Error("Invalid right direction")
	}
}
