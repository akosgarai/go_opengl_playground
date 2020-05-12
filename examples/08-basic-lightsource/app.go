package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - cubes, basic lighting"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed = 0.005
)

var (
	app *application.Application

	lastUpdate int64

	cameraDistance       = 0.1
	cameraDirectionSpeed = float32(0.00500)
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0, 0, 10.0}, mgl32.Vec3{0, 1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
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

// It generates the colored cube.
func GenerateWhiteCube(shaderProgram *shader.Shader) {
	whiteBottomCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{-3.5, -0.5, -3.5},
		mgl32.Vec3{-3.5, -0.5, -2.5},
		mgl32.Vec3{-2.5, -0.5, -2.5},
		mgl32.Vec3{-2.5, -0.5, -3.5},
	}
	whiteBottomColor := [4]mgl32.Vec3{
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
	}
	bottomRect := rectangle.New(whiteBottomCoordinates, whiteBottomColor, shaderProgram)
	cube := cuboid.New(bottomRect, 1.0, shaderProgram)
	mat := material.New(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, 36.0)
	cube.SetMaterial(mat)
	cube.DrawMode(cuboid.DRAW_MODE_LIGHT)
	app.AddItem(cube)
}

// It generates the colored cube.
func GenerateColoredCube(shaderProgram *shader.Shader) {
	coloredBottomCoordinates := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, -0.5, -0.5},
		mgl32.Vec3{-0.5, -0.5, 0.5},
		mgl32.Vec3{0.5, -0.5, 0.5},
		mgl32.Vec3{0.5, -0.5, -0.5},
	}
	coloredBottomColor := [4]mgl32.Vec3{
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
	}
	bottomRect := rectangle.New(coloredBottomCoordinates, coloredBottomColor, shaderProgram)
	cube := cuboid.New(bottomRect, 1.0, shaderProgram)
	mat := material.New(mgl32.Vec3{0.0, 0.3, 0.3}, mgl32.Vec3{0, 1, 1}, mgl32.Vec3{0, 1, 1}, 36.0)
	cube.SetMaterial(mat)
	cube.DrawMode(cuboid.DRAW_MODE_LIGHT)
	app.AddItem(cube)
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano - lastUpdate)
	moveTime := delta / float64(time.Millisecond)
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
	wrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	lightSource := light.NewPointLight([4]mgl32.Vec3{mgl32.Vec3{-3, 0, -3}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	shaderProgramColored := shader.NewShader("examples/08-basic-lightsource/vertexshader.vert", "examples/08-basic-lightsource/fragmentshader.frag")
	shaderProgramColored.AddPointLightSource(lightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})
	GenerateColoredCube(shaderProgramColored)
	shaderProgramWhite := shader.NewShader("examples/08-basic-lightsource/vertexshader.vert", "examples/08-basic-lightsource/fragmentshader.frag")
	shaderProgramWhite.AddPointLightSource(lightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})
	GenerateWhiteCube(shaderProgramWhite)

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)
	wrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		shaderProgramColored.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
		shaderProgramWhite.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
