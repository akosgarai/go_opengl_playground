package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - textured rectangle"
)

var (
	app *application.Application
)

// It generates a square.
func GenerateSquare(shaderProgram *shader.Shader) {
	squareColor := mgl32.Vec3{0.58, 0.29, 0}
	coords := [4]mgl32.Vec3{
		mgl32.Vec3{0.5, 0.5, 0},
		mgl32.Vec3{0.5, -0.5, 0},
		mgl32.Vec3{-0.5, -0.5, 0},
		mgl32.Vec3{-0.5, 0.5, 0},
	}
	colors := [4]mgl32.Vec3{squareColor, squareColor, squareColor, squareColor}
	square := rectangle.New(coords, colors, shaderProgram)
	square.SetPrecision(1)
	app.AddItem(square)
}

// Update the z coordinates of the vectors.
func Update() {
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	shaderProgram := shader.NewShader("examples/07-textured-rectangle/vertexshader.vert", "examples/07-textured-rectangle/fragmentshader.frag")
	shaderProgram.AddTexture("examples/07-textured-rectangle/image-texture.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	GenerateSquare(shaderProgram)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
