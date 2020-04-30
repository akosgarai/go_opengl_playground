package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Example - static square"
)

var (
	coordinates = [4]mgl32.Vec3{
		mgl32.Vec3{0.25, -0.25, 0}, // top-left
		mgl32.Vec3{0.25, -0.75, 0}, // bottom-left
		mgl32.Vec3{0.75, -0.75, 0}, // bottom-right
		mgl32.Vec3{0.75, -0.25, 0}, // top-right
	}
	colors = [4]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}

	app *application.Application
)

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	shaderProgram := shader.NewShader("examples/static-square/vertexshader.vert", "examples/static-square/fragmentshader.frag")

	square := rectangle.New(coordinates, colors, shaderProgram)
	app.AddItem(square)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
