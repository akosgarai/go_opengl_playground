package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - static triangle"
)

var (
	color = []mgl32.Vec3{mgl32.Vec3{0, 1, 0}}

	app *application.Application
)

func GenerateColoredMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	triang := triangle.New(60, 60, 60)
	v, i := triang.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-static-triangle/shaders/vertexshader.vert", "examples/02-static-triangle/shaders/fragmentshader.frag")
	app.AddShader(shaderProgram)

	triang := GenerateColoredMesh(color)
	triang.SetRotationAngle(mgl32.DegToRad(90))
	triang.SetRotationAxis(mgl32.Vec3{1, 0, 0})
	app.AddMeshToShader(triang, shaderProgram)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
