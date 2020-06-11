package main

import (
	"runtime"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/primitives/triangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

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

	glWrapper glwrapper.Wrapper
)

func GenerateColoredRectangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := rectangle.NewSquare()
	v, i, _ := square.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}
func GenerateColoredTriangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := triangle.New(30, 60, 90)
	v, i, _ := square.ColoredMeshInput(col)
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
	triang.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	triang.SetPosition(mgl32.Vec3{-0.4, 0.2, 0})
	mod.AddMesh(triang)

	square := GenerateColoredRectangleMesh(color)
	square.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	square.SetPosition(mgl32.Vec3{0.4, -0.2, 0})
	mod.AddMesh(square)
	mod.RotateX(90)
	app.AddModelToShader(mod, shaderProgram)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
