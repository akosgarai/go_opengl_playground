package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - plane with ball"

	moveSpeed = 1.0 / 100.0
	ballSpeed = float32(1.0 / 10000000.0 / 5.0)

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE
)

var (
	app           *application.Application
	ball          *sphere.Sphere
	BallPrecision = 10

	lastUpdate           int64
	cameraDistance       = 0.1
	cameraDirectionSpeed = float32(0.005)

	BallInitialDirection = mgl32.Vec3{0, -1, 0}
	BallTopPosition      = float32(-10)
	BallBottomPosition   = float32(-2)
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-10, -4, 22.0}, mgl32.Vec3{0, 1, 0}, 300.0, 16.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

// It generates a Sphere.
func GenerateSphere(shaderProgram *shader.Shader) {
	ball = sphere.New(mgl32.Vec3{0, -5, 0}, mgl32.Vec3{1, 0, 0}, 2.0, shaderProgram)
	ball.SetDirection(BallInitialDirection)
	ball.SetSpeed(ballSpeed)
	ball.SetPrecision(BallPrecision)
	app.AddItem(ball)
}

// It generates a square.
func GenerateSquare(shaderProgram *shader.Shader) {
	squareColor := mgl32.Vec3{0, 1, 0}
	coords := [4]mgl32.Vec3{
		mgl32.Vec3{-20, 0, -20},
		mgl32.Vec3{20, 0, -20},
		mgl32.Vec3{20, 0, 20},
		mgl32.Vec3{-20, 0, 20},
	}
	colors := [4]mgl32.Vec3{squareColor, squareColor, squareColor, squareColor}
	square := rectangle.New(coords, colors, shaderProgram)
	square.SetPrecision(1)
	app.AddItem(square)
}

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano - lastUpdate)
	moveTime := delta / float64(time.Millisecond)
	// handle ball
	if ball.GetCenter().Y() <= BallTopPosition {
		ball.SetCenter(mgl32.Vec3{ball.GetCenter().X(), BallTopPosition, ball.GetCenter().Z()})
		ball.SetDirection(BallInitialDirection.Mul(-1.0))
	}
	if ball.GetCenter().Y() >= BallBottomPosition {
		ball.SetCenter(mgl32.Vec3{ball.GetCenter().X(), BallBottomPosition, ball.GetCenter().Z()})
		ball.SetDirection(BallInitialDirection)
	}
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
	shader.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/05-ball-with-camera/vertexshader.vert", "examples/05-ball-with-camera/fragmentshader.frag")
	GenerateSquare(shaderProgram)
	GenerateSphere(shaderProgram)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		Update()
		app.DrawWithUniforms()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
