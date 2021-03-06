package main

import (
	"math/rand"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/transformations"
	"github.com/akosgarai/playground_engine/pkg/window"

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

	glWrapper glwrapper.Wrapper
)

func init() {
	runtime.LockOSThread()
}

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["forward"] = []glfw.Key{glfw.KeyW}
	cm["back"] = []glfw.Key{glfw.KeyS}
	cm["rotateLeft"] = []glfw.Key{glfw.KeyA}
	cm["rotateRight"] = []glfw.Key{glfw.KeyD}
	cm["rotateOnEdgeDistance"] = float32(0.0)
	cm["mode"] = "default"
	return cm
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.DefaultCamera {
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
		mX, mY := transformations.MouseCoordinates(app.MousePosX, app.MousePosY, WindowWidth, WindowHeight)
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
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.PROGRAM_POINT_SIZE)
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
}

func main() {
	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	scrn.SetupCamera(CreateCamera(), CameraMovementOptions())

	cameraLastUpdate = time.Now().UnixNano()

	Shader := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(Shader)

	PointMesh = mesh.NewPointMesh(glWrapper)
	Model.AddMesh(PointMesh)
	scrn.AddModelToShader(Model, Shader)
	scrn.Setup(setupApp)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
