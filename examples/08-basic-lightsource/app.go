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
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - cubes, basic lighting"

	CameraMoveSpeed      = 0.005
	CameraDistance       = 0.1
	CameraDirectionSpeed = float32(0.00500)
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
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}

// It generates the material cube mesh.
func CreateMaterialCubeMesh(position mgl32.Vec3, mat *material.Material) *mesh.MaterialMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(position)
	return m
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
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

	lightSource := light.NewPointLight([4]mgl32.Vec3{mgl32.Vec3{-3, 0, -3}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	app.AddPointLightSource(lightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})
	shaderProgram := shader.NewShader("examples/08-basic-lightsource/shaders/vertexshader.vert", "examples/08-basic-lightsource/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)
	whiteMat := material.New(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, 36.0)
	whiteCube := CreateMaterialCubeMesh(mgl32.Vec3{-3.0, 0.0, -3.0}, whiteMat)
	Model.AddMesh(whiteCube)
	colorMat := material.New(mgl32.Vec3{0.0, 0.3, 0.3}, mgl32.Vec3{0, 1, 1}, mgl32.Vec3{0, 1, 1}, 36.0)
	coloredCube := CreateMaterialCubeMesh(mgl32.Vec3{0.0, 0.0, 0.0}, colorMat)
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
