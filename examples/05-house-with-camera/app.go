package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
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
	epsilon      = 1000.0
	precision    = 20
	// buttons
	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Turn 90 deg. left
	RIGHT    = glfw.KeyD // Turn 90 deg. right
)

var (
	cameraLastUpdate int64
	app              *application.Application
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{75, 30, 0.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 1000.0)
	return camera
}

// Create the keymap
func SetupKeyMap() map[glfw.Key]bool {
	keyDowns := make(map[glfw.Key]bool)
	keyDowns[FORWARD] = false
	keyDowns[LEFT] = false
	keyDowns[RIGHT] = false
	keyDowns[BACKWARD] = false

	return keyDowns
}
func addRect(coordinates [4]mgl32.Vec3, rectColor mgl32.Vec3, shaderProg *shader.Shader) {
	color := [4]mgl32.Vec3{rectColor, rectColor, rectColor, rectColor}
	rect := rectangle.New(coordinates, color, shaderProg)
	rect.SetPrecision(precision)
	app.AddItem(rect)
}

// the path
func Path(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{215.0 / 255.0, 100.0 / 255.0, 30.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 0, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the wall left of the path
func LeftFullWall(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{90, 50, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 0, 30},
		mgl32.Vec3{90, 50, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the wall front of the path
func FrontPathWall(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 50, 80},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the roof of the path
func PathRoof(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 50, 30},
		mgl32.Vec3{90, 50, 30},
		mgl32.Vec3{90, 50, 80},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the floor of the room
func RoomFloor(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{60, 0, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the roof of the room
func RoomRoof(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 50, 30},
		mgl32.Vec3{10, 50, 80},
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 50, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the front wall of the room
func RoomFront(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{10, 50, 80},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the back wall of the room
func RoomBack(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 30},
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 50, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the left wall of the room
func RoomLeft(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 50, 30},
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{10, 50, 80},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the right wall of the room 1x5
func RoomRight1(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{60, 0, 70},
		mgl32.Vec3{60, 50, 70},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the right wall of the room 2x5
func RoomRight2(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 50},
		mgl32.Vec3{60, 0, 50},
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{60, 50, 30},
	}
	addRect(coordinates, rectColor, shaderProg)
}

// the right wall of the room 2x2
func RoomRight3(shaderProg *shader.Shader) {
	rectColor := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 70},
		mgl32.Vec3{60, 30, 70},
		mgl32.Vec3{60, 30, 50},
		mgl32.Vec3{60, 50, 50},
	}
	addRect(coordinates, rectColor, shaderProg)
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
		dX = -90
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		dX = 90
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
	wrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/05-house-with-camera/vertexshader.vert", "examples/05-house-with-camera/fragmentshader.frag")

	app.SetCamera(CreateCamera())
	cameraLastUpdate = time.Now().UnixNano()

	app.SetKeys(SetupKeyMap())
	Path(shaderProgram)
	LeftFullWall(shaderProgram)
	FrontPathWall(shaderProgram)
	PathRoof(shaderProgram)
	RoomFloor(shaderProgram)
	RoomRoof(shaderProgram)
	RoomFront(shaderProgram)
	RoomBack(shaderProgram)
	RoomLeft(shaderProgram)
	RoomRight1(shaderProgram)
	RoomRight2(shaderProgram)
	RoomRight3(shaderProgram)

	wrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
	wrapper.Viewport(0, 0, WindowWidth, WindowHeight)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(window.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
