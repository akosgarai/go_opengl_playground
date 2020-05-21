package main

import (
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/texture"
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

	SquareColor = []mgl32.Vec3{mgl32.Vec3{0.58, 0.29, 0}}
	Model       = model.New()

	glWrapper wrapper.Wrapper
)

// It generates a square.
func GenerateSquareMesh(t texture.Textures) *mesh.TexturedColoredMesh {
	square := rectangle.NewSquare()
	v, i := square.TexturedColoredMeshInput(SquareColor)
	return mesh.NewTexturedColoredMesh(v, i, t, SquareColor, glWrapper)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/07-textured-rectangle/shaders/vertexshader.vert", "examples/07-textured-rectangle/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	var tex texture.Textures
	tex.AddTexture("examples/07-textured-rectangle/assets/image-texture.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "textureOne", glWrapper)
	squareMesh := GenerateSquareMesh(tex)
	squareMesh.SetRotationAngle(mgl32.DegToRad(90))
	squareMesh.SetRotationAxis(mgl32.Vec3{1, 0, 0})
	Model.AddMesh(squareMesh)
	app.AddModelToShader(Model, shaderProgram)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
