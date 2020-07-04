package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
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
	WindowTitle  = "Example - multiple light source"

	MoveSpeed            = 0.005
	CameraDirectionSpeed = float32(0.100)
	CameraMoveSpeed      = 0.005
	CameraDistance       = 0.1
)

var (
	app                       *application.Application
	Bug1                      *model.Bug
	Bug2                      *mesh.TexturedMesh
	BugOneLastRotate          int64
	lastUpdate                int64
	ShaderProgramsWithViewPos []*shader.Shader
	DirectionalLightSource    *light.Light
	PointLightSource_1        *light.Light
	PointLightSource_2        *light.Light
	SpotLightSource_1         *light.Light
	SpotLightSource_2         *light.Light

	BugOneForwardMove         = float64(1000)
	DirectionalLightDirection = (mgl32.Vec3{0.7, 0.7, 0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightDiffuse   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightSpecular  = mgl32.Vec3{0.1, 0.1, 0.1}
	PointLightAmbient         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular        = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition_1      = mgl32.Vec3{8, -0.5, -1.0}
	PointLightPosition_2      = mgl32.Vec3{8, -5, -30}
	LightConstantTerm         = float32(1.0)
	LightLinearTerm           = float32(0.14)
	LightQuadraticTerm        = float32(0.07)
	SpotLightAmbient          = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse          = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular         = mgl32.Vec3{1, 1, 1}
	SpotLightDirection_1      = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightDirection_2      = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightPosition_1       = mgl32.Vec3{0.20, -6, -0.65}
	SpotLightPosition_2       = mgl32.Vec3{10.20, -6, -0.65}
	SpotLightCutoff_1         = float32(4)
	SpotLightCutoff_2         = float32(4)
	SpotLightOuterCutoff_1    = float32(5)
	SpotLightOuterCutoff_2    = float32(5)
	TexModel                  = model.New()
	MatModel                  = model.New()

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

func CreateGrassMesh(t texture.Textures) *mesh.TexturedMesh {
	square := rectangle.NewSquare()
	v, i, bo := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetScale(mgl32.Vec3{1000, 1, 1000})
	m.SetBoundingObject(bo)
	return m
}
func CreateCubeMesh(t texture.Textures, pos mgl32.Vec3) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i, bo := cube.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(pos)
	m.SetBoundingObject(bo)
	return m
}

// It generates the lamp. Now it uses the StreetLamp model for creating it.
func StreetLamp(position mgl32.Vec3) *model.StreetLamp {
	StreetLamp := model.NewMaterialStreetLamp(position, 6, glWrapper)
	StreetLamp.RotateX(90)
	StreetLamp.RotateY(-90)
	return StreetLamp
}
func TexturedStreetLamp(position mgl32.Vec3) *model.StreetLamp {
	StreetLamp := model.NewTexturedStreetLamp(position, 6, glWrapper)
	StreetLamp.RotateX(90)
	StreetLamp.RotateY(-90)
	return StreetLamp
}

func TexturedBug(t texture.Textures) {
	sph := sphere.New(15)
	v, i, bo := sph.TexturedMeshInput()
	Bug2 = mesh.NewTexturedMesh(v, i, t, glWrapper)
	Bug2.SetPosition(PointLightPosition_2)
	Bug2.SetDirection(mgl32.Vec3{0, 0, 1})
	Bug2.SetSpeed(MoveSpeed)
	Bug2.SetBoundingObject(bo)
	TexModel.AddMesh(Bug2)
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-11.2, -5.0, 4.2}, mgl32.Vec3{0, 1, 0}, -37.0, -2.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.01, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func RotateBugOne(now int64) {
	moveTime := float64(now-BugOneLastRotate) / float64(time.Millisecond)
	if moveTime > BugOneForwardMove {
		BugOneLastRotate = now
		// rotate 45 deg
		Bug1.RotateY(-45)
	}
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	RotateBugOne(nowNano)
	PointLightSource_1.SetPosition(Bug1.GetBottomPosition())
	PointLightSource_2.SetPosition(Bug2.GetPosition())
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

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource = light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	PointLightSource_1 = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition_1,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	SpotLightSource_1 = light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition_1,
		SpotLightDirection_1,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_1, SpotLightOuterCutoff_1})
	SpotLightSource_2 = light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition_2,
		SpotLightDirection_2,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_2, SpotLightOuterCutoff_2})
	PointLightSource_2 = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition_2,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})

	// Add the lightources to the application
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	scrn.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	scrn.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	scrn.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	// grass textures
	var grassTexture texture.Textures
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	grassMesh := CreateGrassMesh(grassTexture)
	TexModel.AddMesh(grassMesh)

	// box textures
	var boxTexture texture.Textures
	boxTexture.AddTexture(baseDir()+"/assets/box-diffuse.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	boxTexture.AddTexture(baseDir()+"/assets/box-specular.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	// we have 3 boxes in the following coordinates.
	boxPositions := []mgl32.Vec3{
		mgl32.Vec3{-5.0, -0.51, 0.0},
		mgl32.Vec3{0.0, -0.51, 0.0},
		mgl32.Vec3{5.0, -0.51, 0.0},
	}
	for _, pos := range boxPositions {
		box := CreateCubeMesh(boxTexture, pos)
		TexModel.AddMesh(box)
	}

	// Shader application for the lamp
	shaderProgramMaterial := shader.NewShader(baseDir()+"/shaders/lamp.vert", baseDir()+"/shaders/lamp.frag", glWrapper)
	scrn.AddShader(shaderProgramMaterial)
	shaderProgramTextureMat := shader.NewShader(baseDir()+"/shaders/texturemat.vert", baseDir()+"/shaders/texturemat.frag", glWrapper)
	scrn.AddShader(shaderProgramTextureMat)

	lamp1 := TexturedStreetLamp(mgl32.Vec3{0.4, -6.0, -1.3})
	scrn.AddModelToShader(lamp1, shaderProgramTextureMat)
	lamp2 := StreetLamp(mgl32.Vec3{10.4, -6.0, -1.3})
	scrn.AddModelToShader(lamp2, shaderProgramMaterial)

	Bug1 = model.NewBug(mgl32.Vec3{9, -0.5, -1.0}, mgl32.Vec3{0.2, 0.2, 0.2}, glWrapper)
	Bug1.SetDirection(mgl32.Vec3{1, 0, 0})
	Bug1.SetSpeed(MoveSpeed)

	scrn.AddModelToShader(Bug1, shaderProgramMaterial)

	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	TexturedBug(sunTexture)
	scrn.AddModelToShader(TexModel, shaderProgramTexture)
	scrn.AddModelToShader(MatModel, shaderProgramMaterial)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.0, 0.0, 1.0)

	lastUpdate = time.Now().UnixNano()
	BugOneLastRotate = lastUpdate
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
