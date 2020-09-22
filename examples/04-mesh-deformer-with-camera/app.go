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
	"github.com/akosgarai/playground_engine/pkg/primitives/triangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - mesh deformer - with moving camera"

	rows   = 10
	cols   = 10
	length = 10

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

func InitSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("TriangleColorFront", "Front Color", "The front color of the triangle.", mgl32.Vec3{0.0, 0.0, 1.0}, colorValidator)
	Settings.AddConfig("TriangleColorBack", "Back Color", "The back color of the square.", mgl32.Vec3{0.0, 0.5, 1.0}, colorValidator)
	Settings.AddConfig("Rows", "Triangle/row", "The number of the triangles in one row.", int(10), nil)
	Settings.AddConfig("Columns", "Triangle/column", "The number of the triangles in one column.", int(10), nil)
	Settings.AddConfig("Length", "Side length", "The length of the longest side of the triangle.", float32(10.0), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{-4.5, -3.0, 8}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The initial position of the camera.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(0.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(3.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.01), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(200.0), nil)
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

// It generates a bunch of triangles and sets their color to static blue.
func GenerateModel() *model.BaseModel {
	Model := model.New()
	triang := triangle.New(90, 45, 45)
	frontColor := []mgl32.Vec3{Settings["TriangleColorFront"].GetCurrentValue().(mgl32.Vec3)}
	backColor := []mgl32.Vec3{Settings["TriangleColorBack"].GetCurrentValue().(mgl32.Vec3)}
	v1, indices1, _ := triang.ColoredMeshInput(frontColor)
	v2, indices2, _ := triang.ColoredMeshInput(backColor)
	rows := Settings["Rows"].GetCurrentValue().(int)
	cols := Settings["Columns"].GetCurrentValue().(int)
	length := Settings["Length"].GetCurrentValue().(float32)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j) * length
			topZ := float32(i) * length
			topY := float32(0.0)

			m := mesh.NewColorMesh(v1, indices1, frontColor, glWrapper)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetScale(mgl32.Vec3{length, length, length})

			Model.AddMesh(m)

			m2 := mesh.NewColorMesh(v2, indices2, backColor, glWrapper)
			m2.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.RotateY(180)

			Model.AddMesh(m2)
		}
	}
	return Model
}

// Update the z coordinates of the vectors.
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	app.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
	app.Update(delta)
	lastUpdate = nowNano
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
		"TriangleColorFront",
		"TriangleColorBack",
		"Rows", "Columns",
		"Length",
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
func AddFormScreen() bool {
	val := os.Getenv(FORM_ENV_NAME)
	if val == ON_VALUE {
		return true
	}
	return false
}
func mainScreen() *screen.Screen {
	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.SetCamera(CreateCameraFromSettings())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(Settings["CameraRotationEdge"].GetCurrentValue().(float32))
	scrn.AddShader(shaderProgram)

	scrn.AddModelToShader(GenerateModel(), shaderProgram)
	scrn.Setup(setupApp)
	return scrn
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
