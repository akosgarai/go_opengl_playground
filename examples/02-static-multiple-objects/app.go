package main

import (
	"os"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/primitives/triangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth   = 800
	WindowHeight  = 800
	WindowTitle   = "Example - static triangle and square"
	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app *application.Application

	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	glWrapper glwrapper.Wrapper

	lastUpdate int64
	startTime  int64
)

func InitSettings() {
	var colorValidator, widthValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	widthValidator = func(f float32) bool { return f >= 0 && f < 2 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("SquareColor", "Square color", "The color of the square item.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("TriangleColor", "Triangle color", "The color of the triangle item.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("TriangleScale", "Triangle scale", "The scale of the triangle item.", mgl32.Vec3{0.5, 0.5, 0.5}, nil)
	Settings.AddConfig("SquareScale", "Square scale", "The scale of the square item.", mgl32.Vec3{0.5, 0.5, 0.5}, nil)
	Settings.AddConfig("TrianglePosition", "Triangle position", "The position of the triangle item.", mgl32.Vec3{-0.4, 0.2, 0}, nil)
	Settings.AddConfig("SquarePosition", "Square position", "The position of the square item.", mgl32.Vec3{0.4, -0.2, 0}, nil)
	Settings.AddConfig("Width", "Width", "The width of the square.", float32(1.0), widthValidator)
}
func GenerateModel() *model.BaseModel {
	mod := model.New()
	col := []mgl32.Vec3{Settings["SquareColor"].GetCurrentValue().(mgl32.Vec3)}
	width := Settings["Width"].GetCurrentValue().(float32)
	square := rectangle.NewExact(width, width)
	v, i, _ := square.ColoredMeshInput(col)
	squareMesh := mesh.NewColorMesh(v, i, col, glWrapper)
	squareMesh.SetScale(Settings["SquareScale"].GetCurrentValue().(mgl32.Vec3))
	squareMesh.SetPosition(Settings["SquarePosition"].GetCurrentValue().(mgl32.Vec3))
	mod.AddMesh(squareMesh)

	col = []mgl32.Vec3{Settings["TriangleColor"].GetCurrentValue().(mgl32.Vec3)}
	triang := triangle.New(30, 60, 90)
	v, indices, _ := triang.ColoredMeshInput(col)
	triangleMesh := mesh.NewColorMesh(v, indices, col, glWrapper)
	triangleMesh.SetScale(Settings["TriangleScale"].GetCurrentValue().(mgl32.Vec3))
	triangleMesh.SetPosition(Settings["TrianglePosition"].GetCurrentValue().(mgl32.Vec3))
	mod.AddMesh(triangleMesh)
	mod.RotateX(90)
	return mod
}
func mainScreen() *screen.Screen {
	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)

	scrn.AddModelToShader(GenerateModel(), shaderProgram)
	scrn.Setup(setupApp)
	return scrn
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
func AddFormScreen() bool {
	val := os.Getenv(FORM_ENV_NAME)
	if val == ON_VALUE {
		return true
	}
	return false
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
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
func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"ClearCol",
		"SquareColor",
		"SquareScale",
		"SquarePosition",
		"Width",
		"TriangleColor",
		"TriangleScale",
		"TrianglePosition",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Settings")
}
func main() {
	runtime.LockOSThread()

	InitSettings()
	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := mainScreen()
	app.AddScreen(scrn)
	if AddFormScreen() {
		app.GetWindow().SetKeyCallback(app.KeyCallback)
		app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
		app.GetWindow().SetCharCallback(app.CharCallback)
		MenuScreen = createMenu()
		app.AddScreen(MenuScreen)
		app.MenuScreen(MenuScreen)
		SettingsScreen = createSettings(Settings)
		app.AddScreen(SettingsScreen)
		app.ActivateScreen(SettingsScreen)
	} else {
		app.ActivateScreen(scrn)
	}

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		Update()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
