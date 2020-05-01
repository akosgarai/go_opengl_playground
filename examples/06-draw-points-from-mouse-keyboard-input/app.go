package main

import (
	"math/rand"
	"runtime"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/point"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs and keyboard colors"
)

var (
	addPoint = false

	app    *application.Application
	points *point.Points

	RED   = glfw.KeyR // red color component
	GREEN = glfw.KeyG // green color component
	BLUE  = glfw.KeyB // blue color component

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

func Update() {
	if !app.GetMouseButtonState(LEFT_MOUSE_BUTTON) && addPoint {
		var r, g, b float32
		if app.GetKeyState(RED) {
			r = 1
		} else {
			r = 0
		}
		if app.GetKeyState(GREEN) {
			g = 1
		} else {
			g = 0
		}
		if app.GetKeyState(BLUE) {
			b = 1
		} else {
			b = 0
		}
		mX, mY := trans.MouseCoordinates(app.MousePosX, app.MousePosY, windowWidth, windowHeight)
		coords := mgl32.Vec3{float32(mX), float32(mY), 0.0}
		color := mgl32.Vec3{r, g, b}
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
	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	shaderProgram := shader.NewShader("examples/06-draw-points-from-mouse-keyboard-input/vertexshader.vert", "examples/06-draw-points-from-mouse-keyboard-input/fragmentshader.frag")
	points = point.New(shaderProgram)
	app.AddItem(points)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

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
