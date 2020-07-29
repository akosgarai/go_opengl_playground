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
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth   = 800
	WindowHeight  = 800
	WindowTitle   = "Example - static square"
	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
	FontFile      = "/../../assets/fonts/Desyrel/desyrel.regular.ttf"
)

var (
	app *application.Application

	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	glWrapper glwrapper.Wrapper

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby
	lastUpdate        int64
	startTime         int64
)

func InitSettings() {
	var colorValidator, widthValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	widthValidator = func(f float32) bool { return f >= 0 && f < 2 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, colorValidator)
	Settings.AddConfig("ItemColor", "Item color", "The color of the item.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	Settings.AddConfig("Width", "Width", "The width of the square.", float32(1.0), widthValidator)
}

func GenerateModel() *model.BaseModel {
	mod := model.New()
	col := []mgl32.Vec3{Settings["ItemColor"].GetCurrentValue().(mgl32.Vec3)}
	width := Settings["Width"].GetCurrentValue().(float32)
	square := rectangle.NewExact(width, width)
	v, i, _ := square.ColoredMeshInput(col)
	msh := mesh.NewColorMesh(v, i, col, glWrapper)
	mod.AddMesh(msh)
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
	app.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
	app.Update(delta)
}
func createMenu() *screen.MenuScreen {
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 0, "paper", glWrapper)
	builder := screen.NewMenuScreenBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	builder.SetFrameMaterial(DefaultMaterial)
	builder.SetBackgroundColor(DefaultMaterial.GetAmbient())
	builder.SetMenuItemSurfaceTexture(tex)
	builder.SetMenuItemDefaultMaterial(DefaultMaterial)
	builder.SetMenuItemHighlightMaterial(HighlightMaterial)
	builder.SetMenuItemFontColor(mgl32.Vec3{0, 0, 0})

	MenuFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	MenuFonts.SetTransparent(true)
	builder.SetCharset(MenuFonts)
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
	cont := screen.NewMenuScreenOption("continue", contS, continueEvent)
	builder.AddOption(*cont) // continue
	start := screen.NewMenuScreenOption("start", contNS, startEvent)
	builder.AddOption(*start) // start
	restart := screen.NewMenuScreenOption("restart", contS, restartEvent)
	builder.AddOption(*restart) // restart
	settings := screen.NewMenuScreenOption("settings", contAll, settingsEvent)
	builder.AddOption(*settings) // settings
	exit := screen.NewMenuScreenOption("exit", contAll, exitEvent)
	builder.AddOption(*exit) // exit
	return builder.Build()
}
func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"ClearCol",
		"ItemColor",
		"Width",
	}
	builder := screen.NewFormScreenBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	builder.SetFrameMaterial(material.Ruby)
	builder.SetHeaderLabel("Settings")
	builder.SetConfig(defaults)
	builder.SetConfigOrder(formItemOrders)
	return builder.Build()
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
