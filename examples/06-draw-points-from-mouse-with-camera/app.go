package main

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
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
	WindowTitle  = "Example - draw points from mouse inputs and keyboard colors"

	// buttons
	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
	RED               = glfw.KeyR // red color component
	GREEN             = glfw.KeyG // green color component
	BLUE              = glfw.KeyB // blue color component

	CameraMoveSpeed = 1.0 / 100.0
	Epsilon         = 100.0
)

var (
	app *application.Application

	Shader    *shader.Shader
	PointMesh *mesh.PointMesh

	cameraLastUpdate int64

	addPoint = false
	Model    = model.New()

	glWrapper wrapper.Wrapper
)

// Setup keymap for the camera movement
func CameraMovementMap() map[string]glfw.Key {
	cm := make(map[string]glfw.Key)
	cm["forward"] = glfw.KeyW
	cm["back"] = glfw.KeyS
	cm["rotateLeft"] = glfw.KeyA
	cm["rotateRight"] = glfw.KeyD
	return cm
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -10.0}, mgl32.Vec3{0, -1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(90)
	return camera
}
func Update() {
	updatePointState()
	nowUnix := time.Now().UnixNano()
	delta := float64(nowUnix-cameraLastUpdate) / float64(time.Millisecond)
	if Epsilon > delta {
		return
	}
	cameraLastUpdate = nowUnix
	app.Update(delta)
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
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())
	app.SetCameraMovementMap(CameraMovementMap())

	cameraLastUpdate = time.Now().UnixNano()

	Shader = shader.NewShader("examples/06-draw-points-from-mouse-with-camera/shaders/vertexshader.vert", "examples/06-draw-points-from-mouse-with-camera/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(Shader)

	PointMesh = mesh.NewPointMesh(glWrapper)
	Model.AddMesh(PointMesh)
	app.AddModelToShader(Model, Shader)

	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	glWrapper.Enable(wrapper.PROGRAM_POINT_SIZE)
	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
