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
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - plane with ball"

	moveSpeed          = float32(1.0 / 100.0)
	BallSpeed          = float32(1.0 / 10.0 / 5.0)
	BallPrecision      = 10
	BallTopPosition    = float32(10)
	BallBottomPosition = float32(2)
	CameraDistance     = float32(0.1)
)

var (
	app    *application.Application
	Ball   *mesh.ColorMesh
	Ground *mesh.ColorMesh

	lastUpdate int64

	BallInitialDirection = mgl32.Vec3{0, 1, 0}
	Model                = model.New()

	glWrapper wrapper.Wrapper
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 5.0, -24.0}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(moveSpeed)
	camera.SetRotationStep(0.005)
	return camera
}

func CreateSphereMesh() *mesh.ColorMesh {
	s := sphere.New(BallPrecision)
	cols := []mgl32.Vec3{mgl32.Vec3{1, 0, 0}}
	v, i, _ := s.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(mgl32.Vec3{0, 5, 0})
	m.SetScale(mgl32.Vec3{2, 2, 2})
	m.SetDirection(BallInitialDirection)
	m.SetSpeed(BallSpeed)
	return m
}

// It generates a square.
func CreateSquareMesh() *mesh.ColorMesh {
	squareColor := []mgl32.Vec3{mgl32.Vec3{0, 1, 0}}
	s := rectangle.NewSquare()
	v, i, _ := s.ColoredMeshInput(squareColor)
	m := mesh.NewColorMesh(v, i, squareColor, glWrapper)
	m.SetScale(mgl32.Vec3{40, 40, 40})
	bo := make(map[string]float32)
	bo["width"] = float32(40.0)
	bo["height"] = float32(40.0)
	bo["length"] = float32(0.05)
	m.SetBoundingObjectParams(bo)
	return m
}

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

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	// handle ball
	if Ball.GetPosition().Y() >= BallTopPosition {
		Ball.SetPosition(mgl32.Vec3{Ball.GetPosition().X(), BallTopPosition, Ball.GetPosition().Z()})
		Ball.SetDirection(BallInitialDirection.Mul(-1.0))
	}
	if Ball.GetPosition().Y() <= BallBottomPosition {
		Ball.SetPosition(mgl32.Vec3{Ball.GetPosition().X(), BallBottomPosition, Ball.GetPosition().Z()})
		Ball.SetDirection(BallInitialDirection)
	}
	lastUpdate = nowNano
	app.Update(delta)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())
	app.SetCameraMovementMap(CameraMovementMap())
	app.SetRotateOnEdgeDistance(CameraDistance)

	shaderProgram := shader.NewShader("examples/05-ball-with-camera/shaders/vertexshader.vert", "examples/05-ball-with-camera/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	Ball = CreateSphereMesh()
	Model.AddMesh(Ball)
	Ground = CreateSquareMesh()
	Model.AddMesh(Ground)
	app.AddModelToShader(Model, shaderProgram)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
