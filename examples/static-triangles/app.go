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
	windowHeight = 600
	windowTitle  = "Example - static triangles, lots of them"
)

var (
	triangles []primitives.Triangle
)

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

func generateTriangles() {
	/*
	 * The goal is to draw triangles to the screen. The screen will contain 20 * 20 triangles.
	 * The screen [-1, 1] -> 20 part, one part : 0.10
	 */
	rows := 20
	cols := 20
	length := 0.10
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := 1.0 - (float64(j) * length)
			topY := -1.0 + (float64(i) * length)
			topZ := 0.0

			triangles = append(
				triangles,
				*primitives.NewTriangle(
					primitives.Vector{topX, topY, topZ},
					primitives.Vector{topX, topY - length, topZ},
					primitives.Vector{topX - length, topY - length, topZ},
				),
			)
		}
	}
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// Configure global settings
	gl.UseProgram(program)
	uniform := [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	// mvp - modelview - projection matrix
	mvp := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	gl.UniformMatrix4fv(mvp, 1, false, &uniform[0])

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	generateTriangles()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, item := range triangles {
			item.Draw()
		}
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
