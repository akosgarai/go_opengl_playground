package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - mesh deformer - with moving camera"
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

	triangleColorFront = []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	triangleColorBack  = []mgl32.Vec3{mgl32.Vec3{0, 0.5, 1}}

	lastUpdate int64

	cameraDistance       = 0.1
	cameraDirectionSpeed = float32(0.005)
	Model                = model.New()

	glWrapper wrapper.Wrapper
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 1.0, 9.5}, mgl32.Vec3{0, 0, 1}, -90.0, -34.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

// It generates a bunch of triangles and sets their color to static blue.
func GenerateTriangles() {
	triang := triangle.New(90, 45, 45)
	v1, indicies1 := triang.ColoredMeshInput(triangleColorFront)
	v2, indicies2 := triang.ColoredMeshInput(triangleColorBack)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j * length)
			topY := float32(i * length)
			topZ := float32(0.0)

			m := mesh.NewColorMesh(v1, indicies1, triangleColorFront, glWrapper)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetScale(mgl32.Vec3{length, length, length})
			m.SetDirection(mgl32.Vec3{0, 0, 1})
			m.SetSpeed(float32(1.0) / float32(1000000000.0))

			Model.AddMesh(m)

			m2 := mesh.NewColorMesh(v2, indicies2, triangleColorBack, glWrapper)
			m2.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetDirection(mgl32.Vec3{0, 0, 1})
			m2.SetSpeed(float32(1.0) / float32(1000000000.0))

			Model.AddMesh(m2)
		}
	}
	Model.RotateX(45)
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
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/04-mesh-deformer-with-camera/vertexshader.vert", "examples/04-mesh-deformer-with-camera/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	GenerateTriangles()
	app.AddModelToShader(Model, shaderProgram)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
