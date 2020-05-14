package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/application"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/mesh"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/shader"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/vertex"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - mesh experiment"

	RED   = glfw.KeyR // red color component
	GREEN = glfw.KeyG // green color component
	BLUE  = glfw.KeyB // blue color component

	// buttons
	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Turn 90 deg. left
	RIGHT    = glfw.KeyD // Turn 90 deg. right

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft

	moveSpeed = 1.0 / 100.0
	epsilon   = 100.0
)

var (
	app *application.Application

	Shader    *shader.Shader
	PointMesh *mesh.PointMesh

	cameraLastUpdate int64

	addPoint = false
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -10.0}, mgl32.Vec3{0, 1, 0}, -270.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}
func Update() {
	updatePointState()
	updateCameraState()
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
		mX, mY := trans.MouseCoordinates(app.MousePosX, app.MousePosY, WindowWidth, WindowHeight)
		// to calculate the coordinate of the point, we have to apply the inverse of the camera transformations.
		V := app.GetCamera().GetViewMatrix()
		P := app.GetCamera().GetProjectionMatrix()
		trMat := P.Mul4(V).Inv()
		coords := mgl32.TransformCoordinate(mgl32.Vec3{float32(mX), float32(mY), 0.0}, trMat)
		color := mgl32.Vec3{r, g, b}
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
	wrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	cameraLastUpdate = time.Now().UnixNano()

	Shader = shader.NewShader("examples/model-loading/shaders/point.vert", "examples/model-loading/shaders/point.frag")
	app.AddShader(Shader)

	PointMesh = mesh.NewPointMesh()
	app.AddMeshToShader(PointMesh, Shader)

	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	wrapper.Enable(wrapper.PROGRAM_POINT_SIZE)
	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)
	wrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
