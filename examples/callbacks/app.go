package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - callbacks"
)

type Application struct {
	window  *glfw.Window
	program uint32
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("KeyCallback has been called with the following options: key: '%d', scancode: '%d', action: '%d'!, mods: '%d'\n", key, scancode, action, mods)
}

func initGlfw() *glfw.Window {
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

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program := gl.CreateProgram()
	gl.LinkProgram(program)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	return program
}

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	// Configure global settings
	gl.UseProgram(app.program)
	app.window.SetKeyCallback(app.KeyCallback)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
