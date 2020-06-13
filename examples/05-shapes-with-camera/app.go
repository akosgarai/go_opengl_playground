package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/cylinder"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - shapes with camera"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.005)
	CameraDistance       = 0.1
)

var (
	app *application.Application

	lastUpdate int64
	colors     = []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{1.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 0.0, 1.0},
		mgl32.Vec3{1.0, 0.0, 1.0},
	}
	Model = model.New()

	glWrapper glwrapper.Wrapper
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
	camera := camera.NewCamera(mgl32.Vec3{-3, -5, 18.0}, mgl32.Vec3{0, 1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	app.Update(delta)
	lastUpdate = nowNano
}

// It generates a cube.
func CreateCubeMesh() *mesh.ColorMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.ColoredMeshInput(colors)
	m := mesh.NewColorMesh(v, i, colors, glWrapper)
	m.SetPosition(mgl32.Vec3{-0.5, -0.5, 0.5})
	return m
}

// It generates a Sphere.
func CreateSphereMesh() *mesh.ColorMesh {
	s := sphere.New(20)
	cols := []mgl32.Vec3{colors[4]}
	v, i, _ := s.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(mgl32.Vec3{3, 3, 5})
	m.SetScale(mgl32.Vec3{2, 2, 2})
	return m
}

// It generates a cylinder
func CreateCylinder() *mesh.ColorMesh {
	c := cylinder.New(1.5, 30, 3)
	cols := []mgl32.Vec3{colors[3]}
	v, i, _ := c.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(mgl32.Vec3{-3, -3, 5})
	return m
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
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

	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	sphereMesh := CreateSphereMesh()
	Model.AddMesh(sphereMesh)
	cubeMesh := CreateCubeMesh()
	Model.AddMesh(cubeMesh)
	cyl := CreateCylinder()
	Model.AddMesh(cyl)
	app.AddModelToShader(Model, shaderProgram)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
