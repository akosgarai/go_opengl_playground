package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Example - static button handler"
	epsilon      = 300
)

var (
	triangle = primitives.NewTriangle(
		primitives.Vector{-0.75, 0.75, 0}, // top
		primitives.Vector{-0.75, 0.25, 0}, // left
		primitives.Vector{-0.25, 0.25, 0}, // right
	)
	square = primitives.NewSquare(
		primitives.Vector{0.25, -0.25, 0}, // top-left
		primitives.Vector{0.25, -0.75, 0}, // bottom-left
		primitives.Vector{0.75, -0.75, 0}, // bottom-right
		primitives.Vector{0.75, -0.25, 0}, // top-right
	)
	lastUpdate = time.Now().UnixNano() / 1000000
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderCookBookSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderCookBookSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyT) == glfw.Press {
		triangle.SetColor(primitives.Vector{0, 1, 0})
		square.SetColor(primitives.Vector{1, 0, 0})
	} else {
		triangle.SetColor(primitives.Vector{1, 0, 0})
		square.SetColor(primitives.Vector{0, 1, 0})
	}
	nowUnixM := time.Now().UnixNano() / 1000000
	if lastUpdate+epsilon > nowUnixM {
		return
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		square.A.Y += 0.05
		square.B.Y += 0.05
		square.C.Y += 0.05
		square.D.Y += 0.05
		triangle.A.Y -= 0.05
		triangle.B.Y -= 0.05
		triangle.C.Y -= 0.05
		lastUpdate = nowUnixM
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		square.A.X -= 0.05
		square.B.X -= 0.05
		square.C.X -= 0.05
		square.D.X -= 0.05
		triangle.A.X += 0.05
		triangle.B.X += 0.05
		triangle.C.X += 0.05
		lastUpdate = nowUnixM
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		square.A.Y -= 0.05
		square.B.Y -= 0.05
		square.C.Y -= 0.05
		square.D.Y -= 0.05
		triangle.A.Y += 0.05
		triangle.B.Y += 0.05
		triangle.C.Y += 0.05
		lastUpdate = nowUnixM
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		square.A.X += 0.05
		square.B.X += 0.05
		square.C.X += 0.05
		square.D.X += 0.05
		triangle.A.X -= 0.05
		triangle.B.X -= 0.05
		triangle.C.X -= 0.05
		lastUpdate = nowUnixM
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

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		keyHandler(window)
		gl.UseProgram(program)
		triangle.Draw()
		square.Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
