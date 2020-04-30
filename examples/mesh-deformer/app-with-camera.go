package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	tr "github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - mesh deformer - with moving camera"
	moveSpeed    = 0.005

	rows   = 10
	cols   = 10
	length = 10

	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Go left
	RIGHT    = glfw.KeyD // Go right
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE
)

var (
	app *application.Application

	triangleColorFront = mgl32.Vec3{0, 0, 1}
	triangleColorBack  = mgl32.Vec3{0, 0.5, 1}

	lastUpdate int64

	cameraDistance       = 0.1
	cameraDirectionSpeed = float32(0.500)
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 1.0, 9.5}, mgl32.Vec3{0, 0, 1}, -90.0, -34.0)
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
	keyDowns[UP] = false
	keyDowns[DOWN] = false

	return keyDowns
}

// It generates a bunch of triangles and sets their color to static blue.
func GenerateTriangles(shaderProgram *shader.Shader) {
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j * length)
			topY := float32(i * length)
			topZ := float32(0.0)

			coords := [3]mgl32.Vec3{
				mgl32.Vec3{topX, topY, topZ},
				mgl32.Vec3{topX, topY - float32(length), topZ},
				mgl32.Vec3{topX - float32(length), topY - float32(length), topZ},
			}
			colors := [3]mgl32.Vec3{triangleColorFront, triangleColorFront, triangleColorFront}
			triangle := tr.NewTriangle(coords, colors, shaderProgram)
			triangle.SetDirection(mgl32.Vec3{0, 0, 1})
			triangle.SetSpeed(float32(1.0) / float32(1000000000.0))
			app.AddItem(triangle)

			coords = [3]mgl32.Vec3{
				mgl32.Vec3{topX, topY, topZ},
				mgl32.Vec3{topX - float32(length), topY - float32(length), topZ},
				mgl32.Vec3{topX - float32(length), topY, topZ},
			}
			colors = [3]mgl32.Vec3{triangleColorBack, triangleColorBack, triangleColorBack}
			triangle = tr.NewTriangle(coords, colors, shaderProgram)
			triangle.SetDirection(mgl32.Vec3{0, 0, 1})
			triangle.SetSpeed(float32(1.0) / float32(1000000000.0))
			app.AddItem(triangle)
		}
	}
}

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano - lastUpdate)
	moveTime := delta / float64(time.Millisecond)
	app.Update(delta)
	lastUpdate = nowNano

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = moveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -moveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	horisontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horisontal = -moveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horisontal = moveSpeed * moveTime
	}
	if horisontal != 0 {
		app.GetCamera().Strafe(float32(horisontal))
	}
	vertical := 0.0
	if app.GetKeyState(UP) && !app.GetKeyState(DOWN) {
		vertical = -moveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = moveSpeed * moveTime
	}
	if vertical != 0 {
		app.GetCamera().Lift(float32(vertical))
	}

	currX, currY := app.GetWindow().GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, windowWidth, windowHeight)
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
		dY = cameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = -cameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = -cameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = cameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/mesh-deformer/vertexshader.vert", "examples/mesh-deformer/fragmentshader.frag")

	GenerateTriangles(shaderProgram)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
