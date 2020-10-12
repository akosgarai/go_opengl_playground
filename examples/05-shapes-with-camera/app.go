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
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/cylinder"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - shapes with camera"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app *application.Application

	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	lastUpdate int64
	startTime  int64

	glWrapper glwrapper.Wrapper
)

func init() {
	runtime.LockOSThread()

	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	Settings.AddConfig("Color1", "Cube color 1", "The color of the 1. side of the cube.", mgl32.Vec3{1.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("Color2", "Cube color 2", "The color of the 2. side of the cube.", mgl32.Vec3{1.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("Color3", "Cube color 3", "The color of the 3. side of the cube.", mgl32.Vec3{1.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("Color4", "Cube color 4", "The color of the 4. side of the cube.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("Color5", "Cube color 5", "The color of the 5. side of the cube.", mgl32.Vec3{0.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("Color6", "Cube color 6", "The color of the 6. side of the cube.", mgl32.Vec3{0.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("CubePosition", "Cube position", "The position of the cube.", mgl32.Vec3{-0.5, -0.5, 0.5}, nil)
	// sphere
	Settings.AddConfig("SphereColor", "Sphere color", "The color of the sphere.", mgl32.Vec3{0.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("SpherePosition", "Sphere position", "The position of the sphere.", mgl32.Vec3{3, 3, 5}, nil)
	Settings.AddConfig("SphereScale", "Sphere scale", "The scale of the sphere.", mgl32.Vec3{2, 2, 2}, nil)
	Settings.AddConfig("SpherePrec", "Sphere prec", "The precision of the sphere mesh.", int(30), nil)
	// cylinder
	Settings.AddConfig("CylinderColor", "Cylinder color", "The color of the cylinder.", mgl32.Vec3{1.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("CylinderPosition", "Cylinder position", "The position of the cylinder.", mgl32.Vec3{-3, -3, 5}, nil)
	Settings.AddConfig("CylinderRad", "Cylinder rad", "The radius of the cylinder mesh.", float32(1.5), nil)
	Settings.AddConfig("CylinderLength", "Cylinder length", "The length of the cylinder mesh.", float32(3.0), nil)
	Settings.AddConfig("CylinderPrec", "Cylinder prec", "The precision of the cylinder mesh.", int(30), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{-3, -5, 18.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-90.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.1), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(100.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.005), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
}

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["forward"] = []glfw.Key{glfw.KeyW}
	cm["back"] = []glfw.Key{glfw.KeyS}
	cm["up"] = []glfw.Key{glfw.KeyQ}
	cm["down"] = []glfw.Key{glfw.KeyE}
	cm["left"] = []glfw.Key{glfw.KeyA}
	cm["right"] = []glfw.Key{glfw.KeyD}
	cm["rotateOnEdgeDistance"] = Settings["CameraRotationEdge"].GetCurrentValue().(float32)
	cm["mode"] = "default"
	return cm
}

// It creates a new camera with the necessary setup from settings screen
func CreateCameraFromSettings() *camera.DefaultCamera {
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

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	app.Update(delta)
	lastUpdate = nowNano
}

// It generates a cube.
func CreateCubeMesh() *mesh.ColorMesh {
	colors := []mgl32.Vec3{
		Settings["Color1"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color2"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color3"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color4"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color5"].GetCurrentValue().(mgl32.Vec3),
		Settings["Color6"].GetCurrentValue().(mgl32.Vec3),
	}
	cube := cuboid.NewCube()
	v, i, _ := cube.ColoredMeshInput(colors)
	m := mesh.NewColorMesh(v, i, colors, glWrapper)
	m.SetPosition(Settings["CubePosition"].GetCurrentValue().(mgl32.Vec3))
	return m
}

// It generates a Sphere.
func CreateSphereMesh() *mesh.ColorMesh {
	s := sphere.New(20)
	cols := []mgl32.Vec3{Settings["SphereColor"].GetCurrentValue().(mgl32.Vec3)}
	v, i, _ := s.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(Settings["SpherePosition"].GetCurrentValue().(mgl32.Vec3))
	m.SetScale(Settings["SphereScale"].GetCurrentValue().(mgl32.Vec3))
	return m
}

// It generates a cylinder
func CreateCylinder() *mesh.ColorMesh {
	rad := Settings["CylinderRad"].GetCurrentValue().(float32)
	len := Settings["CylinderLength"].GetCurrentValue().(float32)
	prec := Settings["CylinderPrec"].GetCurrentValue().(int)
	c := cylinder.New(rad, prec, len)
	cols := []mgl32.Vec3{Settings["CylinderColor"].GetCurrentValue().(mgl32.Vec3)}
	v, i, _ := c.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(Settings["CylinderPosition"].GetCurrentValue().(mgl32.Vec3))
	return m
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
		"Color1",
		"Color2",
		"Color3",
		"Color4",
		"Color5",
		"Color6",
		"CubePosition",
		"SphereColor",
		"SpherePosition",
		"SphereScale",
		"SpherePrec",
		"CylinderColor",
		"CylinderPosition",
		"CylinderRad", "CylinderLength",
		"CylinderPrec",

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
func GenerateModel() *model.BaseModel {
	mod := model.New()
	mod.AddMesh(CreateCubeMesh())
	mod.AddMesh(CreateSphereMesh())
	mod.AddMesh(CreateCylinder())
	return mod
}

func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCameraFromSettings(), CameraMovementOptions())

	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)
	scrn.AddModelToShader(GenerateModel(), shaderProgram)
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
