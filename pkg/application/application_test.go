package application

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type WindowMock struct {
}

func (wm WindowMock) GetCursorPos() (float64, float64) {
	return 0.0, 0.0
}
func (wm WindowMock) SetKeyCallback(cb glfw.KeyCallback) glfw.KeyCallback {
	return cb
}
func (wm WindowMock) SetMouseButtonCallback(cb glfw.MouseButtonCallback) glfw.MouseButtonCallback {
	return cb
}
func (wm WindowMock) ShouldClose() bool {
	return false
}
func (wm WindowMock) SwapBuffers() {
}

var wm WindowMock

type CameraMock struct {
}

func (cm CameraMock) Log() string {
	return ""
}
func (cm CameraMock) GetViewMatrix() mgl32.Mat4 {
	return mgl32.Ident4()
}
func (cm CameraMock) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Ident4()
}
func (cm CameraMock) Walk(float32) {
}
func (cm CameraMock) Strafe(float32) {
}
func (cm CameraMock) Lift(float32) {
}
func (cm CameraMock) UpdateDirection(float32, float32) {
}

var cm CameraMock

func TestNew(t *testing.T) {
	app := New()
	if len(app.items) != 0 {
		t.Error("Invalid application - items should be empty")
	}
	if app.cameraSet {
		t.Error("Camera shouldn't be set")
	}
}
func TestLog(t *testing.T) {
	app := New()
	log := app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
}
func TestSetWindow(t *testing.T) {
	app := New()
	app.SetWindow(wm)

	if app.window != wm {
		t.Error("Invalid window setup.")
	}
}
func TestGetWindow(t *testing.T) {
	app := New()
	app.SetWindow(wm)

	if app.GetWindow() != wm {
		t.Error("Invalid window setup.")
	}
}
func TestSetCamera(t *testing.T) {
	app := New()
	app.SetCamera(cm)

	if app.camera != cm {
		t.Error("Invalid camera setup.")
	}
}
func TestGetCamera(t *testing.T) {
	app := New()
	app.SetCamera(cm)

	if app.GetCamera() != cm {
		t.Error("Invalid camera setup.")
	}
}
func TestSetMouseButtons(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestGetMouseButtons(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestSetKeys(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestGetKeys(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestAddItem(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestDraw(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestUpdate(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestDrawWithUniforms(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestKeyCallback(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestMouseButtonCallback(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestSetKeyState(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestSetButtonState(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestGetMouseButtonState(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestGetKeyState(t *testing.T) {
	t.Skip("Unimplemented")
}
