package main

import (
	"math/rand"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/point"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs"

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

var (
	addPoint = false

	app    *application.Application
	points *point.Points
)

func Update() {
	if !app.GetMouseButtonState(LEFT_MOUSE_BUTTON) && addPoint {
		mX, mY := trans.MouseCoordinates(app.MousePosX, app.MousePosY, windowWidth, windowHeight)
		coords := mgl32.Vec3{float32(mX), float32(mY), 0.0}
		color := mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
		size := float32(3 + rand.Intn(17))
		points.Add(coords, color, size)
		addPoint = false
	} else if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
		addPoint = true
	}
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	shaderProgram := shader.NewShader("examples/06-draw-points-from-mouse-input/vertexshader.vert", "examples/06-draw-points-from-mouse-input/fragmentshader.frag")
	points = point.New(shaderProgram)
	app.AddItem(points)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		if points.Count() > 0 {
			app.Draw()
		}
		app.GetWindow().SwapBuffers()
	}
}
