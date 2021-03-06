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
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
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
	WindowTitle  = "Example - plane with ball"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app            *application.Application
	Ball           *mesh.ColorMesh
	Ground         *mesh.ColorMesh
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	lastUpdate int64
	startTime  int64

	BallInitialDirection = mgl32.Vec3{0, 1, 0}

	glWrapper glwrapper.Wrapper
)

func init() {
	runtime.LockOSThread()

	var colorValidator, heightValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	heightValidator = func(f float32) bool { return f >= 0.0 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	Settings.AddConfig("SphereColor", "Sphere Color", "The color of the sphere.", mgl32.Vec3{1.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("SpherePosition", "Sphere position", "The position of the sphere item.", mgl32.Vec3{0.0, 5.0, 0.0}, nil)
	Settings.AddConfig("SphereScale", "Sphere scale", "The scale of the sphere item.", mgl32.Vec3{2.0, 2.0, 2.0}, nil)
	Settings.AddConfig("SphereSpeed", "Sphere speed", "The velocity of the sphere item.", float32(0.02), nil)
	Settings.AddConfig("SpherePrecision", "Sphere precision", "The precision of the sphere item.", int(10), nil)
	Settings.AddConfig("SphereMaxHeight", "Sphere height", "The top position of the sphere item.", float32(10.0), heightValidator)
	Settings.AddConfig("SquareColor", "Surface Color", "The color of the square surface.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("SquareScale", "Surface scale", "The scale of the square surface.", mgl32.Vec3{40.0, 40.0, 40.0}, nil)

	Settings.AddConfig("SquarePosition", "Square position", "The position of the square item.", mgl32.Vec3{0.0, 0.0, 0}, nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0.0, 5.0, -24.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, -1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(90.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.1), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(100.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.01), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.005), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
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

func CreateSphereMesh() *mesh.ColorMesh {
	s := sphere.New(Settings["SpherePrecision"].GetCurrentValue().(int))
	cols := []mgl32.Vec3{Settings["SphereColor"].GetCurrentValue().(mgl32.Vec3)}
	v, i, _ := s.ColoredMeshInput(cols)
	m := mesh.NewColorMesh(v, i, cols, glWrapper)
	m.SetPosition(Settings["SpherePosition"].GetCurrentValue().(mgl32.Vec3))
	m.SetScale(Settings["SphereScale"].GetCurrentValue().(mgl32.Vec3))
	m.SetDirection(BallInitialDirection)
	m.SetSpeed(Settings["SphereSpeed"].GetCurrentValue().(float32))
	return m
}

// It generates a square.
func CreateSquareMesh() *mesh.ColorMesh {
	squareColor := []mgl32.Vec3{Settings["SquareColor"].GetCurrentValue().(mgl32.Vec3)}
	s := rectangle.NewSquare()
	v, i, bo := s.ColoredMeshInput(squareColor)
	m := mesh.NewColorMesh(v, i, squareColor, glWrapper)
	m.SetScale(Settings["SquareScale"].GetCurrentValue().(mgl32.Vec3))
	m.SetPosition(Settings["SquarePosition"].GetCurrentValue().(mgl32.Vec3))
	m.SetBoundingObject(bo)
	return m
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

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	// handle ball
	if Ball.GetPosition().Y() >= Settings["SphereMaxHeight"].GetCurrentValue().(float32) {
		Ball.SetPosition(mgl32.Vec3{Ball.GetPosition().X(), Settings["SphereMaxHeight"].GetCurrentValue().(float32), Ball.GetPosition().Z()})
		Ball.SetDirection(BallInitialDirection.Mul(-1.0))
	}
	bottomPos := Settings["SquarePosition"].GetCurrentValue().(mgl32.Vec3).Y() + Settings["SphereScale"].GetCurrentValue().(mgl32.Vec3).Y()
	if Ball.GetPosition().Y() <= bottomPos {
		Ball.SetPosition(mgl32.Vec3{Ball.GetPosition().X(), bottomPos, Ball.GetPosition().Z()})
		Ball.SetDirection(BallInitialDirection)
	}
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
		"SphereColor",
		"SpherePosition",
		"SphereScale",
		"SphereSpeed", "SpherePrecision",
		"SquareColor",
		"SquareScale",
		"SquarePosition",

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
	Ball = CreateSphereMesh()
	mod.AddMesh(Ball)
	Ground = CreateSquareMesh()
	mod.AddMesh(Ground)
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
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
