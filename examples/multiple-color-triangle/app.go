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
	windowTitle  = "Example - static triangle with multiple color"
)

var (
	triangle = tr.NewTriangle(
		vec.Vector{-0.75, 0.75, 0}, // top
		vec.Vector{-0.75, 0.25, 0}, // left
		vec.Vector{-0.25, 0.25, 0}, // right
	)
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

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// Configure global settings
	gl.UseProgram(program)

	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))

	P := mat.UnitMatrix()
	MV := mat.UnitMatrix()
	mvp := P.Dot(MV).GetMatrix()
	gl.UniformMatrix4fv(mvpLocation, 1, false, &mvp[0])

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	// red
	triangle.A.SetColor(vec.Vector{1, 0, 0})
	// green
	triangle.B.SetColor(vec.Vector{0, 1, 0})
	// blue
	triangle.C.SetColor(vec.Vector{0, 0, 1})

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		triangle.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
