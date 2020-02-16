package main

import (
	"fmt"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs"
)

var (
	DebugPrint         = false
	mouseButtonPressed = false
	mousePositionX     = 0.0
	mousePositionY     = 0.0
)

type Application struct {
	Points []primitives.Point
}

func (a *Application) AddPoint(point primitives.Point) {
	a.Points = append(a.Points, point)
}

var app Application

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderDirectOutputSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderConstantSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

/*
* Mouse click handler logic:
* - if the left mouse button is not pressed, and the button is just released, App.AddPoint(), clean up the temp.point.
* - if the button is just pressed, set the point that needs to be added.
 */
func mouseHandler(window *glfw.Window) {
	x, y := window.GetCursorPos()

	if window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = x
			mousePositionY = y
			mouseButtonPressed = true
		}
	} else {
		if mouseButtonPressed {
			mouseButtonPressed = false
			x, y := convertMouseCoordinates()
			app.AddPoint(
				primitives.Point{
					primitives.Vector{x, y, 0.0},
					primitives.Vector{1, 1, 1},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyD) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Println(app.Points)
		}
	} else {
		DebugPrint = false
	}
}
func convertMouseCoordinates() (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (mousePositionX - halfWidth) / (halfWidth)
	y := (halfHeight - mousePositionY) / (halfHeight)
	return x, y
}
func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		mouseHandler(window)
		keyHandler(window)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
