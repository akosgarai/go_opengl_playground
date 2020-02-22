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
	windowTitle  = "Example - mesh deformer"
)

type Application struct {
	triangles             []primitives.Triangle
	defaultTriangleLength int
	triangleColorFront    primitives.Vector
	triangleColorBack     primitives.Vector
	camera                *primitives.Camera
	worldWidth            int
	worldHeight           int
	worldDepth            int
	worldUpDirection      *primitives.Vector
}

func NewApplication() *Application {
	var app Application
	app.worldWidth = 1000
	app.worldHeight = 1000
	app.worldDepth = 1000
	app.defaultTriangleLength = 10
	app.triangleColorFront = primitives.Vector{0, 0, 1}
	app.triangleColorBack = primitives.Vector{0, 0.5, 1}
	app.camera = primitives.NewCamera()
	app.camera.SetPosition(primitives.Vector{0, 0, 100})
	app.camera.TargetCameraSetTarget(primitives.Vector{0, 0, 0})
	app.camera.SetupProjection(45, float64(windowWidth/windowHeight))
	app.camera.UpDirection = primitives.Vector{0, 0, 1}
	return &app
}

func (a *Application) AddTriangle(triangle primitives.Triangle) {
	a.triangles = append(a.triangles, triangle)
}

// It generates a bunch of triangles and sets their color to static blue.
func (a *Application) GenerateTriangles() {
	rows := a.worldDepth / a.defaultTriangleLength
	cols := a.worldWidth / a.defaultTriangleLength
	length := a.defaultTriangleLength
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float64(j * length)
			topY := float64(i * length)
			topZ := 0.0

			triangle := *primitives.NewTriangle(
				primitives.Vector{topX, topY, topZ},
				primitives.Vector{topX, topY - float64(length), topZ},
				primitives.Vector{topX - float64(length), topY - float64(length), topZ},
			)
			triangle.SetColor(a.triangleColorFront)
			a.AddTriangle(triangle)
			triangle = *primitives.NewTriangle(
				primitives.Vector{topX, topY, topZ},
				primitives.Vector{topX - float64(length), topY - float64(length), topZ},
				primitives.Vector{topX - float64(length), topY, topZ},
			)
			triangle.SetColor(a.triangleColorBack)
			a.AddTriangle(triangle)
		}
	}
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderDeformVertexPositionSource, gl.VERTEX_SHADER)
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

	app := NewApplication()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	// Configure global settings
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	app.GenerateTriangles()

	nowUnix := time.Now().UnixNano()

	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// time
		elapsedTimeNano := time.Now().UnixNano() - nowUnix
		time := gl.GetUniformLocation(program, gl.Str("time\x00"))
		gl.Uniform1f(time, float32(elapsedTimeNano/10000000))
		// mvp - modelview - projection matrix
		MV := app.camera.GetViewMatrix()
		P := app.camera.GetProjectionMatrix()
		mvpValue := (P.Dot(MV)).Points
		gl.UniformMatrix4fv(mvpLocation, 1, false, &mvpValue[0])

		for _, item := range app.triangles {
			item.Draw()
		}
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
