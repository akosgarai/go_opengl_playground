package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - frame screen generator"

	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	AppScreen      *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	glWrapper      glwrapper.Wrapper
	lastUpdate     int64
	startTime      int64
	Settings       = config.New()

	DirectionalLightDirection = mgl32.Vec3{0.5, 0.5, -0.7}.Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightDiffuse   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightSpecular  = mgl32.Vec3{0.3, 0.3, 0.3}
)

func InitSettings() {
	Settings.AddConfig("FrameWidth", "Width", "The width of the screen.", float32(2.0), nil)
	Settings.AddConfig("FrameLength", "Frame Width", "The width of the frame.", float32(0.02), nil)
	Settings.AddConfig("FrameLeftOffset", "Left width", "The width of the left offset.", float32(0.1), nil)
	Settings.AddConfig("LabelColor", "Header Label Col", "The color of the header label text.", mgl32.Vec3{1, 0, 0}, nil)
	Settings.AddConfig("ItemLabelColor", "Item Label Col", "The color of the item label text.", mgl32.Vec3{0, 1, 0}, nil)
	Settings.AddConfig("ItemInputColor", "Item Input Col", "The color of the item input text.", mgl32.Vec3{0, 0, 0}, nil)
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func createMenu() *screen.MenuScreen {
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 1, "paper", glWrapper)
	builder := screen.NewMenuScreenBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	builder.SetFrameMaterial(material.Ruby)
	builder.SetMenuItemSurfaceTexture(tex)
	builder.SetMenuItemFontColor(mgl32.Vec3{0, 0, 0})

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
		AppScreen = createGame(Settings, SettingsScreen)
		app.ActivateScreen(AppScreen)
	}
	startEvent := func() {
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = createGame(Settings, SettingsScreen)
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
		"FrameWidth", "FrameLength",
		"FrameLeftOffset",
		"LabelColor",
		"ItemLabelColor",
		"ItemInputColor",
	}
	builder := screen.NewFormScreenBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	builder.SetFrameMaterial(material.Ruby)
	builder.SetHeaderLabel("Settings")
	builder.SetConfig(Settings)
	builder.SetConfigOrder(formItemOrders)
	return builder.Build()
}
func createGame(conf config.Config, form *screen.FormScreen) *screen.FormScreen {
	formItemOrders := []string{
		"FrameWidth", "FrameLength",
		"FrameLeftOffset", "TextWidth",
	}
	builder := screen.NewFormScreenBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	builder.SetFrameMaterial(material.Ruby)
	builder.SetFrameSize(conf["FrameWidth"].GetCurrentValue().(float32), conf["FrameLength"].GetCurrentValue().(float32), conf["FrameLeftOffset"].GetCurrentValue().(float32))
	builder.SetHeaderLabelColor(conf["LabelColor"].GetCurrentValue().(mgl32.Vec3))
	builder.SetFormItemLabelColor(conf["ItemLabelColor"].GetCurrentValue().(mgl32.Vec3))
	builder.SetFormItemInputColor(conf["ItemInputColor"].GetCurrentValue().(mgl32.Vec3))
	builder.SetHeaderLabel("Settings")
	builder.SetConfig(Settings)
	builder.SetConfigOrder(formItemOrders)
	AppScreen := builder.Build()

	return AppScreen
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
	app.Update(delta)
}

func main() {
	runtime.LockOSThread()
	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	InitSettings()
	MenuScreen = createMenu()
	SettingsScreen = createSettings(Settings)
	AppScreen = createGame(Settings, SettingsScreen)

	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.SetWrapper(glWrapper)
	app.GetWindow().SetCharCallback(app.CharCallback)
	app.AddScreen(AppScreen)
	app.AddScreen(MenuScreen)
	app.AddScreen(SettingsScreen)
	app.MenuScreen(MenuScreen)
	app.ActivateScreen(MenuScreen)

	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		Update()
		app.GetWindow().SwapBuffers()
	}
}
