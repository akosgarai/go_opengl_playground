package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
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
	WindowTitle  = "Example - cubes with light source"

	CameraMoveSpeed      = 0.005
	cameraDistance       = 0.1
	cameraDirectionSpeed = float32(0.00500)
)

var (
	app *application.Application

	lastUpdate int64

	Model = model.New()

	glWrapper wrapper.Wrapper
)

// Setup keymap for the camera movement
func CameraMovementMap() map[string]glfw.Key {
	cm := make(map[string]glfw.Key)
	cm["forward"] = glfw.KeyW
	cm["back"] = glfw.KeyS
	cm["up"] = glfw.KeyQ
	cm["down"] = glfw.KeyE
	cm["left"] = glfw.KeyA
	cm["right"] = glfw.KeyD
	return cm
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0, 0, 10.0}, mgl32.Vec3{0, 1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	return camera
}

func CreateColoredCubeMesh(pos mgl32.Vec3, col []mgl32.Vec3) *mesh.ColorMesh {
	cube := cuboid.NewCube()
	v, i := cube.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(pos)
	return m
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)

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
	app.SetCameraMovementMap(CameraMovementMap())

	shaderProgram := shader.NewShader("examples/08-colors/shaders/vertexshader.vert", "examples/08-colors/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	lightSource := light.NewPointLight([4]mgl32.Vec3{mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	app.AddPointLightSource(lightSource, [7]string{"", "light.ambient", "", "", "", "", ""})

	whiteCube := CreateColoredCubeMesh(mgl32.Vec3{-3.0, -0.5, -3.0}, []mgl32.Vec3{mgl32.Vec3{1.0, 1.0, 1.0}})
	Model.AddMesh(whiteCube)
	colors := []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{1.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 0.0, 1.0},
		mgl32.Vec3{1.0, 0.0, 1.0},
	}
	coloredCube := CreateColoredCubeMesh(mgl32.Vec3{0.0, 0.0, 0.0}, colors)
	Model.AddMesh(coloredCube)
	app.AddModelToShader(Model, shaderProgram)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
