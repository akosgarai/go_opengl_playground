package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - static triangle and square"
)

var (
	app *application.Application

	color = []mgl32.Vec3{mgl32.Vec3{0, 1, 0}}

	glWrapper wrapper.Wrapper
)

func GenerateColoredRectangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := rectangle.NewSquare()
	v, i := square.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}
func GenerateColoredTriangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := triangle.New(30, 60, 90)
	v, i := square.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}
func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-static-multiple-objects/shaders/vertexshader.vert", "examples/02-static-multiple-objects/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	mod := model.New()
	triang := GenerateColoredTriangleMesh(color)
	triang.SetRotationAngle(mgl32.DegToRad(90))
	triang.SetRotationAxis(mgl32.Vec3{1, 1, 0})
	triang.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	triang.SetPosition(mgl32.Vec3{-0.4, 0.2, 0})
	mod.AddMesh(triang)

	square := GenerateColoredRectangleMesh(color)
	square.SetRotationAngle(mgl32.DegToRad(90))
	square.SetRotationAxis(mgl32.Vec3{1, 0, 0})
	square.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	square.SetPosition(mgl32.Vec3{0.4, -0.2, 0})
	mod.AddMesh(square)
	app.AddModelToShader(mod, shaderProgram)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
