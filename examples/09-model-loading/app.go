package main

import (
	"os"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/modelimport"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Model loading example"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	CameraMoveSpeed       = 0.005
	CameraDirectionSpeed  = float32(0.00500)
	DefaultModelDirectory = "examples/09-model-loading/assets"
	DefaultModelFilename  = "object.obj"
)

var (
	app      *application.Application
	Importer *modelimport.Import

	lastUpdate int64

	cameraDistance = 0.1

	glWrapper wrapper.Wrapper

	Model = model.New()
)

// Importer init. If we have 2 or more command line arguments,
// the first one is used as model directory, the second one
// as the model filename.
func Init() {
	args := os.Args[1:]
	if len(args) == 0 {
		Importer = modelimport.New(DefaultModelDirectory, DefaultModelFilename)
	} else if len(args) == 1 {
		Importer = modelimport.New(args[0], DefaultModelFilename)
	} else if len(args) > 1 {
		Importer = modelimport.New(args[0], args[1])
	}
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -5.0}, mgl32.Vec3{0, 1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	app.Update(moveTime)

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = CameraMoveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -CameraMoveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	horizontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horizontal = -CameraMoveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horizontal = CameraMoveSpeed * moveTime
	}
	if horizontal != 0 {
		app.GetCamera().Strafe(float32(horizontal))
	}
	vertical := 0.0
	if app.GetKeyState(UP) && !app.GetKeyState(DOWN) {
		vertical = -CameraMoveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = CameraMoveSpeed * moveTime
	}
	if vertical != 0 {
		app.GetCamera().Lift(float32(vertical))
	}
	currX, currY := app.GetWindow().GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, WindowWidth, WindowHeight)
	KeyDowns := make(map[string]bool)
	// dUp
	if y > 1.0-cameraDistance && y < 1.0 {
		KeyDowns["dUp"] = true
	} else {
		KeyDowns["dUp"] = false
	}
	// dDown
	if y < -1.0+cameraDistance && y > -1.0 {
		KeyDowns["dDown"] = true
	} else {
		KeyDowns["dDown"] = false
	}
	// dLeft
	if x < -1.0+cameraDistance && x > -1.0 {
		KeyDowns["dLeft"] = true
	} else {
		KeyDowns["dLeft"] = false
	}
	// dRight
	if x > 1.0-cameraDistance && x < 1.0 {
		KeyDowns["dRight"] = true
	} else {
		KeyDowns["dRight"] = false
	}

	dX := float32(0.0)
	dY := float32(0.0)
	if KeyDowns["dUp"] && !KeyDowns["dDown"] {
		dY = CameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = -CameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = -CameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = CameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}

func main() {
	Init()
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/09-model-loading/shaders/vertexshader.vert", "examples/09-model-loading/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	lightSource := light.NewPointLight([4]mgl32.Vec3{mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	app.AddPointLightSource(lightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})
	Importer.Import()
	meshes := Importer.GetMeshes()
	for _, m := range meshes {
		Model.AddMesh(m)
	}
	app.AddModelToShader(Model, shaderProgram)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	lastUpdate = time.Now().UnixNano()
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
