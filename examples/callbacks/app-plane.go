package main

import (
	"fmt"
	"runtime"

	"github.com/akosgarai/opengl_playground/examples/callbacks/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - plane with ball"
)

type Application struct {
	window  *glfw.Window
	program uint32
	// camera related parameters
	camera *primitives.Camera
	square *primitives.Square
}

func NewApplication() *Application {
	var app Application
	app.camera = primitives.NewCamera(
		mgl32.Vec3{-10, -5, 20.0},
		mgl32.Vec3{0, 1, 0},
		-90.0,
		0.0)
	app.camera.SetupProjection(
		45,
		float32(windowWidth)/float32(windowHeight),
		0.1,
		100.0)
	// square
	squareColor := mgl32.Vec3{0, 1, 0}
	app.square = primitives.NewSquare(
		primitives.Point{mgl32.Vec3{-20, 0, -20}, squareColor},
		primitives.Point{mgl32.Vec3{20, 0, -20}, squareColor},
		primitives.Point{mgl32.Vec3{20, 0, 20}, squareColor},
		primitives.Point{mgl32.Vec3{-20, 0, 20}, squareColor})
	app.square.SetPrecision(10)
	return &app
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("KeyCallback has been called with the following options: key: '%d', scancode: '%d', action: '%d'!, mods: '%d'\n", key, scancode, action, mods)
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("MouseButtonCallback has been called with the following options: button: '%d', action: '%d'!, mods: '%d'\n", button, action, mods)
}

// Draw is responsible for drawing the screen.
func (a *Application) Draw() {
	modelLocation := gl.GetUniformLocation(a.program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(a.program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(a.program, gl.Str("projection\x00"))

	// mvp - modelview - projection matrix
	V := a.camera.GetViewMatrix()
	gl.UniformMatrix4fv(viewLocation, 1, false, &V[0])
	P := a.camera.GetProjectionMatrix()
	gl.UniformMatrix4fv(projectionLocation, 1, false, &P[0])
	M := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
	a.square.Draw()
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderModelViewProjectionSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderBasicSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)
	return program
}

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	gl.UseProgram(app.program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// register keyboard button callback
	app.window.SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.window.SetMouseButtonCallback(app.MouseButtonCallback)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
