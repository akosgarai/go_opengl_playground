package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth     = 800
	WindowHeight    = 800
	WindowTitle     = "Example - the house"
	CameraMoveSpeed = 1.0 / 100.0
	Epsilon         = 100.0
)

var (
	cameraLastUpdate int64
	app              *application.Application

	glWrapper glwrapper.Wrapper

	Model = model.New()
)

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
	camera := camera.NewCamera(mgl32.Vec3{75, 30, 0.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 1000.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(90)
	return camera
}

// the path
func Path() {
	rect := rectangle.New(30, 50)
	col := []mgl32.Vec3{mgl32.Vec3{215.0 / 255.0, 100.0 / 255.0, 30.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{75, 0, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the wall left of the path
func LeftFullWall() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{90, 25, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateZ(90)

	Model.AddMesh(m)
}

// the wall front of the path
func FrontPathWall() {
	rect := rectangle.New(30, 50)
	col := []mgl32.Vec3{mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{75, 25, 80})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateX(90)
	Model.AddMesh(m)
}

// the roof of the path
func PathRoof() {
	rect := rectangle.New(30, 50)
	col := []mgl32.Vec3{mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{75, 50, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	Model.AddMesh(m)
}

// the floor of the room
func RoomFloor() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 0, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the roof of the room
func RoomRoof() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 50, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the front wall of the room
func RoomFront() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 25, 80})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateX(90)
	Model.AddMesh(m)
}

// the back wall of the room
func RoomBack() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 25, 30})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateX(90)
	Model.AddMesh(m)
}

// the left wall of the room
func RoomLeft() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i, _ := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{10, 25, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateZ(90)
	Model.AddMesh(m)
}

func Update() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := float64(nowUnix-cameraLastUpdate) / float64(time.Millisecond)

	if Epsilon > delta {
		return
	}
	cameraLastUpdate = nowUnix

	app.Update(delta)
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
}

func main() {
	runtime.LockOSThread()

	app = application.New(glWrapper)

	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)

	scrn.SetupCamera(CreateCamera(), CameraMovementOptions())
	cameraLastUpdate = time.Now().UnixNano()

	Path()
	LeftFullWall()
	FrontPathWall()
	PathRoof()
	RoomFloor()
	RoomRoof()
	RoomFront()
	RoomBack()
	RoomLeft()
	scrn.AddModelToShader(Model, shaderProgram)
	scrn.Setup(setupApp)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)
	glWrapper.Viewport(0, 0, WindowWidth, WindowHeight)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(window.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
