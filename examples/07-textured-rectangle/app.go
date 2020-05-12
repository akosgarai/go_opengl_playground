package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

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
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/07-textured-rectangle/vertexshader.vert", "examples/07-textured-rectangle/fragmentshader.frag")
	shaderProgram.AddTexture("examples/07-textured-rectangle/image-texture.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "textureOne")
	GenerateSquare(shaderProgram)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)
	wrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
