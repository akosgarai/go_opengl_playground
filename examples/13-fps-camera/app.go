package main

import (
	"os"
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
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowTitle       = "Example - FPS camera application"
	LEFT_MOUSE_BUTTON = glfw.MouseButtonLeft
	Epsilon           = float64(200)
)

var (
	glWrapper      glwrapper.Wrapper
	app            *application.Application
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	SettingsScreen *screen.FormScreen
	// window related variables
	Builder          *window.WindowBuilder
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowFullScreen = false
	// config
	Settings = config.New()

	lastUpdate     int64
	startTime      int64
	LampOn         bool
	LampLastToggle float64
)

func setupWindowBuilder() {
	Builder = window.NewWindowBuilder()
	fullScreen := os.Getenv("FULL")
	if fullScreen == "1" {
		WindowFullScreen = true
		WindowDecorated = false
		WindowWidth, WindowHeight = Builder.GetCurrentMonitorResolution()
	} else {
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
	}
	Builder.SetFullScreen(WindowFullScreen)
	Builder.SetDecorated(WindowDecorated)
	Builder.SetTitle(WindowTitle)
	Builder.SetWindowSize(WindowWidth, WindowHeight)
}
func addCameraConfigToSettings() {
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0.0, -0.5, 3.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0, 1, 0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-90.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.001), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(20.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.5), nil)
}
func addTerrainConfigToSettings() {
	Settings.AddConfig("GroundWidth", "Ground width", "The value is used for generating map. The GroundBuilder will generate width*width tiles total.", int(10), nil)
	Settings.AddConfig("GroundScale", "Ground scale", "The value is used for generating map. The tile size will be scale*scale unit.", float32(2.0), nil)
}
func addRoomConfigToSettings() {
	Settings.AddConfig("RoomPosition", "Position", "The center point of the floor.", mgl32.Vec3{2, 0, 1}, nil)
	Settings.AddConfig("RoomWidth", "Width", "The width of the room. The size in the X axis", float32(1.0), nil)
	Settings.AddConfig("RoomLength", "Length", "The length of the room. The size in the Z axis.", float32(1.0), nil)
	Settings.AddConfig("RoomHeight", "Height", "The height of the room. The size in the Y axis.", float32(1.0), nil)
	Settings.AddConfig("RoomWallWidth", "Wall width", "The width of the walls.", float32(0.005), nil)
	Settings.AddConfig("RoomDoorWidth", "Door width", "The width of the door. The size in the X axis", float32(0.4), nil)
	Settings.AddConfig("RoomDoorHeight", "Door height", "The height of the door. The size in the Y axis.", float32(0.6), nil)
	Settings.AddConfig("RoomWindowWidth", "Window width", "The width of the window. The size in the X axis", float32(0.2), nil)
	Settings.AddConfig("RoomWindowHeight", "Window height", "The height of the window. The size in the Y axis.", float32(0.4), nil)
	// Room rotation degrees
	Settings.AddConfig("RoomXDeg", "Room X rot.", "The rotation on the X axis.", float32(0.0), nil)
	Settings.AddConfig("RoomYDeg", "Room Y rot.", "The rotation on the Y axis.", float32(0.0), nil)
	Settings.AddConfig("RoomZDeg", "Room Z rot.", "The rotation on the Z axis.", float32(180.0), nil)
	Settings.AddConfig("RoomTextured", "Textured", "If this flag is active, the room will be textured, otherwise material.", false, nil)
	Settings.AddConfig("RoomFrontWindow", "Front window", "If this flag is active, a window will be set to the front wall. Only works in textured mode.", false, nil)
	Settings.AddConfig("RoomDoorOpened", "Door opened", "If this flag is active, the door will start in opened state.", false, nil)
}
func addLampConfigToSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("LampPosition", "Position", "The center point of the lamp bottom surface.", mgl32.Vec3{0, -3, 0}, nil)
	Settings.AddConfig("PoleLength", "Pole length", "The length of the pole. The size in the Z axis.", float32(3.0), nil)
	Settings.AddConfig("LampBulbShininess", "Bulb mat sh.", "The shininess of the bulb material.", float32(36.0), nil)
	// spot light
	Settings.AddConfig("LampBulbAmbient", "Spot ambient", "The ambient color component of the bulb. Also the ambinet color of the spot light.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LampBulbDiffuse", "Spot diffuse", "The diffuse color component of the bulb. Also the diffuse color of the spot light.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LampBulbSpecular", "Spot specular", "The specular color component of the bulb. Also the specular color of the spot light.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	Settings.AddConfig("SLCutoff", "Spot cutoff", "The inner cutoff component of the spot light equasion.", float32(4.0), nil)
	Settings.AddConfig("SLOuterCutoff", "Spot o. cutoff", "The inner cutoff component of the spot light equasion.", float32(5.0), nil)
	// type of the lamp: material or textured
	Settings.AddConfig("LampTextured", "Textured", "If this flag is active, the lamp will be textured, otherwise material.", false, nil)
	// turned on by default.
	Settings.AddConfig("LampOn", "Light", "If this flag is active, the color of the spotlight will be the color of the bulb material. otherwise black.", true, nil)
	// rotations
	Settings.AddConfig("LampRotateXDeg", "Lamp X rot.", "The rotation on the X axis.", float32(90.0), nil)
	Settings.AddConfig("LampRotateYDeg", "Lamp Y rot.", "The rotation on the Y axis.", float32(0.0), nil)
	Settings.AddConfig("LampRotateZDeg", "Lamp Z rot.", "The rotation on the Z axis.", float32(0.0), nil)
}
func addDirectionalLightConfigToSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("DLAmbient", "DL ambient", "The ambient color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDiffuse", "DL diffuse", "The diffuse color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLSpecular", "DL specular", "The specular color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDirection", "DL direction", "The direction vector of the directional lightsource.", mgl32.Vec3{0.7, 0.7, 0.7}, nil)
}

