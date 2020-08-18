package main

import (
	"os"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
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

	FORM_ENV_NAME = "SETTINGS"
	OFF_VALUE     = "off"
)

var (
	app                    *application.Application
	Bug1                   *model.Bug
	Bug2                   *mesh.TexturedMesh
	BugOneLastRotate       int64
	lastUpdate             int64
	startTime              int64
	DirectionalLightSource *light.Light
	PointLightSource_1     *light.Light
	PointLightSource_2     *light.Light

	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	glWrapper glwrapper.Wrapper
)

func InitSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("SurfaceSize", "Surface size", "The size of the surface area. This value is used for scaling", float32(100.0), nil)
	// light sources
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	// directional light
	Settings.AddConfig("DLAmbient", "DL ambient", "The ambient color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDiffuse", "DL diffuse", "The diffuse color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLSpecular", "DL specular", "The specular color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDirection", "DL direction", "The direction vector of the directional lightsource.", mgl32.Vec3{0.7, 0.7, 0.7}, nil)
	// point light
	Settings.AddConfig("PLAmbient", "PL ambient", "The ambient color component of the point lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("PLDiffuse", "PL diffuse", "The diffuse color component of the point lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("PLSpecular", "PL specular", "The specular color component of the point lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("PL1Position", "PL1 position", "The position vector of the 1. point lightsource.", mgl32.Vec3{8, -0.5, -1.0}, nil)
	Settings.AddConfig("PL2Position", "PL2 position", "The position vector of the 2. point lightsource.", mgl32.Vec3{8, -5, -30}, nil)
	// spot light
	Settings.AddConfig("SLAmbient", "SL ambient", "The ambient color component of the spot lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("SLDiffuse", "SL diffuse", "The diffuse color component of the spot lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("SLSpecular", "SL specular", "The specular color component of the spot lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("SLCutoff", "SL cutoff", "The inner cutoff component of the spot light equasion.", float32(4.0), nil)
	Settings.AddConfig("SLOuterCutoff", "SL o. cutoff", "The inner cutoff component of the spot light equasion.", float32(5.0), nil)
	// boxes
	Settings.AddConfig("Box1Position", "Box 1 pos", "The position vector of the 1. box.", mgl32.Vec3{-5.0, -0.51, 0.0}, nil)
	Settings.AddConfig("Box2Position", "Box 2 pos", "The position vector of the 2. box.", mgl32.Vec3{0.0, -0.51, 0.0}, nil)
	Settings.AddConfig("Box3Position", "Box 3 pos", "The position vector of the 3. box.", mgl32.Vec3{5.0, -0.51, 0.0}, nil)
	// streetlamps
	Settings.AddConfig("StreetLamp1Position", "St. lamp 1 pos", "The position vector of the 1. street lamp (textured).", mgl32.Vec3{0.4, -6.0, -1.3}, nil)
	Settings.AddConfig("StreetLamp2Position", "St. lamp 2 pos", "The position vector of the 2. street lamp (material).", mgl32.Vec3{10.4, -6.0, -1.3}, nil)
	// bug that moves around.
	Settings.AddConfig("BugPosition", "Bug position", "The position vector of the material bug.", mgl32.Vec3{9, -0.5, -1.0}, nil)
	Settings.AddConfig("BugScale", "Bug scale", "The scale vector of the material bug.", mgl32.Vec3{0.2, 0.2, 0.2}, nil)
	Settings.AddConfig("BugVelocity", "Bug velocity", "The speed of the material bug.", float32(0.005), nil)
	Settings.AddConfig("BugRotationAngle", "Bug rotation", "The rotation angle of the material bug.", float32(-45.0), nil)
	Settings.AddConfig("BugForwardTime", "Bug forward ms", "The time in ms during the material bug goes into the same direction.", float32(1000.0), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{-11.2, -5.0, 4.2}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-37.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(-2.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.001), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(50.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.005), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
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

func CreateGrassMesh(t texture.Textures) *mesh.TexturedMesh {
	square := rectangle.NewSquare()
	v, i, bo := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	scale := Settings["SurfaceSize"].GetCurrentValue().(float32)
	m.SetScale(mgl32.Vec3{scale, 1, scale})
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
func StreetLamp() *model.StreetLamp {
	builder := model.NewStreetLampBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetPosition(Settings["StreetLamp1Position"].GetCurrentValue().(mgl32.Vec3))
	c := Settings["LSConstantTerm"].GetCurrentValue().(float32)
	l := Settings["LSLinearTerm"].GetCurrentValue().(float32)
	q := Settings["LSQuadraticTerm"].GetCurrentValue().(float32)
	builder.SetLightTerms(c, l, q)
	cutoff := Settings["SLCutoff"].GetCurrentValue().(float32)
	outerCutoff := Settings["SLOuterCutoff"].GetCurrentValue().(float32)
	builder.SetCutoff(cutoff, outerCutoff)
	builder.SetPoleLength(6.0)
	builder.SetRotation(90, -90, 0)
	mat := material.New(Settings["SLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["SLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["SLSpecular"].GetCurrentValue().(mgl32.Vec3),
		256.0)
	builder.SetBulbMaterial(mat)
	return builder.BuildMaterial()
}
func TexturedStreetLamp() *model.StreetLamp {
	builder := model.NewStreetLampBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetPosition(Settings["StreetLamp2Position"].GetCurrentValue().(mgl32.Vec3))
	c := Settings["LSConstantTerm"].GetCurrentValue().(float32)
	l := Settings["LSLinearTerm"].GetCurrentValue().(float32)
	q := Settings["LSQuadraticTerm"].GetCurrentValue().(float32)
	builder.SetLightTerms(c, l, q)
	cutoff := Settings["SLCutoff"].GetCurrentValue().(float32)
	outerCutoff := Settings["SLOuterCutoff"].GetCurrentValue().(float32)
	builder.SetCutoff(cutoff, outerCutoff)
	builder.SetPoleLength(6.0)
	builder.SetRotation(90, -90, 0)
	mat := material.New(Settings["SLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["SLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["SLSpecular"].GetCurrentValue().(mgl32.Vec3),
		256.0)
	builder.SetBulbMaterial(mat)
	return builder.BuildTexture()
}

func TexturedBug(t texture.Textures) *mesh.TexturedMesh {
	sph := sphere.New(15)
	v, i, bo := sph.TexturedMeshInput()
	Bug2 = mesh.NewTexturedMesh(v, i, t, glWrapper)
	Bug2.SetPosition(Settings["PL2Position"].GetCurrentValue().(mgl32.Vec3))
	Bug2.SetDirection(mgl32.Vec3{0, 0, 1})
	Bug2.SetSpeed(Settings["BugVelocity"].GetCurrentValue().(float32))
	Bug2.SetBoundingObject(bo)
	return Bug2
}

// It creates a new camera with the necessary setup from settings screen
func CreateCameraFromSettings() *camera.Camera {
	cameraPosition := Settings["CameraPos"].GetCurrentValue().(mgl32.Vec3)
	worldUp := Settings["WorldUp"].GetCurrentValue().(mgl32.Vec3)
	yawAngle := Settings["CameraYaw"].GetCurrentValue().(float32)
	pitchAngle := Settings["CameraPitch"].GetCurrentValue().(float32)
	fov := Settings["CameraFov"].GetCurrentValue().(float32)
	near := Settings["CameraNear"].GetCurrentValue().(float32)
	far := Settings["CameraFar"].GetCurrentValue().(float32)
	moveSpeed := Settings["CameraVelocity"].GetCurrentValue().(float32)
	directionSpeed := Settings["CameraRotation"].GetCurrentValue().(float32)
	camera := camera.NewCamera(cameraPosition, worldUp, yawAngle, pitchAngle)
	camera.SetupProjection(fov, float32(WindowWidth)/float32(WindowHeight), near, far)
	camera.SetVelocity(moveSpeed)
	camera.SetRotationStep(directionSpeed)
	return camera
}
func RotateBugOne(now int64) {
	moveTime := float64(now-BugOneLastRotate) / float64(time.Millisecond)
	if moveTime > float64(Settings["BugForwardTime"].GetCurrentValue().(float32)) {
		BugOneLastRotate = now
		// rotate 45 deg
		Bug1.RotateY(Settings["BugRotationAngle"].GetCurrentValue().(float32))
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
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	clearColor := Settings["ClearCol"].GetCurrentValue().(mgl32.Vec3)
	glWrapper.ClearColor(clearColor.X(), clearColor.Y(), clearColor.Z(), 1.0)
}
func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"ClearCol",
		"LSConstantTerm", "LSLinearTerm",
		"LSQuadraticTerm",

		"DLAmbient",
		"DLDiffuse",
		"DLSpecular",
		"DLDirection",

		"PLAmbient",
		"PLDiffuse",
		"PLSpecular",
		"PL1Position",
		"PL2Position",

		"SLAmbient",
		"SLDiffuse",
		"SLSpecular",
		"SLCutoff", "SLOuterCutoff",

		"Box1Position",
		"Box2Position",
		"Box3Position",

		"StreetLamp1Position",
		"StreetLamp2Position",

		"BugPosition",
		"BugScale",
		"BugVelocity", "BugRotationAngle",
		"BugForwardTime",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Settings")
}
func createMenu() *screen.MenuScreen {
	contS := func(m map[string]bool) bool {
		return m["world-started"]
	}
	contNS := func(m map[string]bool) bool {
		return !m["world-started"]
	}
	contAll := func(m map[string]bool) bool {
		return true
	}
	restartEvent := func() {
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = mainScreen()
		app.ActivateScreen(AppScreen)
	}
	startEvent := func() {
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = mainScreen()
		app.ActivateScreen(AppScreen)
		MenuScreen.SetState("world-started", true)
		MenuScreen.BuildScreen()
	}
	settingsEvent := func() {
		app.ActivateScreen(SettingsScreen)
	}
	continueEvent := func() {
		app.ActivateScreen(AppScreen)
	}
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	options := []screen.Option{
		*screen.NewMenuScreenOption("continue", contS, continueEvent),
		*screen.NewMenuScreenOption("start", contNS, startEvent),
		*screen.NewMenuScreenOption("restart", contS, restartEvent),
		*screen.NewMenuScreenOption("settings", contAll, settingsEvent),
		*screen.NewMenuScreenOption("exit", contAll, exitEvent),
	}
	return app.BuildMenuScreen(options)
}

func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetCamera(CreateCameraFromSettings())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(Settings["CameraRotationEdge"].GetCurrentValue().(float32))

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource = light.NewDirectionalLight([4]mgl32.Vec3{
		Settings["DLDirection"].GetCurrentValue().(mgl32.Vec3).Normalize(),
		Settings["DLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["DLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["DLSpecular"].GetCurrentValue().(mgl32.Vec3),
	})
	PointLightSource_1 = light.NewPointLight([4]mgl32.Vec3{
		Settings["PL1Position"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLSpecular"].GetCurrentValue().(mgl32.Vec3)},
		[3]float32{
			Settings["LSConstantTerm"].GetCurrentValue().(float32),
			Settings["LSLinearTerm"].GetCurrentValue().(float32),
			Settings["LSQuadraticTerm"].GetCurrentValue().(float32),
		})
	PointLightSource_2 = light.NewPointLight([4]mgl32.Vec3{
		Settings["PL2Position"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["PLSpecular"].GetCurrentValue().(mgl32.Vec3)},
		[3]float32{
			Settings["LSConstantTerm"].GetCurrentValue().(float32),
			Settings["LSLinearTerm"].GetCurrentValue().(float32),
			Settings["LSQuadraticTerm"].GetCurrentValue().(float32),
		})

	// Add the lightources to the application
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	scrn.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	TexModel := model.New()
	MatModel := model.New()

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
		Settings["Box1Position"].GetCurrentValue().(mgl32.Vec3),
		Settings["Box2Position"].GetCurrentValue().(mgl32.Vec3),
		Settings["Box3Position"].GetCurrentValue().(mgl32.Vec3),
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

	lamp1 := TexturedStreetLamp()
	scrn.AddModelToShader(lamp1, shaderProgramTextureMat)
	lamp2 := StreetLamp()
	scrn.AddModelToShader(lamp2, shaderProgramMaterial)
	scrn.AddSpotLightSource(lamp1.GetLightSource(), [10]string{
		"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient",
		"spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant",
		"spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	scrn.AddSpotLightSource(lamp2.GetLightSource(), [10]string{
		"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient",
		"spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant",
		"spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})

	Bug1 = model.NewBug(
		Settings["BugPosition"].GetCurrentValue().(mgl32.Vec3),
		Settings["BugScale"].GetCurrentValue().(mgl32.Vec3),
		glWrapper)
	Bug1.SetDirection(mgl32.Vec3{1, 0, 0})
	Bug1.SetSpeed(Settings["BugVelocity"].GetCurrentValue().(float32))

	scrn.AddModelToShader(Bug1, shaderProgramMaterial)

	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	TexModel.AddMesh(TexturedBug(sunTexture))
	scrn.AddModelToShader(TexModel, shaderProgramTexture)
	scrn.AddModelToShader(MatModel, shaderProgramMaterial)
	scrn.Setup(setupApp)
	return scrn
}
func SkipFormScreen() bool {
	val := os.Getenv(FORM_ENV_NAME)
	if val == OFF_VALUE {
		return true
	}
	return false
}

func main() {
	runtime.LockOSThread()
	InitSettings()

	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	AppScreen = mainScreen()
	app.AddScreen(AppScreen)
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	if SkipFormScreen() {
		app.ActivateScreen(AppScreen)
	} else {
		app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
		app.GetWindow().SetCharCallback(app.CharCallback)
		MenuScreen = createMenu()
		app.AddScreen(MenuScreen)
		app.MenuScreen(MenuScreen)
		SettingsScreen = createSettings(Settings)
		app.AddScreen(SettingsScreen)
		app.ActivateScreen(SettingsScreen)
	}
	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
