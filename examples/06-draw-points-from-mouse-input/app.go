package main

import (
	"math/rand"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - draw points from mouse inputs"

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

var (
	addPoint = false
	Model    = model.New()

	app *application.Application

	PointMesh *mesh.PointMesh

	glWrapper wrapper.Wrapper
)

func Update() {
	if !app.GetMouseButtonState(LEFT_MOUSE_BUTTON) && addPoint {
		mX, mY := trans.MouseCoordinates(app.MousePosX, app.MousePosY, WindowWidth, WindowHeight)
		coords := mgl32.Vec3{float32(mX), float32(mY), 0.0}
		color := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
		size := float32(3 + rand.Intn(17))
		vert := vertex.Vertex{
			Position:  coords,
			Color:     color,
			PointSize: size,
		}
		PointMesh.AddVertex(vert)

		addPoint = false
	} else if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
		addPoint = true
	}
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/06-draw-points-from-mouse-input/shaders/vertexshader.vert", "examples/06-draw-points-from-mouse-input/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	PointMesh = mesh.NewPointMesh(glWrapper)
	Model.AddMesh(PointMesh)
	app.AddModelToShader(Model, shaderProgram)

	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	glWrapper.Enable(wrapper.PROGRAM_POINT_SIZE)
	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
