package main

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/window"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Hello World Window"
)

var (
	ClearColor = [4]float32{0.3, 0.3, 0.3, 1.0}
)

func main() {
	runtime.LockOSThread()

	app := application.New()

	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	program := wrapper.CreateProgram()
	wrapper.LinkProgram(program)
	wrapper.UseProgram(program)
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(window.DummyKeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(window.DummyMouseButtonCallback)

	wrapper.ClearColor(ClearColor[0], ClearColor[1], ClearColor[2], ClearColor[3])

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
