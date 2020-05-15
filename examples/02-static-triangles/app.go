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
	WindowHeight = 600
	WindowTitle  = "Example - static triangles, lots of them"
)

var (
	app *application.Application

	color = mgl32.Vec3{0, 1, 0}
)

// GenerateTriangles fills up the triangles.
func GenerateTriangles(rows int, shaderProgram *shader.Shader) {
	triang := triangle.New(60, 60, 60)
	v, indicies := triang.ColoredMeshInput(color)

	length := 2.0 / float32(rows)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= rows; j++ {
			topX := 1.0 - (float32(j) * length)
			topY := -1.0 + (float32(i) * length)
			topZ := float32(0.0)

			m := mesh.NewColorMesh(v, indicies)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetRotationAngle(mgl32.DegToRad(90))
			m.SetRotationAxis(mgl32.Vec3{1, 0, 0})
			m.SetScale(mgl32.Vec3{length, length, length})
			app.AddMeshToShader(m, shaderProgram)

		}
	}
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/02-static-triangles/shaders/vertexshader.vert", "examples/02-static-triangles/shaders/fragmentshader.frag")
	app.AddShader(shaderProgram)

	GenerateTriangles(50, shaderProgram)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)

		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
