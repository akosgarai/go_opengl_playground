package application

import (
	"reflect"
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

type DrawableMock struct {
}

func (dm DrawableMock) Draw() {
}
func (dm DrawableMock) DrawWithUniforms(m1, m2 mgl32.Mat4) {
}
func (dm DrawableMock) Update(float64) {
}
func (dm DrawableMock) Log() string {
	return ""
}

var dm DrawableMock

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
	app.SetCamera(cm)
	log = app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
	app.AddItem(dm)
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
func TestSetMouseButtons(t *testing.T) {
	mbs := make(map[glfw.MouseButton]bool)
	mbs[glfw.MouseButtonLeft] = true
	app := New()
	app.SetMouseButtons(mbs)
	if !reflect.DeepEqual(app.mouseDowns, mbs) {
		t.Error("Invalid mouse states")
	}
}
func TestGetMouseButtons(t *testing.T) {
	mbs := make(map[glfw.MouseButton]bool)
	mbs[glfw.MouseButtonLeft] = true
	app := New()
	app.SetMouseButtons(mbs)
	if !reflect.DeepEqual(app.GetMouseButtons(), mbs) {
		t.Error("Invalid mouse states")
	}
}
func TestSetKeys(t *testing.T) {
	ks := make(map[glfw.Key]bool)
	ks[glfw.KeyW] = true
	app := New()
	app.SetKeys(ks)
	if !reflect.DeepEqual(app.keyDowns, ks) {
		t.Error("Invalid key states")
	}
}
func TestGetKeys(t *testing.T) {
	ks := make(map[glfw.Key]bool)
	ks[glfw.KeyW] = true
	app := New()
	app.SetKeys(ks)
	if !reflect.DeepEqual(app.GetKeys(), ks) {
		t.Error("Invalid key states")
	}
}
func TestAddItem(t *testing.T) {
	app := New()
	if len(app.items) != 0 {
		t.Error("Invalid item length")
	}
	app.AddItem(dm)
	if len(app.items) != 1 {
		t.Error("Invalid item length")
	}
}
func TestDraw(t *testing.T) {
	app := New()
	app.Draw()
	app.AddItem(dm)
	app.Draw()
}
func TestUpdate(t *testing.T) {
	app := New()
	app.Update(10)
	app.AddItem(dm)
	app.Update(10)
}
func TestDrawWithUniforms(t *testing.T) {
	app := New()
	app.DrawWithUniforms()
	app.AddItem(dm)
	app.DrawWithUniforms()
	app.SetCamera(cm)
	app.DrawWithUniforms()
}
func TestKeyCallback(t *testing.T) {
	t.Skip("Unimplemented - glfw needed")
}
func TestMouseButtonCallback(t *testing.T) {
	t.Skip("Unimplemented - glfw needed")
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
