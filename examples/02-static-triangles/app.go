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
	WindowTitle   = "Example - static triangles, lots of them"
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
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("ItemColor", "Item color", "The color of the item.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	Settings.AddConfig("Width", "Rows", "The number of triangles in one row.", 4, nil)
	Settings.AddConfig("Length", "Cols", "The number of triangles in on column.", 4, nil)
}

// GenerateTrianglesModel returns a model filled up with the triangles.
func GenerateTrianglesModel() *model.BaseModel {
	mod := model.New()
	triang := triangle.New(60, 60, 60)
	col := []mgl32.Vec3{Settings["ItemColor"].GetCurrentValue().(mgl32.Vec3)}
	v, indicies, _ := triang.ColoredMeshInput(col)

	rows := Settings["Width"].GetCurrentValue().(int)
	cols := Settings["Length"].GetCurrentValue().(int)
	length := 2.0 / float32(cols)
	width := 2.0 / float32(rows)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := 1.0 - (float32(j) * length)
			topY := -1.0 + (float32(i) * width)
			topZ := float32(0.0)

			m := mesh.NewColorMesh(v, indicies, col, glWrapper)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetScale(mgl32.Vec3{length, length, length})
			mod.AddMesh(m)

		}
	}
	mod.RotateX(90)
	return mod
}
func mainScreen() *screen.Screen {
	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)

	scrn.AddModelToShader(GenerateTrianglesModel(), shaderProgram)
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
	app.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
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
		"ItemColor",
		"Width", "Length",
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
