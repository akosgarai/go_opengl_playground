package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Static triangle with multiple color"
)

var (
	color = []mgl32.Vec3{
		mgl32.Vec3{0, 1, 0}, // top
		mgl32.Vec3{1, 0, 0}, // left
		mgl32.Vec3{0, 0, 1}, // right
	}

	app *application.Application

	glWrapper wrapper.Wrapper
)

func GenerateColoredTriangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	triang := triangle.New(30, 60, 90)
	v, i := triang.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-multiple-color-triangle/shaders/vertexshader.vert", "examples/02-multiple-color-triangle/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	triang := GenerateColoredTriangleMesh(color)
	triang.SetRotationAngle(mgl32.DegToRad(90))
	triang.SetRotationAxis(mgl32.Vec3{1, 1, 0})
	mod := model.New()
	mod.AddMesh(triang)
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
