package application

import (
	"fmt"

	"github.com/akosgarai/opengl_playground/examples/light/primitives"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG = glfw.KeyH
)

type Drawable interface {
	Draw()
	DrawWithUniforms(mgl32.Mat4, mgl32.Mat4)
	Update(float64)
	Log() string
}

type Application struct {
	window   *glfw.Window
	program  uint32
	camera   *primitives.Camera
	keyDowns map[glfw.Key]bool

	items []Drawable
}

// New returns an application instance
func New() *Application {
	return &Application{
		keyDowns: make(map[glfw.Key]bool),
		items:    []Drawable{},
	}
}

// Log returns the string representation of this object.
func (a *Application) Log() string {
	logString := "Application:\n"
	logString += " - camera : " + a.camera.Log() + "\n"
	logString += " - items :\n"
	for _, item := range a.items {
		logString += item.Log()
	}
	return logString
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
func (a *Application) SetKeys(m map[glfw.Key]bool) {
	a.keyDowns = m
}

// GetKeys returns the current keyDowns of the application.
func (a *Application) GetKeys() map[glfw.Key]bool {
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

// DrawWithUniforms calls DrawWithUniforms function in every drawable item with the calculated V & P.
func (a *Application) DrawWithUniforms() {
	V := a.camera.GetViewMatrix()
	P := a.camera.GetProjectionMatrix()

	for _, item := range a.items {
		item.DrawWithUniforms(V, P)
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

// InitOpenGL is for initializing the gl lib. It also prints out the gl version.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

// DummyKeyCallback is responsible for the keyboard event handling with log.
// So this function does nothing but printing out the input parameters.
func (a *Application) DummyKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("KeyCallback has been called with the following options: key: '%d', scancode: '%d', action: '%d'!, mods: '%d'\n", key, scancode, action, mods)
}

// DummyMouseButtonCallback is responsible for the mouse button event handling with log.
// So this function does nothing but printing out the input parameters.
func (a *Application) DummyMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("MouseButtonCallback has been called with the following options: button: '%d', action: '%d'!, mods: '%d'\n", button, action, mods)
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case DEBUG:
		if action != glfw.Release {
			fmt.Printf("%s\n", a.Log())
		}
		break
	default:
		a.SetKeyState(key, action)
		break
	}
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetKeyState(key glfw.Key, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.keyDowns[key] = isButtonPressed
}

// SetKeyState returns the state of the given key
func (a *Application) GetKeyState(key glfw.Key) bool {
	return a.keyDowns[key]
}
