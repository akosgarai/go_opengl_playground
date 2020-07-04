package main

import (
	"path"
	"runtime"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

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

	glWrapper glwrapper.Wrapper
)

// It generates a square.
func GenerateSquareMesh(t texture.Textures) *mesh.TexturedColoredMesh {
	square := rectangle.NewSquare()
	v, i, _ := square.TexturedColoredMeshInput(SquareColor)
	return mesh.NewTexturedColoredMesh(v, i, t, SquareColor, glWrapper)
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/image-texture.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "textureOne", glWrapper)
	squareMesh := GenerateSquareMesh(tex)
	squareMesh.RotateX(90)
	Model.AddMesh(squareMesh)
	scrn.AddModelToShader(Model, shaderProgram)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
