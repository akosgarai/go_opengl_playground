package main

import (
	"fmt"
	"runtime"

	P "github.com/akosgarai/opengl_playground/pkg/primitives/point"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
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
	Points []P.Point

	window  *glfw.Window
	program uint32
}

func NewApplication() *Application {
	return &Application{}
}

// AddPoint inserts a new point to the points.
func (a *Application) AddPoint(point P.Point) {
	a.Points = append(a.Points, point)
}

// Basic function for glfw initialization.
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderPointSource, gl.VERTEX_SHADER)
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
func (a *Application) MouseHandler() {
	x, y := a.window.GetCursorPos()

	if a.window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = x
			mousePositionY = y
			mouseButtonPressed = true
		}
	} else {
		if mouseButtonPressed {
			mouseButtonPressed = false
			mX, mY := trans.MouseCoordinates(x, y, windowWidth, windowHeight)
			a.AddPoint(
				P.Point{
					vec.Vector{mX, mY, 0.0},
					vec.Vector{1, 1, 1},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func (a *Application) KeyHandler() {
	if a.window.GetKey(glfw.KeyD) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Println(a.Points)
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
func (a *Application) buildVAO() []float32 {
	var vao []float32
	for _, item := range a.Points {
		vao = append(vao, float32(item.Coordinate.X))
		vao = append(vao, float32(item.Coordinate.Y))
		vao = append(vao, float32(item.Coordinate.Z))
	}
	return vao
}
func (a *Application) Draw() {
	if len(a.Points) < 1 {
		return
	}
	points := a.buildVAO()
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(0))

	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.POINTS, 0, int32(len(a.Points)))
}

func main() {
	runtime.LockOSThread()
	app := NewApplication()
	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	gl.UseProgram(app.program)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		app.MouseHandler()
		app.KeyHandler()
		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
