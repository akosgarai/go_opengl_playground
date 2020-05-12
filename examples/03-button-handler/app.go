package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	WindowTitle  = "Example - static button handler"
	epsilon      = 30
	speed        = float32(0.15) / float32(1000.0)

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	SWAP     = glfw.KeyT
	DEBUG    = glfw.KeyH
)

var (
	triangleCoordinates = [3]mgl32.Vec3{
		mgl32.Vec3{-0.75, 0.75, 0}, // top
		mgl32.Vec3{-0.75, 0.25, 0}, // left
		mgl32.Vec3{-0.25, 0.25, 0}, // right
	}
	triangleColors = [3]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0}, // top
		mgl32.Vec3{0, 1, 0}, // left
		mgl32.Vec3{0, 1, 0}, // right
	}
	squareCoordinates = [4]mgl32.Vec3{
		mgl32.Vec3{0.25, -0.25, 0}, // top-left
		mgl32.Vec3{0.25, -0.75, 0}, // bottom-left
		mgl32.Vec3{0.75, -0.75, 0}, // bottom-right
		mgl32.Vec3{0.75, -0.25, 0}, // top-right
	}
	squareColors = [4]mgl32.Vec3{
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	lastUpdate = time.Now().UnixNano() / 1000000

	app *application.Application

	item   *triangle.Triangle
	square *rectangle.Rectangle
)

func SetupKeyMap() map[glfw.Key]bool {
	keyDowns := make(map[glfw.Key]bool)
	keyDowns[FORWARD] = false
	keyDowns[LEFT] = false
	keyDowns[RIGHT] = false
	keyDowns[BACKWARD] = false
	keyDowns[SWAP] = false
	keyDowns[DEBUG] = false

	return keyDowns
}
func Update() {
	if app.GetKeyState(SWAP) {
		item.SetColor(mgl32.Vec3{0, 1, 0})
		square.SetColor(mgl32.Vec3{1, 0, 0})
	} else {
		item.SetColor(mgl32.Vec3{1, 0, 0})
		square.SetColor(mgl32.Vec3{0, 1, 0})
	}
	// now in milisec.
	nowUnixM := time.Now().UnixNano() / 1000000
	delta := nowUnixM - lastUpdate
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		square.SetIndexDirection(1, 1.0)
		item.SetIndexDirection(1, -1.0)
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		square.SetIndexDirection(1, -1.0)
		item.SetIndexDirection(1, 1.0)
	} else {
		square.SetIndexDirection(1, 0.0)
		item.SetIndexDirection(1, 0.0)
	}
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		square.SetIndexDirection(0, -1.0)
		item.SetIndexDirection(0, 1.0)
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		square.SetIndexDirection(0, 1.0)
		item.SetIndexDirection(0, -1.0)
	} else {
		square.SetIndexDirection(0, 0.0)
		item.SetIndexDirection(0, 0.0)
	}
	if epsilon > delta {
		return
	}
	lastUpdate = nowUnixM
	app.Update(float64(delta))
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	app.SetKeys(SetupKeyMap())
	shaderProgram := shader.NewShader("examples/03-button-handler/vertexshader.vert", "examples/03-button-handler/fragmentshader.frag")

	item = triangle.New(triangleCoordinates, triangleColors, shaderProgram)
	item.SetSpeed(speed)
	app.AddItem(item)
	square = rectangle.New(squareCoordinates, squareColors, shaderProgram)
	square.SetSpeed(speed)
	app.AddItem(square)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
