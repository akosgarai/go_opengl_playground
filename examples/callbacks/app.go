package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/akosgarai/opengl_playground/pkg/application"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - callbacks"
)

func main() {
	runtime.LockOSThread()

	app := application.New()

	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	program := gl.CreateProgram()
	gl.LinkProgram(program)
	gl.UseProgram(program)
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.DummyKeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(app.DummyMouseButtonCallback)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
