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
	"github.com/akosgarai/opengl_playground/pkg/primitives/cylinder"
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
	WindowTitle  = "Example - material light - with rotation"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.00500)
	CameraDistance       = 0.1

	LightSourceRoundSpeed = 3000.0
)

var (
	app *application.Application

	lastUpdate int64

	LightSource     *light.Light
	LightSourceCube *mesh.MaterialMesh
	JadeCube        *mesh.MaterialMesh

	InitialCenterPointLight = mgl32.Vec3{-3, 0, -3}
	CenterPointObject       = mgl32.Vec3{0, 0, 0}
	Model                   = model.New()

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
	camera := camera.NewCamera(mgl32.Vec3{3.3, -10, 14.0}, mgl32.Vec3{0, -1, 0}, -101.0, 21.5)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}

func CreateMaterialCube(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	cube := cuboid.NewCube()
	v, i := cube.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	return m
}
func CreateMaterialCylinder(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	c := cylinder.New(0.75, 30, 3)
	v, i := c.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	return m
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	lightSourceRotationAngleRadian := mgl32.DegToRad(float32((360 / LightSourceRoundSpeed) * delta))
	lightDirectionRotationMatrix := mgl32.HomogRotate3D(lightSourceRotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentLightSourceDirection := LightSourceCube.GetDirection()
	LightSourceCube.SetDirection(mgl32.TransformNormal(currentLightSourceDirection, lightDirectionRotationMatrix))
	LightSource.SetPosition(LightSourceCube.GetPosition())

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

	LightSource = light.NewPointLight([4]mgl32.Vec3{InitialCenterPointLight, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	app.AddPointLightSource(LightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})

	materialShader := shader.NewShader("examples/08-material-light/shaders/vertexshader.vert", "examples/08-material-light/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(materialShader)

	JadeCube = CreateMaterialCube(material.Jade, mgl32.Vec3{0.0, 0.0, 0.0})
	Model.AddMesh(JadeCube)

	rpCube := CreateMaterialCube(material.Redplastic, mgl32.Vec3{-6.5, -3.5, -4.5})
	rpCube.SetScale(mgl32.Vec3{2.0, 2.0, 2.0})
	Model.AddMesh(rpCube)

	obsidianCube := CreateMaterialCube(material.Obsidian, mgl32.Vec3{-7.5, -4.5, -0.5})
	Model.AddMesh(obsidianCube)

	copperCube := CreateMaterialCube(material.Copper, mgl32.Vec3{2.0, -4.5, -0.5})
	Model.AddMesh(copperCube)

	silverCube := CreateMaterialCube(material.Silver, mgl32.Vec3{2.0, -2.5, -1.5})
	Model.AddMesh(silverCube)

	turqCylinder := CreateMaterialCylinder(material.Turquoise, mgl32.Vec3{4, 3, -3})
	Model.AddMesh(turqCylinder)

	mat := material.New(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, 144.0)
	LightSourceCube = CreateMaterialCube(mat, mgl32.Vec3{-3.0, -1.5, -3.0})
	LightSourceCube.SetDirection((mgl32.Vec3{9, 0, -3}).Normalize())

	distance := (LightSourceCube.GetPosition().Sub(JadeCube.GetPosition())).Len()
	LightSourceCube.SetSpeed((float32(2) * float32(3.1415) * distance) / LightSourceRoundSpeed)
	Model.AddMesh(LightSourceCube)
	app.AddModelToShader(Model, materialShader)

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
