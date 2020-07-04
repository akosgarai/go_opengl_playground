package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
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
	WindowTitle  = "Example - textured spheres"

	CameraDirectionSpeed = float32(0.010)
	CameraMoveSpeed      = 0.005
	CameraDistance       = 0.1

	EarthRoundSpeed       = 3000.0
	OtherPlanetRoundSpeed = 7000.0
	SunRoundSpeed         = 0.01
)

var (
	app              *application.Application
	Sun              *mesh.TexturedMesh
	Earth            *mesh.TexturedMesh
	MatPlanet        *mesh.TexturedMaterialMesh
	lastUpdate       int64
	PointLightSource *light.Light

	rotationAngle      = float32(0.0)
	PointLightAmbient  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition = mgl32.Vec3{0.0, 0.0, 0.0}
	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.14)
	LightQuadraticTerm = float32(0.07)
	spherePrimitive    = sphere.New(20)
	TexModel           = model.New()
	TexMatModel        = model.New()
	CmModel            = model.New()

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

func TexturedSphere(t texture.Textures, position mgl32.Vec3, scale float32, shaderProgram *shader.Shader) *mesh.TexturedMesh {
	v, i, _ := spherePrimitive.TexturedMeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(position)
	m.SetScale(mgl32.Vec3{scale, scale, scale})
	return m
}
func TexturedMaterialSphere(t texture.Textures, mat *material.Material, position mgl32.Vec3, scale float32, shaderProgram *shader.Shader) *mesh.TexturedMaterialMesh {
	v, i, _ := spherePrimitive.TexturedMeshInput()
	m := mesh.NewTexturedMaterialMesh(v, i, t, mat, glWrapper)
	m.SetPosition(position)
	m.SetScale(mgl32.Vec3{scale, scale, scale})
	return m
}

// It generates a cube map.
func CubeMap(t texture.Textures) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	return m
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -10.0}, mgl32.Vec3{0, 1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.01, 200.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func updateSun(moveTime float64) {
	Sun.RotateY(float32(moveTime) * SunRoundSpeed)
}
func updatePlanets(moveTime float64) {
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	rotationAngleRadian := mgl32.DegToRad(float32((360 / EarthRoundSpeed) * moveTime))
	rotationMatrix := mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentDirection := Earth.GetDirection()
	Earth.SetDirection(mgl32.TransformNormal(currentDirection, rotationMatrix))

	rotationAngleRadian = mgl32.DegToRad(float32((360 / OtherPlanetRoundSpeed) * moveTime))
	rotationMatrix = mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentDirection = MatPlanet.GetDirection()
	MatPlanet.SetDirection(mgl32.TransformNormal(currentDirection, rotationMatrix))
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	updatePlanets(delta)
	updateSun(delta)

	app.Update(delta)
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

	scrn := screen.New()
	scrn.SetCamera(CreateCamera())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(CameraDistance)

	PointLightSource = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})

	// Add the lightources to the application
	scrn.AddPointLightSource(PointLightSource, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	Sun = TexturedSphere(sunTexture, mgl32.Vec3{0.0, 0.0, 0.0}, 1, shaderProgramTexture)
	TexModel.AddMesh(Sun)
	// sun texture
	var earthTexture texture.Textures
	earthTexture.AddTexture(baseDir()+"/assets/earth.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	earthTexture.AddTexture(baseDir()+"/assets/earth.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	Earth = TexturedSphere(earthTexture, mgl32.Vec3{3.0, 0.0, 0.0}, 0.1, shaderProgramTexture)
	distance := Earth.GetPosition().Len()
	Earth.SetSpeed((float32(2) * float32(3.1415) * distance) / EarthRoundSpeed)
	Earth.SetDirection((mgl32.Vec3{0, 0, 1}).Normalize())
	TexModel.AddMesh(Earth)
	scrn.AddModelToShader(TexModel, shaderProgramTexture)
	// other planet texture
	shaderProgramTextureMaterial := shader.NewShader(baseDir()+"/shaders/texturemat.vert", baseDir()+"/shaders/texturemat.frag", glWrapper)
	scrn.AddShader(shaderProgramTextureMaterial)
	var materialTexture texture.Textures
	materialTexture.AddTexture(baseDir()+"/assets/venus.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	materialTexture.AddTexture(baseDir()+"/assets/venus.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", glWrapper)
	MatPlanet = TexturedMaterialSphere(materialTexture, material.Gold, mgl32.Vec3{2.0, 0.0, 0.0}, 0.15, shaderProgramTextureMaterial)
	distance = MatPlanet.GetPosition().Len()
	MatPlanet.SetSpeed((float32(2) * float32(3.1415) * distance) / OtherPlanetRoundSpeed)
	MatPlanet.SetDirection((mgl32.Vec3{0, 0, 1}).Normalize())
	TexMatModel.AddMesh(MatPlanet)
	scrn.AddModelToShader(TexMatModel, shaderProgramTextureMaterial)

	shaderProgramCubeMap := shader.NewShader(baseDir()+"/shaders/cubeMap.vert", baseDir()+"/shaders/cubeMap.frag", glWrapper)
	scrn.AddShader(shaderProgramCubeMap)
	var cubeMapTexture texture.Textures
	cubeMapTexture.AddCubeMapTexture(baseDir()+"/assets", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "skybox", glWrapper)
	cubeMap := CubeMap(cubeMapTexture)
	cubeMap.SetScale(mgl32.Vec3{100.0, 100.0, 100.0})
	CmModel.AddMesh(cubeMap)
	scrn.AddModelToShader(CmModel, shaderProgramCubeMap)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.0, 0.0, 1.0)

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
