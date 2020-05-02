package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
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
	windowTitle  = "Example - draw points from mouse inputs and keyboard colors"
)

var (
	addPoint = false

	cameraLastUpdate int64
	moveSpeed        = 1.0 / 100.0
	epsilon          = 100.0

	app    *application.Application
	points *point.Points

	RED   = glfw.KeyR // red color component
	GREEN = glfw.KeyG // green color component
	BLUE  = glfw.KeyB // blue color component

	// buttons
	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Turn 90 deg. left
	RIGHT    = glfw.KeyD // Turn 90 deg. right

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -10.0}, mgl32.Vec3{0, 1, 0}, -270.0, 0.0)
	camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 100.0)
	return camera
}

func updatePointState() {
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
		// to calculate the coordinate of the point, we have to apply the inverse of the camera transformations.
		V := app.GetCamera().GetViewMatrix()
		P := app.GetCamera().GetProjectionMatrix()
		trMat := P.Mul4(V).Inv()
		coords := mgl32.TransformCoordinate(mgl32.Vec3{float32(mX), float32(mY), 0.0}, trMat)
		color := mgl32.Vec3{r, g, b}
		size := float32(3 + rand.Intn(17))
		points.Add(coords, color, size)
		addPoint = false
	} else if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
		addPoint = true
	}
}
func updateCameraState() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))

	if epsilon > moveTime {
		return
	}
	cameraLastUpdate = nowUnix

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = moveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -moveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	dX := float32(0.0)
	dY := float32(0.0)
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		dX = -90
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		dX = 90
	}
	if dX != 0.0 {
		app.GetCamera().UpdateDirection(dX, dY)
	}
}
func Update() {
	updatePointState()
	updateCameraState()
}
func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()
	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/06-draw-points-from-mouse-with-camera/vertexshader.vert", "examples/06-draw-points-from-mouse-with-camera/fragmentshader.frag")
	points = point.New(shaderProgram)
	app.AddItem(points)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		if points.Count() > 0 {
			app.DrawWithUniforms()
		}
		app.GetWindow().SwapBuffers()
	}
}
