package main

import (
	"fmt"
	"runtime"

	mat "github.com/akosgarai/opengl_playground/pkg/primitives/matrix"
	tr "github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Example - static triangles, lots of them"
)

type Application struct {
	triangles []tr.Triangle
	rows      int
	length    float64

	window  *glfw.Window
	program uint32
}

// addTriangle inserts a new triangle to the application.
func (a *Application) addTriangle(t tr.Triangle) {
	a.triangles = append(a.triangles, t)
}

// GenerateTriangles fills up the triangles.
func (a *Application) GenerateTriangles() {
	for i := 0; i <= a.rows; i++ {
		for j := 0; j <= a.rows; j++ {
			topX := 1.0 - (float64(j) * a.length)
			topY := -1.0 + (float64(i) * a.length)
			topZ := 0.0

			a.addTriangle(
				*tr.NewTriangle(
					vec.Vector{topX, topY, topZ},
					vec.Vector{topX, topY - a.length, topZ},
					vec.Vector{topX - a.length, topY - a.length, topZ},
				))
		}
	}
}

// Draw calls the traingle draw.
func (a *Application) Draw() {
	for _, item := range a.triangles {
		item.Draw()
	}
}

func NewApplication(rows int) *Application {
	var app Application
	app.rows = rows
	/*
	 * The goal is to draw triangles to the screen. The screen will contain rows * rows triangles.
	 * The screen [-1, 1] -> rows part, one part : 2.0 / rows
	 */
	app.length = 2.0 / float64(app.rows)
	app.GenerateTriangles()
	return &app
}

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

	vertexShader, err := shader.CompileShader(shader.VertexShaderBasicSource, gl.VERTEX_SHADER)
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
	return program
}

func main() {
	runtime.LockOSThread()

	app := NewApplication(50)
	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	gl.UseProgram(app.program)
	mvpLocation := gl.GetUniformLocation(app.program, gl.Str("MVP\x00"))

	P := mat.UnitMatrix()
	MV := mat.UnitMatrix()
	mvp := P.Dot(MV).GetMatrix()
	gl.UniformMatrix4fv(mvpLocation, 1, false, &mvp[0])

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
