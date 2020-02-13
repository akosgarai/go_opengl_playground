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
	windowTitle  = "Example - static triangle"
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
func generateTrianglesModelCoordinates() {
	/*
	 * The goal is to draw triangles to the screen. The screen will contain 20 * 20 triangles.
	 * one part : 10 width,
	 */
	rows := 100
	cols := 100
	length := 10.0
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := (float64(j) * length)
			topY := (float64(i) * length)
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

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	generateTrianglesModelCoordinates()

	nowUnix := time.Now().UnixNano()

	// mvp - modelview - projection matrix
	angelOfView := float64(270)
	near := float64(0.1)
	far := float64(100)
	P := primitives.ProjectionMatrix4x4(angelOfView, near, far)
	scaleMatrix := primitives.ScaleMatrix4x4(0.01, 0.01, 0.01)
	translationMatrix := primitives.TranslationMatrix4x4(-1, -1, -100)
	rotationMatrix := primitives.RotationZMatrix4x4(90)
	MV := (scaleMatrix.Dot(translationMatrix)).Dot(rotationMatrix)
	mvpPoints := (P.Dot(MV)).Points

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// time
		elapsedTimeNano := time.Now().UnixNano() - nowUnix
		time := gl.GetUniformLocation(program, gl.Str("time\x00"))
		gl.Uniform1f(time, float32(elapsedTimeNano/10000000))
		/*
		 * MPV = P * M, where P is the projection matrix and M is the model-view matrix (aka. object to world space * world space to camera space)
		 * P supposed to be something like this : [16]float32{
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		 }
		 * MV supposed to be something like this : [16]float32{
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		     0.0,0.0,0.0,0.0,
		 }
		 * https://www.scratchapixel.com/lessons/3d-basic-rendering/perspective-and-orthographic-projection-matrix/projection-matrix-introduction
		*/
		mvp := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
		gl.UniformMatrix4fv(mvp, 1, false, &mvpPoints[0])

		for _, item := range triangles {
			item.Draw()
		}
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
