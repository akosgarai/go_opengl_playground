package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - material light"

	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Go left
	RIGHT    = glfw.KeyD // Go right
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
	app.SetWindow(window.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	app.SetCamera(CreateCamera())

	lightSource := light.New(mgl32.Vec3{-3, 0, -3}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1})
	shaderProgramColored := shader.NewShader("examples/08-basic-lightsource/vertexshader.vert", "examples/08-basic-lightsource/fragmentshader.frag")
	shaderProgramColored.SetLightSource(lightSource, "light.position", "light.ambient", "light.diffuse", "light.specular")
	GenerateColoredCube(shaderProgramColored)
	shaderProgramWhite := shader.NewShader("examples/08-basic-lightsource/vertexshader.vert", "examples/08-basic-lightsource/fragmentshader.frag")
	shaderProgramWhite.SetLightSource(lightSource, "light.position", "light.ambient", "light.diffuse", "light.specular")
	GenerateWhiteCube(shaderProgramWhite)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		shaderProgramColored.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
		shaderProgramWhite.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
