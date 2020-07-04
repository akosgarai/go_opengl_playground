package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/cylinder"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - textured lighting map"

	CameraDirectionSpeed = float32(0.00500)
	CameraMoveSpeed      = 0.005
	CameraDistance       = 0.1

	LightSourceRoundSpeed = 3000.0
)

var (
	app  *application.Application
	cube *cuboid.Cuboid

	lastUpdate int64

	InitialCenterPointLight = mgl32.Vec3{-3, 0, -3}

	LightSource       *light.Light
	LightSourceSphere *mesh.MaterialMesh
	CubeMesh          *mesh.TexturedMesh

	TexModel = model.New()
	MatModel = model.New()

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
	camera := camera.NewCamera(mgl32.Vec3{3.3, -10, 14.0}, mgl32.Vec3{0, 1, 0}, -101.0, 21.5)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}

// It generates the lightsource sphere.
func CreateWhiteSphere() {
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 64.0)
	sph := sphere.New(15)
	v, i, _ := sph.MaterialMeshInput()
	LightSourceSphere = mesh.NewMaterialMesh(v, i, mat, glWrapper)
	LightSourceSphere.SetPosition(mgl32.Vec3{-3.0, -0.5, -3.0})
	LightSourceSphere.SetDirection((mgl32.Vec3{9, 0, -3}).Normalize())
	distance := (LightSourceSphere.GetPosition().Sub(CubeMesh.GetPosition())).Len()
	LightSourceSphere.SetSpeed((float32(2) * float32(3.1415) * distance) / LightSourceRoundSpeed)
	LightSourceSphere.SetScale(mgl32.Vec3{0.15, 0.15, 0.15})
}

// It generates a cube.
func CreateCubeMesh(t texture.Textures) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	return m
}
func CreateCylinderMesh(t texture.Textures) *mesh.TexturedMesh {
	c := cylinder.New(0.75, 30, 3)
	v, i, _ := c.TexturedMeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
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
	currentLightSourceDirection := LightSourceSphere.GetDirection()
	LightSourceSphere.SetDirection(mgl32.TransformNormal(currentLightSourceDirection, lightDirectionRotationMatrix))
	LightSource.SetPosition(LightSourceSphere.GetPosition())

	app.Update(delta)
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	scrn.SetCamera(CreateCamera())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(CameraDistance)

	LightSource = light.NewPointLight([4]mgl32.Vec3{InitialCenterPointLight, mgl32.Vec3{0.2, 0.2, 0.2}, mgl32.Vec3{0.5, 0.5, 0.5}, mgl32.Vec3{1, 1, 1}}, [3]float32{1.0, 1.0, 1.0})
	scrn.AddPointLightSource(LightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "", "", ""})

	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/colored-image-for-texture-testing-diffuse.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	tex.AddTexture(baseDir()+"/assets/colored-image-for-texture-testing-specular.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	CubeMesh = CreateCubeMesh(tex)
	TexModel.AddMesh(CubeMesh)
	CylinderMesh := CreateCylinderMesh(tex)
	CylinderMesh.SetPosition(mgl32.Vec3{2, 2, 2})
	TexModel.AddMesh(CylinderMesh)
	scrn.AddModelToShader(TexModel, shaderProgramTexture)

	shaderProgramWhite := shader.NewShader(baseDir()+"/shaders/lightsource.vert", baseDir()+"/shaders/lightsource.frag", glWrapper)
	scrn.AddShader(shaderProgramWhite)

	CreateWhiteSphere()
	MatModel.AddMesh(LightSourceSphere)
	scrn.AddModelToShader(MatModel, shaderProgramWhite)
	scrn.Setup(setupApp)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
