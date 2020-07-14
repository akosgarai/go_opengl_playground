package main

import (
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - terrain generator with settings"
	FontFile     = "/assets/fonts/Desyrel/desyrel.regular.ttf"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
	LEFT_MOUSE_BUTTON    = glfw.MouseButtonLeft
)

type Config struct {
	FormItemType interface{}
	Index        int
	Label        string
	Default      interface{}
}
type Conf map[string]*Config

func (c Conf) AddConfig(name string, itemType interface{}, def interface{}, label string) {
	c[name] = &Config{}
	c[name].Default = def
	c[name].FormItemType = itemType
	c[name].Label = label
}

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	glWrapper      glwrapper.Wrapper
	lastUpdate     int64
	startTime      int64
	Settings       = make(Conf)

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby

	DirectionalLightDirection = mgl32.Vec3{0.5, 0.5, -0.7}.Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightDiffuse   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightSpecular  = mgl32.Vec3{0.3, 0.3, 0.3}
)

func InitSettings() {
	Settings.AddConfig("Width", &model.FormItemInt{}, "4", "Rows (i)")
	Settings.AddConfig("Length", &model.FormItemInt{}, "4", "Cols (i)")
	Settings.AddConfig("Iterations", &model.FormItemInt{}, "10", "Iter (i)")
	Settings.AddConfig("Scale", &model.FormItemFloat{}, "5.0", "Scale (f)")
	Settings.AddConfig("PeakProb", &model.FormItemInt{}, "5", "Peak (i)")
	Settings.AddConfig("CliffProb", &model.FormItemInt{}, "5", "Cliff (i)")
	Settings.AddConfig("TerrainMinHeight", &model.FormItemFloat{}, "-1.0", "MinH (f)")
	Settings.AddConfig("TerrainMaxHeight", &model.FormItemFloat{}, "3.0", "MaxH (f)")
	Settings.AddConfig("TerrainPosY", &model.FormItemFloat{}, "1.003", "PosY (f)")
	Settings.AddConfig("Seed", &model.FormItemInt64{}, "0", "Seed (i64)")
	Settings.AddConfig("RandomSeed", &model.FormItemBool{}, false, "Rand Seed")
	Settings.AddConfig("TerrainTexture", &model.FormItemText{}, "Grass", "Terr tex")
	Settings.AddConfig("NeedLiquid", &model.FormItemBool{}, true, "Has Liquid")
	Settings.AddConfig("LiquidEta", &model.FormItemFloat{}, "0.75", "Leta (f)")
	Settings.AddConfig("LiquidAmplitude", &model.FormItemFloat{}, "0.0625", "Lampl (f)")
	Settings.AddConfig("LiquidFrequency", &model.FormItemFloat{}, "1.0", "LFreq (f)")
	Settings.AddConfig("LiquidDetail", &model.FormItemInt{}, "10", "Ldetail (i)")
	Settings.AddConfig("WaterLevel", &model.FormItemFloat{}, "0.25", "W Lev (f)")
	Settings.AddConfig("LiquidTexture", &model.FormItemText{}, "Water", "Liq tex")
	Settings.AddConfig("Debug", &model.FormItemBool{}, false, "Debug mode")
	Settings.AddConfig("ClearR", &model.FormItemFloat{}, "1.0", "BG R (f)")
	Settings.AddConfig("ClearG", &model.FormItemFloat{}, "0.0", "BG G (f)")
	Settings.AddConfig("ClearB", &model.FormItemFloat{}, "0.0", "BG B (f)")
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

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, -0.5, 3.0}, mgl32.Vec3{0, 1, 0}, -85.0, -0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 20.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func createMenu() *screen.MenuScreen {
	MenuFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	MenuFonts.SetTransparent(true)
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)
	return screen.NewMenuScreen(tex, DefaultMaterial, HighlightMaterial, MenuFonts, mgl32.Vec3{1, 0, 0}, mgl32.Vec3{1, 0, 0}, glWrapper)
}
func createSettings(defaults Conf) *screen.FormScreen {
	form := screen.NewFormScreen(material.Ruby, "Settings", glWrapper, float32(WindowWidth), float32(WindowHeight))
	if _, ok := defaults["Width"]; ok {
		index := form.AddFormItemInt(defaults["Width"].Label, glWrapper, defaults["Width"].Default.(string))
		defaults["Width"].Index = index
	}
	if _, ok := defaults["Length"]; ok {
		index := form.AddFormItemInt(defaults["Length"].Label, glWrapper, defaults["Length"].Default.(string))
		defaults["Length"].Index = index
	}
	if _, ok := defaults["Iterations"]; ok {
		index := form.AddFormItemInt(defaults["Iterations"].Label, glWrapper, defaults["Iterations"].Default.(string))
		defaults["Iterations"].Index = index
	}
	if _, ok := defaults["Scale"]; ok {
		index := form.AddFormItemFloat(defaults["Scale"].Label, glWrapper, defaults["Scale"].Default.(string))
		defaults["Scale"].Index = index
	}
	if _, ok := defaults["PeakProb"]; ok {
		index := form.AddFormItemInt(defaults["PeakProb"].Label, glWrapper, defaults["PeakProb"].Default.(string))
		defaults["PeakProb"].Index = index
	}
	if _, ok := defaults["CliffProb"]; ok {
		index := form.AddFormItemInt(defaults["CliffProb"].Label, glWrapper, defaults["CliffProb"].Default.(string))
		defaults["CliffProb"].Index = index
	}
	if _, ok := defaults["TerrainMinHeight"]; ok {
		index := form.AddFormItemFloat(defaults["TerrainMinHeight"].Label, glWrapper, defaults["TerrainMinHeight"].Default.(string))
		defaults["TerrainMinHeight"].Index = index
	}
	if _, ok := defaults["TerrainMaxHeight"]; ok {
		index := form.AddFormItemFloat(defaults["TerrainMaxHeight"].Label, glWrapper, defaults["TerrainMaxHeight"].Default.(string))
		defaults["TerrainMaxHeight"].Index = index
	}
	if _, ok := defaults["TerrainPosY"]; ok {
		index := form.AddFormItemFloat(defaults["TerrainPosY"].Label, glWrapper, defaults["TerrainPosY"].Default.(string))
		defaults["TerrainPosY"].Index = index
	}
	if _, ok := defaults["Seed"]; ok {
		index := form.AddFormItemInt64(defaults["Seed"].Label, glWrapper, defaults["Seed"].Default.(string))
		defaults["Seed"].Index = index
	}
	if _, ok := defaults["RandomSeed"]; ok {
		index := form.AddFormItemBool(defaults["RandomSeed"].Label, glWrapper, defaults["RandomSeed"].Default.(bool))
		defaults["RandomSeed"].Index = index
	}
	if _, ok := defaults["TerrainTexture"]; ok {
		index := form.AddFormItemText(defaults["TerrainTexture"].Label, glWrapper, defaults["TerrainTexture"].Default.(string))
		defaults["TerrainTexture"].Index = index
	}
	if _, ok := defaults["NeedLiquid"]; ok {
		index := form.AddFormItemBool(defaults["NeedLiquid"].Label, glWrapper, defaults["NeedLiquid"].Default.(bool))
		defaults["NeedLiquid"].Index = index
	}
	if _, ok := defaults["LiquidEta"]; ok {
		index := form.AddFormItemFloat(defaults["LiquidEta"].Label, glWrapper, defaults["LiquidEta"].Default.(string))
		defaults["LiquidEta"].Index = index
	}
	if _, ok := defaults["LiquidAmplitude"]; ok {
		index := form.AddFormItemFloat(defaults["LiquidAmplitude"].Label, glWrapper, defaults["LiquidAmplitude"].Default.(string))
		defaults["LiquidAmplitude"].Index = index
	}
	if _, ok := defaults["LiquidFrequency"]; ok {
		index := form.AddFormItemFloat(defaults["LiquidFrequency"].Label, glWrapper, defaults["LiquidFrequency"].Default.(string))
		defaults["LiquidFrequency"].Index = index
	}
	if _, ok := defaults["LiquidDetail"]; ok {
		index := form.AddFormItemInt(defaults["LiquidDetail"].Label, glWrapper, defaults["LiquidDetail"].Default.(string))
		defaults["LiquidDetail"].Index = index
	}
	if _, ok := defaults["WaterLevel"]; ok {
		index := form.AddFormItemFloat(defaults["WaterLevel"].Label, glWrapper, defaults["WaterLevel"].Default.(string))
		defaults["WaterLevel"].Index = index
	}
	if _, ok := defaults["LiquidTexture"]; ok {
		index := form.AddFormItemText(defaults["LiquidTexture"].Label, glWrapper, defaults["LiquidTexture"].Default.(string))
		defaults["LiquidTexture"].Index = index
	}
	if _, ok := defaults["Debug"]; ok {
		index := form.AddFormItemBool(defaults["Debug"].Label, glWrapper, defaults["Debug"].Default.(bool))
		defaults["Debug"].Index = index
	}
	if _, ok := defaults["ClearR"]; ok {
		index := form.AddFormItemFloat(defaults["ClearR"].Label, glWrapper, defaults["ClearR"].Default.(string))
		defaults["ClearR"].Index = index
	}
	if _, ok := defaults["ClearG"]; ok {
		index := form.AddFormItemFloat(defaults["ClearG"].Label, glWrapper, defaults["ClearG"].Default.(string))
		defaults["ClearG"].Index = index
	}
	if _, ok := defaults["ClearB"]; ok {
		index := form.AddFormItemFloat(defaults["ClearB"].Label, glWrapper, defaults["ClearB"].Default.(string))
		defaults["ClearB"].Index = index
	}
	return form
}
func createGame(preSets Conf, form *screen.FormScreen) *screen.Screen {
	AppScreen := screen.New()
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShaderBlendingWithFog(glWrapper)
	AppScreen.AddShader(shaderProgramTexture)
	// Shader application for the liquid surface
	shaderProgramLiquid := shader.NewTextureShaderLiquidWithFog(glWrapper)
	AppScreen.AddShader(shaderProgramLiquid)

	AppScreen.SetCamera(CreateCamera())
	AppScreen.SetCameraMovementMap(CameraMovementMap())
	AppScreen.SetRotateOnEdgeDistance(CameraDistance)

	gb := model.NewTerrainBuilder()
	// terrain related ones
	gb.SetWidth(form.GetFormItem(Settings["Width"].Index).(*model.FormItemInt).GetValue())
	gb.SetLength(form.GetFormItem(Settings["Length"].Index).(*model.FormItemInt).GetValue())
	gb.SetIterations(form.GetFormItem(Settings["Iterations"].Index).(*model.FormItemInt).GetValue())
	scaleValue := form.GetFormItem(Settings["Scale"].Index).(*model.FormItemFloat).GetValue()
	gb.SetScale(mgl32.Vec3{scaleValue, 1, scaleValue})
	gb.SetPeakProbability(form.GetFormItem(Settings["PeakProb"].Index).(*model.FormItemInt).GetValue())
	gb.SetCliffProbability(form.GetFormItem(Settings["CliffProb"].Index).(*model.FormItemInt).GetValue())
	gb.SetMinHeight(form.GetFormItem(Settings["TerrainMinHeight"].Index).(*model.FormItemFloat).GetValue())
	gb.SetMaxHeight(form.GetFormItem(Settings["TerrainMaxHeight"].Index).(*model.FormItemFloat).GetValue())
	posY := form.GetFormItem(Settings["TerrainPosY"].Index).(*model.FormItemFloat).GetValue()
	gb.SetPosition(mgl32.Vec3{0.0, posY, 0.0})
	if form.GetFormItem(Settings["RandomSeed"].Index).(*model.FormItemBool).GetValue() {
		seed := gb.RandomSeed()
		form.SetFormItemValue(Settings["Seed"].Index, strconv.FormatInt(seed, 10), glWrapper)
		form.SetFormItemValue(Settings["RandomSeed"].Index, false, glWrapper)
	} else {
		gb.SetSeed(form.GetFormItem(Settings["Seed"].Index).(*model.FormItemInt64).GetValue())
	}
	gb.SetGlWrapper(glWrapper)
	switch form.GetFormItem(Settings["TerrainTexture"].Index).(*model.FormItemText).GetValue() {
	case "Grass", "grass", "GRASS":
		gb.SurfaceTextureGrass()
		break
	default:
		gb.SurfaceTextureGrass()
		break
	}
	if form.GetFormItem(Settings["Debug"].Index).(*model.FormItemBool).GetValue() {
		gb.SetDebugMode(false)
	}
	if form.GetFormItem(Settings["NeedLiquid"].Index).(*model.FormItemBool).GetValue() {
		gb.SetLiquidEta(form.GetFormItem(Settings["LiquidEta"].Index).(*model.FormItemFloat).GetValue())
		gb.SetLiquidAmplitude(form.GetFormItem(Settings["LiquidAmplitude"].Index).(*model.FormItemFloat).GetValue())
		gb.SetLiquidFrequency(form.GetFormItem(Settings["LiquidFrequency"].Index).(*model.FormItemFloat).GetValue())
		gb.SetLiquidDetailMultiplier(form.GetFormItem(Settings["LiquidDetail"].Index).(*model.FormItemInt).GetValue())
		gb.SetLiquidWaterLevel(form.GetFormItem(Settings["WaterLevel"].Index).(*model.FormItemFloat).GetValue())
		switch form.GetFormItem(Settings["LiquidTexture"].Index).(*model.FormItemText).GetValue() {
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
	r := SettingsScreen.GetFormItem(Settings["ClearR"].Index).(*model.FormItemFloat).GetValue()
	g := SettingsScreen.GetFormItem(Settings["ClearG"].Index).(*model.FormItemFloat).GetValue()
	b := SettingsScreen.GetFormItem(Settings["ClearB"].Index).(*model.FormItemFloat).GetValue()
	glWrapper.ClearColor(r, g, b, 1.0)
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
	menuScreen := createMenu()
	SettingsScreen = createSettings(Settings)
	appScreen := createGame(Settings, SettingsScreen)

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
		appScreen.SetCamera(CreateCamera())
		appScreen = createGame(Settings, SettingsScreen)
		app.ActivateScreen(appScreen)
	}
	startEvent := func() {
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		appScreen = createGame(Settings, SettingsScreen)
		app.ActivateScreen(appScreen)
		menuScreen.SetState("world-started", true)
		menuScreen.BuildScreen(glWrapper, 3/float32(WindowWidth))
	}
	settingsEvent := func() {
		app.ActivateScreen(SettingsScreen)
	}
	continueEvent := func() {
		app.ActivateScreen(appScreen)
	}
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	cont := screen.NewMenuScreenOption("continue", contS, continueEvent)
	menuScreen.AddOption(*cont) // continue
	start := screen.NewMenuScreenOption("start", contNS, startEvent)
	menuScreen.AddOption(*start) // start
	restart := screen.NewMenuScreenOption("restart", contS, restartEvent)
	menuScreen.AddOption(*restart) // restart
	settings := screen.NewMenuScreenOption("settings", contAll, settingsEvent)
	menuScreen.AddOption(*settings) // settings
	exit := screen.NewMenuScreenOption("exit", contAll, exitEvent)
	menuScreen.AddOption(*exit) // exit
	menuScreen.BuildScreen(glWrapper, 3/float32(WindowWidth))

	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.SetWrapper(glWrapper)
	app.GetWindow().SetCharCallback(app.CharCallback)
	app.AddScreen(appScreen)
	app.AddScreen(menuScreen)
	app.AddScreen(SettingsScreen)
	app.MenuScreen(menuScreen)
	app.ActivateScreen(menuScreen)

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