func init() {
	// lock thread
	runtime.LockOSThread()
	// init the WindowBuilder
	setupWindowBuilder()
	// Camera configuration with initial values
	addCameraConfigToSettings()
	// Terrain configuration with initial values
	addTerrainConfigToSettings()
	// Room configuration with initial values
	addRoomConfigToSettings()
	// Lamp configuration with initial values
	addLampConfigToSettings()
	// Directional light configuration with initial values
	addDirectionalLightConfigToSettings()
}

// It creates a new fps camera with the necessary setup from settings screen
func CreateCameraFromSettings() interfaces.Camera {
	cameraPosition := Settings["CameraPos"].GetCurrentValue().(mgl32.Vec3)
	worldUp := Settings["WorldUp"].GetCurrentValue().(mgl32.Vec3)
	yawAngle := Settings["CameraYaw"].GetCurrentValue().(float32)
	pitchAngle := Settings["CameraPitch"].GetCurrentValue().(float32)
	fov := Settings["CameraFov"].GetCurrentValue().(float32)
	near := Settings["CameraNear"].GetCurrentValue().(float32)
	far := Settings["CameraFar"].GetCurrentValue().(float32)
	moveSpeed := Settings["CameraVelocity"].GetCurrentValue().(float32)
	directionSpeed := Settings["CameraRotation"].GetCurrentValue().(float32)
	cam := camera.NewFPSCamera(cameraPosition, worldUp, yawAngle, pitchAngle)
	cam.SetupProjection(fov, float32(WindowWidth)/float32(WindowHeight), near, far)
	cam.SetVelocity(moveSpeed)
	cam.SetRotationStep(directionSpeed)
	return cam
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
	cm["mode"] = "fps"
	return cm
}

// It creates the terrain model. Currently the width and scale is manageable.
func CreateGround() *model.Terrain {
	gb := model.NewTerrainBuilder()
	w := Settings["GroundWidth"].GetCurrentValue().(int)
	gb.SetWidth(w)
	gb.SetLength(w)
	gb.SetIterations(10)
	s := Settings["GroundScale"].GetCurrentValue().(float32)
	gb.SetScale(mgl32.Vec3{s, 1, s})
	gb.SetGlWrapper(glWrapper)
	gb.SurfaceTextureGrass()
	gb.SetPeakProbability(0)
	gb.SetCliffProbability(0)
	gb.SetMinHeight(0)
	gb.SetMaxHeight(0)
	gb.SetPosition(mgl32.Vec3{0.0, 0.0, 0.0})
	gb.SetSeed(0)
	return gb.Build()
}

