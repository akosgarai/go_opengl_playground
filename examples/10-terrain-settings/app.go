package main

import (
	"fmt"
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
	FormItemType string
	Index        int
	Label        string
	Description  string
	Default      interface{}
}
type Conf map[string]*Config

func (c Conf) AddConfig(name, description, label, itemType string, def interface{}) {
	c[name] = &Config{}
	c[name].Default = def
	c[name].FormItemType = itemType
	c[name].Label = label
	c[name].Description = description
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
	Settings.AddConfig("Width", "The width / rows of the heightmap of the terrain.", "Rows (i)", "int", "4")
	Settings.AddConfig("Length", "The length / columns of the heightmap of the terrain.", "Cols (i)", "int", "4")
	Settings.AddConfig("Iterations", "The number of the iterations of the random height generation step.", "Iter (i)", "int", "10")
	Settings.AddConfig("PeakProb", "The probability of a given map position will be a peak.", "Peak (i)", "int", "5")
	Settings.AddConfig("CliffProb", "The probability of a given map position will be a cliff.", "Cliff (i)", "int", "5")
	Settings.AddConfig("TerrainMinHeight", "The minimum height of the terrain surface.", "MinH (f)", "float", "-1.0")
	Settings.AddConfig("TerrainMaxHeight", "The maximum height of the terrain surface.", "MaxH (f)", "float", "3.0")
	Settings.AddConfig("TerrainScale", "The generated terrain is scaled with the components of this vector.", "Scale (f)", "vector", [3]string{"5.0", "1.0", "5.0"})
	Settings.AddConfig("TerrainPos", "The center / middle point of the terrain mesh.", "PosY (f)", "vector", [3]string{"0.0", "1.003", "0.0"})
	Settings.AddConfig("Seed", "This value is used as seed for the random number generation, if the random is not set.", "Seed (i64)", "int64", "0")
	Settings.AddConfig("RandomSeed", "If this is set, the seed will be based on the current timestamp.", "Rand Seed", "bool", false)
	Settings.AddConfig("TerrainTexture", "The texture of the terrain. Currently the 'Grass' is supported.", "Terr tex", "text", "Grass")
	Settings.AddConfig("NeedLiquid", "If this is set, water will also be generated to the terrain.", "Has Liquid", "bool", true)
	Settings.AddConfig("LiquidEta", "The refraction ratio between the air and the liquid surface.", "Leta (f)", "float", "0.75")
	Settings.AddConfig("LiquidAmplitude", "The amplitude of the waves (sin wave) in the liquid surface", "Lampl (f)", "float", "0.0625")
	Settings.AddConfig("LiquidFrequency", "The wavelengts of the waves of the liquid surface.", "LFreq (f)", "float", "1.0")
	Settings.AddConfig("LiquidDetail", "The size of the liquid surface is the same as the terrain, but its height map is bigger this times.", "Ldetail (i)", "int", "10")
	Settings.AddConfig("WaterLevel", "The water level of the liquid surface.", "W Lev (f)", "float", "0.25")
	Settings.AddConfig("LiquidTexture", "The texture of the liquid surface. Currently the 'Water' is supported.", "Liq tex", "text", "Water")
	Settings.AddConfig("Debug", "Turn debug mode on - off. Currently it does nothing.", "Debug mode", "bool", false)
	Settings.AddConfig("ClearCol", "The clear color of the window. It is used as the color of the sky.", "BG color", "vector", [3]string{"0.2", "0.3", "0.8"})
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
	}
	for i := 0; i < len(formItemOrders); i++ {
		itemName := formItemOrders[i]
		if _, ok := defaults[itemName]; ok {
			switch defaults[itemName].FormItemType {
			case "text":
				index := form.AddFormItemText(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.(string), nil)
				defaults[itemName].Index = index
				break
			case "int":
				index := form.AddFormItemInt(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.(string), nil)
				defaults[itemName].Index = index
				break
			case "int64":
				index := form.AddFormItemInt64(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.(string), nil)
				defaults[itemName].Index = index
				break
			case "bool":
				index := form.AddFormItemBool(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.(bool))
				defaults[itemName].Index = index
				break
			case "float":
				index := form.AddFormItemFloat(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.(string), nil)
				defaults[itemName].Index = index
				break
			case "vector":
				index := form.AddFormItemVector(defaults[itemName].Label, defaults[itemName].Description, glWrapper, defaults[itemName].Default.([3]string), nil)
				defaults[itemName].Index = index
				break
			default:
				fmt.Printf("Unhandled form item type '%s'. (%#v)\n", defaults[itemName].FormItemType, defaults[itemName])
				break
			}
		}
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
	gb.SetScale(form.GetFormItem(Settings["TerrainScale"].Index).(*model.FormItemVector).GetValue())
	gb.SetPeakProbability(form.GetFormItem(Settings["PeakProb"].Index).(*model.FormItemInt).GetValue())
	gb.SetCliffProbability(form.GetFormItem(Settings["CliffProb"].Index).(*model.FormItemInt).GetValue())
	gb.SetMinHeight(form.GetFormItem(Settings["TerrainMinHeight"].Index).(*model.FormItemFloat).GetValue())
	gb.SetMaxHeight(form.GetFormItem(Settings["TerrainMaxHeight"].Index).(*model.FormItemFloat).GetValue())
	gb.SetPosition(form.GetFormItem(Settings["TerrainPos"].Index).(*model.FormItemVector).GetValue())
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
	col := SettingsScreen.GetFormItem(Settings["ClearCol"].Index).(*model.FormItemVector).GetValue()
	glWrapper.ClearColor(col.X(), col.Y(), col.Z(), 1.0)
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
