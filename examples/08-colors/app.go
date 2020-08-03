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
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - cubes with light source"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	lastUpdate int64
	startTime  int64

	glWrapper glwrapper.Wrapper
)

func InitSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	// light source
	Settings.AddConfig("LSAmbient", "Light ambient", "The ambient color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSDiffuse", "Light diffuse", "The diffuse color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSSpecular", "Light specular", "The specular color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSPosition", "Light position", "The position vector of the lightsource.", mgl32.Vec3{-3.0, 0.0, -3.0}, nil)
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	// first cube
	Settings.AddConfig("Color1", "Cube color 1", "The color of the 1. side of the cube.", mgl32.Vec3{1.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("Color2", "Cube color 2", "The color of the 2. side of the cube.", mgl32.Vec3{1.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("Color3", "Cube color 3", "The color of the 3. side of the cube.", mgl32.Vec3{1.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("Color4", "Cube color 4", "The color of the 4. side of the cube.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("Color5", "Cube color 5", "The color of the 5. side of the cube.", mgl32.Vec3{0.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("Color6", "Cube color 6", "The color of the 6. side of the cube.", mgl32.Vec3{0.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("CubePosition", "Cube position", "The position of the cube.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("WhiteCubePosition", "White Cube pos", "The position of the white cube.", mgl32.Vec3{-3.0, -0.5, -3.0}, nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0, 0, 10.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-90.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
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

func CreateColoredCubeMesh(pos mgl32.Vec3, col []mgl32.Vec3) *mesh.ColorMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.ColoredMeshInput(col)
	m := mesh.NewColorMesh(v, i, col, glWrapper)
	m.SetPosition(pos)
	return m
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
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
		"LSAmbient",
		"LSDiffuse",
		"LSSpecular",
		"LSPosition",
		"LSConstantTerm", "LSLinearTerm",
		"LSQuadraticTerm",

		"Color1",
		"Color2",
		"Color3",
		"Color4",
		"Color5",
		"Color6",
		"CubePosition",
		"WhiteCubePosition",

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

	PointLightSource := light.NewPointLight([4]mgl32.Vec3{
		Settings["LSPosition"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSSpecular"].GetCurrentValue().(mgl32.Vec3)},
		[3]float32{
			Settings["LSConstantTerm"].GetCurrentValue().(float32),
			Settings["LSLinearTerm"].GetCurrentValue().(float32),
			Settings["LSQuadraticTerm"].GetCurrentValue().(float32),
		})

	// Add the lightources to the application
	scrn.AddPointLightSource(PointLightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "light.constant", "light.linear", "light.quadratic"})

	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)
	Model := model.New()
	whiteCube := CreateColoredCubeMesh(Settings["WhiteCubePosition"].GetCurrentValue().(mgl32.Vec3), []mgl32.Vec3{mgl32.Vec3{1.0, 1.0, 1.0}})
	Model.AddMesh(whiteCube)
	colors := []mgl32.Vec3{
		Settings["Color1"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color2"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color3"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color4"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color5"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color6"].GetCurrentValue().(mgl32.Vec3),
	}
	coloredCube := CreateColoredCubeMesh(Settings["CubePosition"].GetCurrentValue().(mgl32.Vec3), colors)
	Model.AddMesh(coloredCube)
	scrn.AddModelToShader(Model, shaderProgram)
	scrn.Setup(setupApp)
	return scrn
}
func AddFormScreen() bool {
	val := os.Getenv(FORM_ENV_NAME)
	if val == ON_VALUE {
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
	if AddFormScreen() {
		app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
		app.GetWindow().SetCharCallback(app.CharCallback)
		MenuScreen = createMenu()
		app.AddScreen(MenuScreen)
		app.MenuScreen(MenuScreen)
		SettingsScreen = createSettings(Settings)
		app.AddScreen(SettingsScreen)
		app.ActivateScreen(SettingsScreen)
	} else {
		app.ActivateScreen(AppScreen)
	}
	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
