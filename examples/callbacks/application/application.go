package application

import (
	"github.com/akosgarai/opengl_playground/examples/callbacks/primitives"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Drawable interface {
	Draw()
	Update(float64)
}

type Application struct {
	window   *glfw.Window
	program  uint32
	camera   *primitives.Camera
	keyDowns map[string]bool

	items []Drawable
}

// New returns an application instance
func New() *Application {
	return &Application{
		keyDowns: make(map[string]bool),
		items:    []Drawable{},
	}
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

// AddItem inserts a new drawable item
func (a *Application) AddItem(d Drawable) {
	a.items = append(a.items, d)
}

// Draw calls Draw function in every drawable item.
func (a *Application) Draw() {
	for _, item := range a.items {
		item.Draw()
	}
}

// InitGlfw returns a *glfw.Windows instance.
func InitGlfw(windowWidth, windowHeight int, windowTitle string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window
}
