package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Example - static triangle"
)

var (
	coordinates = [3]mgl32.Vec3{
		mgl32.Vec3{-0.75, 0.75, 0}, // top
		mgl32.Vec3{-0.75, 0.25, 0}, // left
		mgl32.Vec3{-0.25, 0.25, 0}, // right
	}
	colors = [3]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0}, // top
		mgl32.Vec3{0, 1, 0}, // left
		mgl32.Vec3{0, 1, 0}, // right
	}

	app *application.Application
)

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-static-triangle/vertexshader.vert", "examples/02-static-triangle/fragmentshader.frag")

	item := triangle.New(coordinates, colors, shaderProgram)
	app.AddItem(item)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
