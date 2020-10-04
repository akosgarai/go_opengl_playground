package main

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/transformations"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	FontFile = "/assets/fonts/Desyrel/desyrel.regular.ttf"
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	glWrapper      glwrapper.Wrapper
	lastUpdate     int64
	startTime      int64
	Settings       = config.New()

	Builder          *window.WindowBuilder
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowTitle      = "Example - terrain generator with settings"
	WindowFullScreen = false
	Aspect           = false

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby

	DirectionalLightDirection = mgl32.Vec3{0.5, 0.5, -0.7}.Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightDiffuse   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightSpecular  = mgl32.Vec3{0.3, 0.3, 0.3}
)

func init() {
	runtime.LockOSThread()

	width := os.Getenv("WIDTH")
	if width != "" {
		val, err := strconv.Atoi(width)
		if err == nil {
			WindowWidth = val
		}
	}
	height := os.Getenv("HEIGHT")
	if height != "" {
		val, err := strconv.Atoi(height)
		if err == nil {
			WindowHeight = val
		}
	}
	decorated := os.Getenv("DECORATED")
	if decorated == "0" {
		WindowDecorated = false
	}
	title := os.Getenv("TITLE")
	if title != "" {
		WindowTitle = title
	}
	Builder = window.NewWindowBuilder()
	fullScreen := os.Getenv("FULL")
	if fullScreen == "1" {
		WindowFullScreen = true
		WindowWidth, WindowHeight = Builder.GetCurrentMonitorResolution()
	}
	aspect := os.Getenv("ASPECT")
	if aspect != "" {
		Aspect = true

	}

	Settings.AddConfig("Width", "Rows (i)", "The width / rows of the heightmap of the terrain.", 4, nil)
	Settings.AddConfig("Length", "Cols (i)", "The length / columns of the heightmap of the terrain.", 4, nil)
	Settings.AddConfig("Iterations", "Iter (i)", "The number of the iterations of the random height generation step.", 10, nil)
	Settings.AddConfig("PeakProb", "Peak (i)", "The probability of a given map position will be a peak.", 5, nil)
	Settings.AddConfig("CliffProb", "Cliff (i)", "The probability of a given map position will be a cliff.", 5, nil)
	Settings.AddConfig("TerrainMinHeight", "MinH (f)", "The minimum height of the terrain surface.", float32(-1.0), nil)
	Settings.AddConfig("TerrainMaxHeight", "MaxH (f)", "The maximum height of the terrain surface.", float32(3.0), nil)
	Settings.AddConfig("TerrainScale", "Scale (f)", "The generated terrain is scaled with the components of this vector.", mgl32.Vec3{5.0, 1.0, 5.0}, nil)
	Settings.AddConfig("TerrainPos", "PosY (f)", "The center / middle point of the terrain mesh.", mgl32.Vec3{0.0, 1.003, 0.0}, nil)
	Settings.AddConfig("Seed", "Seed (i64)", "This value is used as seed for the random number generation, if the random is not set.", int64(0), nil)
	Settings.AddConfig("RandomSeed", "Rand Seed", "If this is set, the seed will be based on the current timestamp.", false, nil)
	Settings.AddConfig("TerrainTexture", "Terr tex", "The texture of the terrain. Currently the 'Grass' is supported.", "Grass", nil)
	Settings.AddConfig("NeedLiquid", "Has Liquid", "If this is set, water will also be generated to the terrain.", true, nil)
	Settings.AddConfig("LiquidEta", "Leta (f)", "The refraction ratio between the air and the liquid surface.", float32(0.75), nil)
	Settings.AddConfig("LiquidAmplitude", "Lampl (f)", "The amplitude of the waves (sin wave) in the liquid surface", float32(0.0625), nil)
	Settings.AddConfig("LiquidFrequency", "LFreq (f)", "The wavelengts of the waves of the liquid surface.", float32(1.0), nil)
	Settings.AddConfig("LiquidDetail", "Ldetail (i)", "The size of the liquid surface is the same as the terrain, but its height map is bigger this times.", 10, nil)
	Settings.AddConfig("WaterLevel", "W Lev (f)", "The water level of the liquid surface.", float32(0.25), nil)
	Settings.AddConfig("LiquidTexture", "Liq tex", "The texture of the liquid surface. Currently the 'Water' is supported.", "Water", nil)
	Settings.AddConfig("Debug", "Debug mode", "Turn debug mode on - off. Currently it does nothing.", false, nil)
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the sky.", mgl32.Vec3{0.2, 0.3, 0.8}, nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0.0, -0.5, 3.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The initial position of the camera.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-85.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.001), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(20.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.0), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.0), nil)
	// - FPS camera
	Settings.AddConfig("CameraFPS", "FPS Camera", "If this flag is true, the camera will be FPS like.", false, nil)
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
func CreateCameraFromSettings(conf config.Config) interfaces.Camera {
	cameraPosition := conf["CameraPos"].GetCurrentValue().(mgl32.Vec3)
	worldUp := conf["WorldUp"].GetCurrentValue().(mgl32.Vec3)
	yawAngle := conf["CameraYaw"].GetCurrentValue().(float32)
	pitchAngle := conf["CameraPitch"].GetCurrentValue().(float32)
	fov := conf["CameraFov"].GetCurrentValue().(float32)
	near := conf["CameraNear"].GetCurrentValue().(float32)
	far := conf["CameraFar"].GetCurrentValue().(float32)
	moveSpeed := conf["CameraVelocity"].GetCurrentValue().(float32)
	directionSpeed := conf["CameraRotation"].GetCurrentValue().(float32)
	if conf["CameraFPS"].GetCurrentValue().(bool) {
		cam := camera.NewFPSCamera(cameraPosition, worldUp, yawAngle, pitchAngle)
		cam.SetupProjection(fov, float32(WindowWidth)/float32(WindowHeight), near, far)
		cam.SetVelocity(moveSpeed)
		cam.SetRotationStep(directionSpeed)
		return cam
	}
	cam := camera.NewCamera(cameraPosition, worldUp, yawAngle, pitchAngle)
	cam.SetupProjection(fov, float32(WindowWidth)/float32(WindowHeight), near, far)
	cam.SetVelocity(moveSpeed)
	cam.SetRotationStep(directionSpeed)
	return cam
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"Width", "Length",
		"Iterations", "PeakProb",
		"CliffProb", "TerrainMinHeight",
		"TerrainMaxHeight",
		"TerrainScale",
		"TerrainPos",
		"RandomSeed", "Seed",
		"TerrainTexture",
		"NeedLiquid", "LiquidEta",
		"LiquidAmplitude", "LiquidFrequency",
		"LiquidDetail", "WaterLevel",
		"LiquidTexture",
		"ClearCol",
		"Debug",
		"WorldUp",
		"CameraPos",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
		"CameraFPS",
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
func createGame(conf config.Config, form *screen.FormScreen) *screen.Screen {
	AppScreen := screen.New()
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShaderBlendingWithFog(glWrapper)
	AppScreen.AddShader(shaderProgramTexture)
	// Shader application for the liquid surface
	shaderProgramLiquid := shader.NewTextureShaderLiquidWithFog(glWrapper)
	AppScreen.AddShader(shaderProgramLiquid)

	AppScreen.SetCamera(CreateCameraFromSettings(conf))
	AppScreen.SetCameraMovementMap(CameraMovementMap())
	AppScreen.SetRotateOnEdgeDistance(conf["CameraRotationEdge"].GetCurrentValue().(float32))

	gb := model.NewTerrainBuilder()
	// terrain related ones
	gb.SetWidth(conf["Width"].GetCurrentValue().(int))
	gb.SetLength(conf["Length"].GetCurrentValue().(int))
	gb.SetIterations(conf["Iterations"].GetCurrentValue().(int))
	gb.SetScale(conf["TerrainScale"].GetCurrentValue().(mgl32.Vec3))
	gb.SetPeakProbability(conf["PeakProb"].GetCurrentValue().(int))
	gb.SetCliffProbability(conf["CliffProb"].GetCurrentValue().(int))
	gb.SetMinHeight(conf["TerrainMinHeight"].GetCurrentValue().(float32))
	gb.SetMaxHeight(conf["TerrainMaxHeight"].GetCurrentValue().(float32))
	gb.SetPosition(conf["TerrainPos"].GetCurrentValue().(mgl32.Vec3))
	if conf["RandomSeed"].GetCurrentValue().(bool) {
		seed := gb.RandomSeed()
		form.SetFormItemValue(form.GetFormItem("Seed"), transformations.Integer64ToString(seed))
		form.SetFormItemValue(form.GetFormItem("RandomSeed"), false)
	} else {
		gb.SetSeed(conf["Seed"].GetCurrentValue().(int64))
	}
	gb.SetGlWrapper(glWrapper)
	switch conf["TerrainTexture"].GetCurrentValue().(string) {
	case "Grass", "grass", "GRASS":
		gb.SurfaceTextureGrass()
		break
	default:
		gb.SurfaceTextureGrass()
		break
	}
	if conf["Debug"].GetCurrentValue().(bool) {
		gb.SetDebugMode(true)
	}
	if conf["NeedLiquid"].GetCurrentValue().(bool) {
		gb.SetLiquidEta(conf["LiquidEta"].GetCurrentValue().(float32))
		gb.SetLiquidAmplitude(conf["LiquidAmplitude"].GetCurrentValue().(float32))
		gb.SetLiquidFrequency(conf["LiquidFrequency"].GetCurrentValue().(float32))
		gb.SetLiquidDetailMultiplier(conf["LiquidDetail"].GetCurrentValue().(int))
		gb.SetLiquidWaterLevel(conf["WaterLevel"].GetCurrentValue().(float32))
		switch conf["LiquidTexture"].GetCurrentValue().(string) {
		case "Water", "water", "WATER":
			gb.LiquidTextureWater()
			break
		default:
			gb.LiquidTextureWater()
		}
		Ground, Water := gb.BuildWithLiquid()
		Water.SetTransparent(true)

		AppScreen.AddModelToShader(Ground, shaderProgramTexture)
		AppScreen.AddModelToShader(Water, shaderProgramLiquid)
	} else {
		Ground := gb.Build()
		AppScreen.AddModelToShader(Ground, shaderProgramTexture)
	}

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	AppScreen.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	AppScreen.Setup(setupApp)

	return AppScreen
}

func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	col := Settings["ClearCol"].GetCurrentValue().(mgl32.Vec3)
	glWrapper.ClearColor(col.X(), col.Y(), col.Z(), 1.0)
	glWrapper.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
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
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)
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

func main() {
	app = application.New(glWrapper)
	Builder.SetFullScreen(WindowFullScreen)
	Builder.SetDecorated(WindowDecorated)
	Builder.SetTitle(WindowTitle)
	Builder.SetWindowSize(WindowWidth, WindowHeight)
	Builder.PrintCurrentMonitorData()

	app.SetWindow(Builder.Build())
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

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
