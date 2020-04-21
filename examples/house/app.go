package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/house/application"
	"github.com/akosgarai/opengl_playground/examples/house/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - the house"
	moveSpeed    = 1.0 / 100.0
	epsilon      = 1000.0
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

// It creates the shader application for the items.
func CreateShader() uint32 {
	vertexShader, err := shader.CompileShader(shader.VertexShaderModelViewProjectionSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderBasicSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

// It creates a new camera with the necessary setup
func CreateCamera() *primitives.Camera {
	camera := primitives.NewCamera(mgl32.Vec3{75, 30, 0.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 100.0)
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

// the path
func Path(shaderProgId uint32) {
	floorCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 0, 30},
	}
	floorColor := mgl32.Vec3{215.0 / 255.0, 100.0 / 255.0, 30.0 / 255.0}
	app.AddItem(primitives.NewRectangle(floorCoordinates, floorColor, 20, shaderProgId))
}

// the wall left of the path
func LeftFullWall(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{90, 50, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 0, 30},
		mgl32.Vec3{90, 50, 30},
	}
	color := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the wall front of the path
func FrontPathWall(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{90, 0, 80},
		mgl32.Vec3{90, 50, 80},
	}
	color := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the roof of the path
func PathRoof(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 50, 30},
		mgl32.Vec3{90, 50, 30},
		mgl32.Vec3{90, 50, 80},
	}
	color := mgl32.Vec3{165.0 / 255.0, 42.0 / 255.0, 42.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the floor of the room
func RoomFloor(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{60, 0, 30},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the roof of the room
func RoomRoof(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 50, 30},
		mgl32.Vec3{10, 50, 80},
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 50, 30},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the front wall of the room
func RoomFront(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{10, 50, 80},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the back wall of the room
func RoomBack(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 30},
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 50, 30},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the left wall of the room
func RoomLeft(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{10, 50, 30},
		mgl32.Vec3{10, 0, 30},
		mgl32.Vec3{10, 0, 80},
		mgl32.Vec3{10, 50, 80},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the right wall of the room 1x5
func RoomRight1(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 80},
		mgl32.Vec3{60, 0, 80},
		mgl32.Vec3{60, 0, 70},
		mgl32.Vec3{60, 50, 70},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the right wall of the room 2x5
func RoomRight2(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 50},
		mgl32.Vec3{60, 0, 50},
		mgl32.Vec3{60, 0, 30},
		mgl32.Vec3{60, 50, 30},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// the right wall of the room 2x2
func RoomRight3(shaderProgramId uint32) {
	coordinates := [4]mgl32.Vec3{
		mgl32.Vec3{60, 50, 70},
		mgl32.Vec3{60, 30, 70},
		mgl32.Vec3{60, 30, 50},
		mgl32.Vec3{60, 50, 50},
	}
	color := mgl32.Vec3{196.0 / 255.0, 196.0 / 255.0, 196.0 / 255.0}
	app.AddItem(primitives.NewRectangle(coordinates, color, 20, shaderProgramId))
}

// The Sphere.
func Sphere(shaderProgId uint32) {
	sphere := primitives.NewSphere()
	sphere.SetCenter(mgl32.Vec3{75, 5, 35})
	sphere.SetColor(mgl32.Vec3{0, 0, 1})
	sphere.SetRadius(1)
	sphere.SetShaderProgram(shaderProgId)
	app.AddItem(sphere)
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

	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	shaderProgramId := CreateShader()

	app.SetCamera(CreateCamera())
	cameraLastUpdate = time.Now().UnixNano()

	app.SetKeys(SetupKeyMap())
	Path(shaderProgramId)
	LeftFullWall(shaderProgramId)
	FrontPathWall(shaderProgramId)
	PathRoof(shaderProgramId)
	RoomFloor(shaderProgramId)
	RoomRoof(shaderProgramId)
	RoomFront(shaderProgramId)
	RoomBack(shaderProgramId)
	RoomLeft(shaderProgramId)
	RoomRight1(shaderProgramId)
	RoomRight2(shaderProgramId)
	RoomRight3(shaderProgramId)
	//Sphere(shaderProgramId)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(app.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
