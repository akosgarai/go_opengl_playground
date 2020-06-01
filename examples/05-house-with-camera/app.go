package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - the house"
	moveSpeed    = 1.0 / 100.0
	epsilon      = 100.0
	// buttons
	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Turn 90 deg. left
	RIGHT    = glfw.KeyD // Turn 90 deg. right
)

var (
	cameraLastUpdate int64
	app              *application.Application

	glWrapper wrapper.Wrapper

	Model = model.New()
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{75, 30, 0.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 1000.0)
	return camera
}

// the path
func Path() {
	rect := rectangle.New(30, 50)
	col := []mgl32.Vec3{mgl32.Vec3{215.0 / 255.0, 100.0 / 255.0, 30.0 / 255.0}}
	v, i := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{75, 0, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the wall left of the path
func LeftFullWall() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}}
	v, i := rect.ColoredMeshInput(col)
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
	v, i := rect.ColoredMeshInput(col)
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
	v, i := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{75, 50, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	Model.AddMesh(m)
}

// the floor of the room
func RoomFloor() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 0, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the roof of the room
func RoomRoof() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{35, 50, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})

	Model.AddMesh(m)
}

// the front wall of the room
func RoomFront() {
	rect := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}}
	v, i := rect.ColoredMeshInput(col)
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
	v, i := rect.ColoredMeshInput(col)
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
	v, i := rect.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(mgl32.Vec3{10, 25, 55})
	m.SetScale(mgl32.Vec3{50, 50, 50})
	m.RotateZ(90)
	Model.AddMesh(m)
}

func Update() {
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
		dX = -9
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		dX = 9
	}
	if dX != 0.0 {
		app.GetCamera().UpdateDirection(dX, dY)
	}
}
func main() {
	runtime.LockOSThread()

	app = application.New()

	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/05-house-with-camera/shaders/vertexshader.vert", "examples/05-house-with-camera/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	app.SetCamera(CreateCamera())
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
	app.AddModelToShader(Model, shaderProgram)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
	glWrapper.Viewport(0, 0, WindowWidth, WindowHeight)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(window.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
