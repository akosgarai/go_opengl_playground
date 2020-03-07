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

// It generates a bunch of triangles and sets their color to static blue.
func generateTrianglesModelCoordinates() {
	rows := 100
	cols := 100
	length := 10.0
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := (float64(j) * length)
			topY := (float64(i) * length)
			topZ := 0.0

			triangle := *primitives.NewTriangle(
				primitives.Vector{topX, topY, topZ},
				primitives.Vector{topX, topY - length, topZ},
				primitives.Vector{topX - length, topY - length, topZ},
			)
			triangle.SetColor(primitives.Vector{0, 0, 1})
			triangles = append(
				triangles,
				triangle,
			)
			triangle = *primitives.NewTriangle(
				primitives.Vector{topX, topY, topZ},
				primitives.Vector{topX - length, topY - length, topZ},
				primitives.Vector{topX - length, topY, topZ},
			)
			triangle.SetColor(primitives.Vector{0, 0.5, 1})
			triangles = append(
				triangles,
				triangle,
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

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	generateTrianglesModelCoordinates()

	nowUnix := time.Now().UnixNano()

	// mvp - modelview - projection matrix
	angelOfView := float64(270)
	near := float64(0.1)
	far := float64(1000)
	// projection matrix
	P := primitives.ProjectionMatrix4x4(angelOfView, near, far)
	// scalematrix - coord / 100
	scaleMatrix := primitives.ScaleMatrix4x4(0.01, 0.01, 0.01)
	// translation matrix
	translationMatrix := primitives.TranslationMatrix4x4(-1, -1, -50)
	// rotationmatrix - rotate on the Z coord.
	rotationMatrix := primitives.RotationZMatrix4x4(90)
	MV := (scaleMatrix.Dot(translationMatrix)).Dot(rotationMatrix)
	mvpPoints := (MV.Dot(P)).Points

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// time
		elapsedTimeNano := time.Now().UnixNano() - nowUnix
		time := gl.GetUniformLocation(program, gl.Str("time\x00"))
		gl.Uniform1f(time, float32(elapsedTimeNano/10000000))
		mvp := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
		gl.UniformMatrix4fv(mvp, 1, false, &mvpPoints[0])

		for _, item := range triangles {
			item.Draw()
		}
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