// It creates the room model based on the Settings.
func GenerateRoom() *model.Room {
	builder := model.NewRoomBuilder()
	builder.SetPosition(Settings["RoomPosition"].GetCurrentValue().(mgl32.Vec3))
	w := Settings["RoomWidth"].GetCurrentValue().(float32)
	h := Settings["RoomHeight"].GetCurrentValue().(float32)
	l := Settings["RoomLength"].GetCurrentValue().(float32)
	dw := Settings["RoomDoorWidth"].GetCurrentValue().(float32)
	dh := Settings["RoomDoorHeight"].GetCurrentValue().(float32)
	ww := Settings["RoomWindowWidth"].GetCurrentValue().(float32)
	wh := Settings["RoomWindowHeight"].GetCurrentValue().(float32)
	rww := Settings["RoomWallWidth"].GetCurrentValue().(float32)
	rx := Settings["RoomXDeg"].GetCurrentValue().(float32)
	ry := Settings["RoomYDeg"].GetCurrentValue().(float32)
	rz := Settings["RoomZDeg"].GetCurrentValue().(float32)
	builder.SetWrapper(glWrapper)
	builder.SetWallWidth(rww)
	builder.SetDoorSize(dw, dh)
	builder.SetWindowSize(ww, wh)
	builder.SetSize(w, h, l)
	builder.SetRotation(rx, ry, rz)
	var room *model.Room
	if Settings["RoomDoorOpened"].GetCurrentValue().(bool) {
		builder.WithOpenedDoor()
	} else {
		builder.WithClosedDoor()
	}
	if Settings["RoomTextured"].GetCurrentValue().(bool) {
		if Settings["RoomFrontWindow"].GetCurrentValue().(bool) {
			builder.WithFrontWindow(true)
		}
		room = builder.BuildTexture()
	} else {
		room = builder.BuildMaterial()
	}
	return room
}

// It creates a StreetLamp model based on the Settings.
func GenerateStreetLamp() *model.StreetLamp {
	builder := model.NewStreetLampBuilder()
	builder.SetWrapper(glWrapper)
	builder.SetPosition(Settings["LampPosition"].GetCurrentValue().(mgl32.Vec3))
	c := Settings["LSConstantTerm"].GetCurrentValue().(float32)
	l := Settings["LSLinearTerm"].GetCurrentValue().(float32)
	q := Settings["LSQuadraticTerm"].GetCurrentValue().(float32)
	builder.SetLightTerms(c, l, q)
	cutoff := Settings["SLCutoff"].GetCurrentValue().(float32)
	outerCutoff := Settings["SLOuterCutoff"].GetCurrentValue().(float32)
	builder.SetCutoff(cutoff, outerCutoff)
	builder.SetPoleLength(Settings["PoleLength"].GetCurrentValue().(float32))
	rx := Settings["LampRotateXDeg"].GetCurrentValue().(float32)
	ry := Settings["LampRotateYDeg"].GetCurrentValue().(float32)
	rz := Settings["LampRotateZDeg"].GetCurrentValue().(float32)
	builder.SetRotation(rx, ry, rz)
	mat := material.New(Settings["LampBulbAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["LampBulbDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["LampBulbSpecular"].GetCurrentValue().(mgl32.Vec3),
		Settings["LampBulbShininess"].GetCurrentValue().(float32))
	builder.SetBulbMaterial(mat)
	LampOn = Settings["LampOn"].GetCurrentValue().(bool)
	builder.SetLampOn(LampOn)
	var lamp *model.StreetLamp
	if Settings["LampTextured"].GetCurrentValue().(bool) {
		lamp = builder.BuildTexture()
	} else {
		lamp = builder.BuildMaterial()
	}
	return lamp
}

// It creates the menu screen.
func CreateMenuScreen() *screen.MenuScreen {
	showAll := func(m map[string]bool) bool {
		return true
	}
	showIfStarted := func(m map[string]bool) bool {
		return m["world-started"]
	}
	showIfNotStarted := func(m map[string]bool) bool {
		return !m["world-started"]
	}
	restartEvent := func() {
		// Disable mouse cursor
		app.GetWindow().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		app.GetWindow().SetInputMode(glfw.RawMouseMotion, 1)

		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = CreateApplicationScreen()
		app.ActivateScreen(AppScreen)
	}
	startEvent := func() {
		// Disable mouse cursor
		app.GetWindow().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		app.GetWindow().SetInputMode(glfw.RawMouseMotion, 1)

		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = CreateApplicationScreen()
		app.ActivateScreen(AppScreen)
		MenuScreen.SetState("world-started", true)
		MenuScreen.BuildScreen()
	}
	settingsEvent := func() {
		// Disable mouse cursor
		app.GetWindow().SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		app.GetWindow().SetInputMode(glfw.RawMouseMotion, 0)
		app.ActivateScreen(SettingsScreen)
	}
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	continueEvent := func() {
		// Disable mouse cursor
		app.GetWindow().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		app.GetWindow().SetInputMode(glfw.RawMouseMotion, 1)
		app.ActivateScreen(AppScreen)
	}
	options := []screen.Option{
		*screen.NewMenuScreenOption("Continue", showIfStarted, continueEvent),
		*screen.NewMenuScreenOption("Start", showIfNotStarted, startEvent),
		*screen.NewMenuScreenOption("Restart", showIfStarted, restartEvent),
		*screen.NewMenuScreenOption("Settings", showAll, settingsEvent),
		*screen.NewMenuScreenOption("Exit", showAll, exitEvent),
	}
	return app.BuildMenuScreen(options)
}

// It creates the Settings screen.
func CreateSettingsScreen(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"GroundWidth", "GroundScale",

		"RoomPosition",
		"RoomWidth", "RoomLength",
		"RoomHeight", "RoomWallWidth",
		"RoomDoorWidth", "RoomDoorHeight",
		"RoomWindowWidth", "RoomWindowHeight",
		"RoomXDeg", "RoomYDeg",
		"RoomZDeg", "RoomTextured",
		"RoomFrontWindow", "RoomDoorOpened",

		"LampPosition",
		"PoleLength", "LampBulbShininess",
		"LampBulbAmbient",
		"LampBulbDiffuse",
		"LampBulbSpecular",
		"LSConstantTerm", "LSLinearTerm",
		"LSQuadraticTerm", "SLCutoff",
		"SLOuterCutoff", "LampRotateXDeg",
		"LampRotateYDeg", "LampRotateZDeg",
		"LampTextured", "LampOn",

		"DLAmbient",
		"DLDiffuse",
		"DLSpecular",
		"DLDirection",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "FPS editor")
}

func CreateApplicationScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCameraFromSettings(), CameraMovementOptions())
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShader(glWrapper)
	scrn.AddShader(shaderProgramTexture)

	scrn.AddModelToShader(CreateGround(), shaderProgramTexture)
	var shaderProgramRoom *shader.Shader
	if Settings["RoomTextured"].GetCurrentValue().(bool) {
		// Shader application for the texture + blending
		shaderProgramRoom = shader.NewTextureShaderBlending(glWrapper)
	} else {
		// Shader application for the material objects
		shaderProgramRoom = shader.NewMaterialShader(glWrapper)
	}
	scrn.AddShader(shaderProgramRoom)
	scrn.AddModelToShader(GenerateRoom(), shaderProgramRoom)

	var shaderProgramLamp *shader.Shader
	if Settings["LampTextured"].GetCurrentValue().(bool) {
		// Shader application for the texture
		shaderProgramLamp = shader.NewTextureShader(glWrapper)
	} else {
		// Shader application for the material objects
		shaderProgramLamp = shader.NewMaterialShader(glWrapper)
	}
	scrn.AddShader(shaderProgramLamp)
	lamp := GenerateStreetLamp()
	scrn.AddSpotLightSource(lamp.GetLightSource(), [10]string{
		"spotLight[0].position", "spotLight[0].direction",
		"spotLight[0].ambient", "spotLight[0].diffuse",
		"spotLight[0].specular", "spotLight[0].constant",
		"spotLight[0].linear", "spotLight[0].quadratic",
		"spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	scrn.AddModelToShader(lamp, shaderProgramLamp)

	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		Settings["DLDirection"].GetCurrentValue().(mgl32.Vec3),
		Settings["DLAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["DLDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["DLSpecular"].GetCurrentValue().(mgl32.Vec3),
	})
	// Add the lightources to the application
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.Setup(setupApp)
	return scrn
}

func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
	LampLastToggle += delta
	mdl, _, _ := app.GetClosestModelMeshDistance()
	switch mdl.(type) {
	case *model.Room:
		if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
			room := mdl.(*model.Room)
			room.PushDoorState()
		}
		break
	case *model.StreetLamp:
		if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) && LampLastToggle > Epsilon {
			LampLastToggle = 0
			LampOn = !LampOn
			light := mdl.(*model.StreetLamp)
			if LampOn {
				light.TurnLampOn()
			} else {
				light.TurnLampOff()
			}

		}
		break
	}
}
func main() {
	app = application.New(glWrapper)
	// Setup the window
	app.SetWindow(Builder.Build())
	// Terminate window at the end.
	defer glfw.Terminate()
	// Init opengl.
	glWrapper.InitOpenGL()

	app.AddScreen(CreateApplicationScreen())

	MenuScreen = CreateMenuScreen()
	app.AddScreen(MenuScreen)
	app.MenuScreen(MenuScreen)
	SettingsScreen = CreateSettingsScreen(Settings)
	app.AddScreen(SettingsScreen)
	app.ActivateScreen(MenuScreen)
	// Register callbacks
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetCharCallback(app.CharCallback)

	// main event loop
	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		Update()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
