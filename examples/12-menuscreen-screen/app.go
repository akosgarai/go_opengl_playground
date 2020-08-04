package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/theme"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - theme editor for menu and form screens"
	FontFile     = "/assets/fonts/Desyrel/desyrel.regular.ttf"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
	LEFT_MOUSE_BUTTON    = glfw.MouseButtonLeft
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	FormScreenApp  *screen.FormScreen
	MenuScreenApp  *screen.MenuScreen
	Settings       = config.New()
	TestSettings   = config.New()

	glWrapper  glwrapper.Wrapper
	lastUpdate int64

	MaterialMap = map[string]*material.Material{
		"Emerald":       material.Emerald,
		"Jade":          material.Jade,
		"Obsidian":      material.Obsidian,
		"Pearl":         material.Pearl,
		"Ruby":          material.Ruby,
		"Turquoise":     material.Turquoise,
		"Brass":         material.Brass,
		"Bronze":        material.Bronze,
		"Chrome":        material.Chrome,
		"Copper":        material.Copper,
		"Gold":          material.Gold,
		"Silver":        material.Silver,
		"Blackplastic":  material.Blackplastic,
		"Cyanplastic":   material.Cyanplastic,
		"Greenplastic":  material.Greenplastic,
		"Redplastic":    material.Redplastic,
		"Whiteplastic":  material.Whiteplastic,
		"Yellowplastic": material.Yellowplastic,
		"Blackrubber":   material.Blackrubber,
		"Cyanrubber":    material.Cyanrubber,
		"Greenrubber":   material.Greenrubber,
		"Redrubber":     material.Redrubber,
		"Whiterubber":   material.Whiterubber,
		"Yellowrubber":  material.Yellowrubber,
	}
)

func InitThemeSettings() {
	Settings.AddConfig("FrameWidth", "Width", "The width of the screen.", float32(2.0), nil)
	Settings.AddConfig("FrameLength", "Frame Width", "The width of the frame.", float32(0.02), nil)
	Settings.AddConfig("FrameLeftOffset", "Left width", "The width of the left offset.", float32(0.1), nil)
	Settings.AddConfig("DCBHeight", "DCB height", "The height of the detail content box.", float32(0.4), nil)
	Settings.AddConfig("HeaderLabelColor", "Header Label Col", "The color of the header label text.", mgl32.Vec3{1, 0, 0}, nil)
	Settings.AddConfig("ItemLabelColor", "Item Label Col", "The color of the item label text.", mgl32.Vec3{0, 1, 0}, nil)
	Settings.AddConfig("ItemInputColor", "Item Input Col", "The color of the item input text.", mgl32.Vec3{0, 0, 0}, nil)
	Settings.AddConfig("BGColor", "Background Col", "The color of the background.", mgl32.Vec3{0, 0, 0}, nil)
	materialList := ""
	for k, _ := range MaterialMap {
		materialList = materialList + k + ", "
	}
	Settings.AddConfig("FrameMaterial", "Frame material", "The material of the frame. Options: "+materialList, "Jade", nil)
	Settings.AddConfig("DefaultMat", "Default material", "The material of the items. Options: "+materialList, "Whiteplastic", nil)
	Settings.AddConfig("HoverMat", "Hover material", "The material of the hovered items. Options: "+materialList, "Ruby", nil)
}

func InitAppSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	TestSettings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	TestSettings.AddConfig("Color1", "Cube color 1", "The color of the 1. side of the cube.", mgl32.Vec3{1.0, 1.0, 0.0}, colorValidator)
	TestSettings.AddConfig("Color2", "Cube color 2", "The color of the 2. side of the cube.", mgl32.Vec3{1.0, 0.0, 1.0}, colorValidator)
	TestSettings.AddConfig("Color3", "Cube color 3", "The color of the 3. side of the cube.", mgl32.Vec3{1.0, 0.0, 0.0}, colorValidator)
	TestSettings.AddConfig("TestBool05", "Bool param", "For checking the bool params with different theme", false, nil)
	TestSettings.AddConfig("Color4", "Cube color 4", "The color of the 4. side of the cube.", mgl32.Vec3{0.0, 1.0, 0.0}, colorValidator)
	TestSettings.AddConfig("TestBool04", "Bool param", "For checking the bool params with different theme", false, nil)
	TestSettings.AddConfig("TestBool06", "Bool param", "For checking the bool params with different theme", false, nil)
	TestSettings.AddConfig("Color5", "Cube color 5", "The color of the 5. side of the cube.", mgl32.Vec3{0.0, 1.0, 1.0}, colorValidator)
	TestSettings.AddConfig("TestBool01", "Bool param", "For checking the bool params with different theme", false, nil)
	TestSettings.AddConfig("TestBool02", "Bool param", "For checking the bool params with different theme", true, nil)
	TestSettings.AddConfig("TestBool03", "Bool param", "For checking the bool params with different theme", false, nil)
	TestSettings.AddConfig("Color6", "Cube color 6", "The color of the 6. side of the cube.", mgl32.Vec3{0.0, 0.0, 1.0}, colorValidator)
	TestSettings.AddConfig("CubePosition", "Cube position", "The position of the cube.", mgl32.Vec3{-0.5, -0.5, 0.5}, nil)
	TestSettings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0, 0, 10.0}, nil)
	TestSettings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	TestSettings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-90.0), nil)
	TestSettings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	TestSettings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.1), nil)
	TestSettings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(100.0), nil)
	TestSettings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	TestSettings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	TestSettings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.005), nil)
	TestSettings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"FrameWidth", "FrameLength",
		"FrameLeftOffset", "DCBHeight",
		"HeaderLabelColor",
		"ItemLabelColor",
		"ItemInputColor",
		"BGColor",
		"FrameMaterial",
		"DefaultMat",
		"HoverMat",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Theme editor")
}
func createMenu() *screen.MenuScreen {
	contAll := func(m map[string]bool) bool {
		return true
	}
	startFormEvent := func() {
		FormScreenApp = formScreenApp(TestSettings)
		app.ActivateScreen(FormScreenApp)
	}
	startMenuEvent := func() {
		MenuScreenApp = menuScreenApp()
		app.ActivateScreen(MenuScreenApp)
	}
	settingsEvent := func() {
		app.ActivateScreen(SettingsScreen)
	}
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	options := []screen.Option{
		*screen.NewMenuScreenOption("Menu screen", contAll, startMenuEvent),
		*screen.NewMenuScreenOption("Form screen", contAll, startFormEvent),
		*screen.NewMenuScreenOption("Settings", contAll, settingsEvent),
		*screen.NewMenuScreenOption("Exit", contAll, exitEvent),
	}
	return app.BuildMenuScreen(options)
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

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
}
func themeFromSettings() theme.Theme {
	var t theme.Theme
	t.SetBackgroundColor(Settings["BGColor"].GetCurrentValue().(mgl32.Vec3))
	t.SetHeaderLabelColor(Settings["HeaderLabelColor"].GetCurrentValue().(mgl32.Vec3))
	t.SetInputColor(Settings["ItemInputColor"].GetCurrentValue().(mgl32.Vec3))
	t.SetLabelColor(Settings["ItemLabelColor"].GetCurrentValue().(mgl32.Vec3))
	t.SetDetailContentBoxHeight(Settings["DCBHeight"].GetCurrentValue().(float32))
	t.SetFrameLength(Settings["FrameLength"].GetCurrentValue().(float32))
	t.SetFrameWidth(Settings["FrameWidth"].GetCurrentValue().(float32))
	t.SetFrameTopLeftWidth(Settings["FrameLeftOffset"].GetCurrentValue().(float32))
	t.SetFrameMaterial(MaterialMap[Settings["FrameMaterial"].GetCurrentValue().(string)])
	t.SetMenuItemDefaultMaterial(MaterialMap[Settings["DefaultMat"].GetCurrentValue().(string)])
	t.SetMenuItemHoverMaterial(MaterialMap[Settings["HoverMat"].GetCurrentValue().(string)])
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 1, "paper", glWrapper)
	t.SetMenuItemSurfaceTexture(tex)
	return t
}
func menuScreenApp() *screen.MenuScreen {
	app.SetTheme(themeFromSettings())
	stateFunction := func(m map[string]bool) bool {
		return true
	}
	eventFunction := func() {
	}
	options := []screen.Option{
		*screen.NewMenuScreenOption("continue", stateFunction, eventFunction),
		*screen.NewMenuScreenOption("start", stateFunction, eventFunction),
		*screen.NewMenuScreenOption("restart", stateFunction, eventFunction),
		*screen.NewMenuScreenOption("settings", stateFunction, eventFunction),
		*screen.NewMenuScreenOption("exit", stateFunction, eventFunction),
	}
	return app.BuildMenuScreen(options)
}
func formScreenApp(defaults config.Config) *screen.FormScreen {
	app.SetTheme(themeFromSettings())
	formItemOrders := []string{
		"ClearCol",
		"Color1",
		"Color2",
		"Color3",
		"TestBool06",
		"Color4",
		"TestBool05",
		"TestBool04",
		"Color5",
		"TestBool03",
		"TestBool02",
		"TestBool01",
		"Color6",
		"CubePosition",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Settings")
}

func main() {
	runtime.LockOSThread()
	InitThemeSettings()
	InitAppSettings()

	app = application.New(glWrapper)
	Window := window.InitGlfw(WindowWidth, WindowHeight, WindowTitle)
	app.SetWindow(Window)
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetWrapper(glWrapper)
	app.SetTheme(*theme.Dark)
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetCharCallback(app.CharCallback)

	MenuScreen = createMenu()
	SettingsScreen = createSettings(Settings)
	app.AddScreen(MenuScreen)
	app.AddScreen(SettingsScreen)
	app.MenuScreen(MenuScreen)
	app.ActivateScreen(MenuScreen)

	FormScreenApp = formScreenApp(TestSettings)
	MenuScreenApp = menuScreenApp()
	app.AddScreen(MenuScreenApp)
	app.AddScreen(FormScreenApp)

	lastUpdate = time.Now().UnixNano()

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		Update()
		app.GetWindow().SwapBuffers()
	}
}
