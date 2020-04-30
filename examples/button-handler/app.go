package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	tr "github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
	windowTitle  = "Example - static button handler"
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

	triangle *tr.Triangle
	square   *rectangle.Rectangle
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
		triangle.SetColor(mgl32.Vec3{0, 1, 0})
		square.SetColor(mgl32.Vec3{1, 0, 0})
	} else {
		triangle.SetColor(mgl32.Vec3{1, 0, 0})
		square.SetColor(mgl32.Vec3{0, 1, 0})
	}
	// now in milisec.
	nowUnixM := time.Now().UnixNano() / 1000000
	delta := nowUnixM - lastUpdate
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		square.SetIndexDirection(1, 1.0)
		triangle.SetIndexDirection(1, -1.0)
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		square.SetIndexDirection(1, -1.0)
		triangle.SetIndexDirection(1, 1.0)
	} else {
		square.SetIndexDirection(1, 0.0)
		triangle.SetIndexDirection(1, 0.0)
	}
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		square.SetIndexDirection(0, -1.0)
		triangle.SetIndexDirection(0, 1.0)
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		square.SetIndexDirection(0, 1.0)
		triangle.SetIndexDirection(0, -1.0)
	} else {
		square.SetIndexDirection(0, 0.0)
		triangle.SetIndexDirection(0, 0.0)
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
	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	app.SetKeys(SetupKeyMap())
	shaderProgram := shader.NewShader("examples/button-handler/vertexshader.vert", "examples/button-handler/fragmentshader.frag")

	triangle = tr.NewTriangle(triangleCoordinates, triangleColors, shaderProgram)
	triangle.SetSpeed(speed)
	app.AddItem(triangle)
	square = rectangle.New(squareCoordinates, squareColors, shaderProgram)
	square.SetSpeed(speed)
	app.AddItem(square)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
