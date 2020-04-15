package application

import (
	"github.com/akosgarai/opengl_playground/examples/callbacks/primitives"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Application struct {
	window   *glfw.Window
	program  uint32
	camera   *primitives.Camera
	keyDowns map[string]bool
}

// New returns an application instance
func New() *Application {
	return &Application{}
}

// SetWindow updates the window with the new one.
func (a *Application) SetWindow(w *glfw.Window) {
	a.window = w
}

// GetWindow returns the current window of the application.
func (a *Application) GetWindow() *glfw.Window {
	return a.window
}

// SetProgram updates the program with the new one.
func (a *Application) SetProgram(p uint32) {
	a.program = p
}

// GetProgram returns the current shader program of the application.
func (a *Application) GetProgram() uint32 {
	return a.program
}

// SetCamera updates the camera with the new one.
func (a *Application) SetCamera(c *primitives.Camera) {
	a.camera = c
}

// GetCamera returns the current camera of the application.
func (a *Application) GetCamera() *primitives.Camera {
	return a.camera
}

// SetKeys updates the keyDowns with the new one.
func (a *Application) SetKeys(m map[string]bool) {
	a.keyDowns = m
}

// GetKeys returns the current keyDowns of the application.
func (a *Application) GetKeys() map[string]bool {
	return a.keyDowns
}
