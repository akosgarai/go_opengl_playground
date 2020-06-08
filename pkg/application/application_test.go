package application

import (
	"testing"

	"github.com/akosgarai/opengl_playground/pkg/testhelper"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var wm testhelper.WindowMock

var cm testhelper.CameraMock

func TestNew(t *testing.T) {
	app := New()
	if len(app.shaderMap) != 0 {
		t.Error("Invalid application - shadermap should be empty")
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
	app.SetCamera(cm)
	log = app.Log()
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
func TestSetKeyState(t *testing.T) {
	app := New()
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.keyDowns[glfw.KeyW] {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.keyDowns[glfw.KeyW] {
		t.Error("W should be released")
	}
}
func TestSetButtonState(t *testing.T) {
	app := New()
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Press)
	if !app.mouseDowns[glfw.MouseButtonLeft] {
		t.Error("LMB should be pressed")
	}
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Release)
	if app.mouseDowns[glfw.MouseButtonLeft] {
		t.Error("LMB should be released")
	}
}
func TestGetMouseButtonState(t *testing.T) {
	app := New()
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Press)
	if !app.GetMouseButtonState(glfw.MouseButtonLeft) {
		t.Error("LMB should be pressed")
	}
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Release)
	if app.GetMouseButtonState(glfw.MouseButtonLeft) {
		t.Error("LMB should be released")
	}
}
func TestGetKeyState(t *testing.T) {
	app := New()
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.GetKeyState(glfw.KeyW) {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.GetKeyState(glfw.KeyW) {
		t.Error("W should be released")
	}
}
func TestSetCameraMovementMap(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestCameraKeyboardMovement(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestCameraKeyboardRotation(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestCameraMouseRotation(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestApplyMouseRotation(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestSetRotateOnEdgeDistance(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestCameraCollisionTest(t *testing.T) {
	t.Skip("Unimplemented")
}
