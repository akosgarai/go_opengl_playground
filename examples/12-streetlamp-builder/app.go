package main

import (
	"runtime"
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
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - streetlamp builder tool."

	TOGGLE_LIGHT_BUTTON = glfw.KeyC
	Epsilon             = float64(200)
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()
	SpotLight      light.Light

	glWrapper glwrapper.Wrapper

	lastUpdate int64
	startTime  int64
	LampOn     bool
	LastToggle float64
)

func InitRoomSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("GroundWidth", "Ground width", "The value is used for generating map. The GroundBuilder will generate width*width tiles total.", int(10), nil)
	Settings.AddConfig("GroundScale", "Ground scale", "The value is used for generating map. The tile size will be scale*scale unit.", float32(2.0), nil)
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
	// directional light
	Settings.AddConfig("DLAmbient", "DL ambient", "The ambient color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDiffuse", "DL diffuse", "The diffuse color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLSpecular", "DL specular", "The specular color component of the directional lightsource.", mgl32.Vec3{0.1, 0.1, 0.1}, colorValidator)
	Settings.AddConfig("DLDirection", "DL direction", "The direction vector of the directional lightsource.", mgl32.Vec3{0.7, 0.7, 0.7}, nil)
	// rotations
	Settings.AddConfig("RotateXDeg", "Rotate X", "The rotation on the X axis.", float32(90.0), nil)
	Settings.AddConfig("RotateYDeg", "Rotate Y", "The rotation on the Y axis.", float32(0.0), nil)
	Settings.AddConfig("RotateZDeg", "Rotate Z", "The rotation on the Z axis.", float32(0.0), nil)
	// type of the lamp: material or textured
	Settings.AddConfig("LampTextured", "Textured", "If this flag is active, the lamp will be textured, otherwise material.", false, nil)
	// turned on by default.
	Settings.AddConfig("LampOn", "Light", "If this flag is active, the color of the spotlight will be the color of the bulb material. otherwise black.", true, nil)

	// camera options:
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
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.05), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
	// - FPS camera
	Settings.AddConfig("CameraFPS", "FPS Camera", "If this flag is true, the camera will be FPS like.", false, nil)
}

func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"GroundWidth", "GroundScale",
		"LampPosition",
		"PoleLength", "LampBulbShininess",
		"LampBulbAmbient",
		"LampBulbDiffuse",
		"LampBulbSpecular",
		"DLAmbient",
		"DLDiffuse",
		"DLSpecular",
		"DLDirection",
		"LSConstantTerm", "LSLinearTerm",
		"LSQuadraticTerm", "SLCutoff",
		"SLOuterCutoff", "RotateXDeg",
		"RotateYDeg", "RotateZDeg",
		"LampTextured", "LampOn",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
		"CameraFPS",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Spotlight editor")
}

func createMenu() *screen.MenuScreen {
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
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	continueEvent := func() {
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
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
	LastToggle += delta
	mdl, _, _ := app.GetClosestModelMeshDistance()
	switch mdl.(type) {
	case *model.StreetLamp:
		if app.GetKeyState(TOGGLE_LIGHT_BUTTON) && LastToggle > Epsilon {
			LastToggle = 0
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
	Ground := gb.Build()
	return Ground
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
}
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
	rx := Settings["RotateXDeg"].GetCurrentValue().(float32)
	ry := Settings["RotateYDeg"].GetCurrentValue().(float32)
	rz := Settings["RotateZDeg"].GetCurrentValue().(float32)
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
func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetCamera(CreateCameraFromSettings())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(Settings["CameraRotationEdge"].GetCurrentValue().(float32))
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShader(glWrapper)
	scrn.AddShader(shaderProgramTexture)

	scrn.AddModelToShader(CreateGround(), shaderProgramTexture)
	var shaderProgramRoom *shader.Shader
	if Settings["LampTextured"].GetCurrentValue().(bool) {
		// Shader application for the texture
		shaderProgramRoom = shader.NewTextureShader(glWrapper)
	} else {
		// Shader application for the material objects
		shaderProgramRoom = shader.NewMaterialShader(glWrapper)
	}
	scrn.AddShader(shaderProgramRoom)
	lamp := GenerateStreetLamp()
	scrn.AddSpotLightSource(lamp.GetLightSource(), [10]string{
		"spotLight[0].position", "spotLight[0].direction",
		"spotLight[0].ambient", "spotLight[0].diffuse",
		"spotLight[0].specular", "spotLight[0].constant",
		"spotLight[0].linear", "spotLight[0].quadratic",
		"spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	scrn.AddModelToShader(lamp, shaderProgramRoom)

	// directional light is coming from the up direction but not from too up.
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

// It creates a new camera with the necessary setup from settings screen
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
	if Settings["CameraFPS"].GetCurrentValue().(bool) {
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

func main() {
	runtime.LockOSThread()
	InitRoomSettings()
	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := mainScreen()
	app.AddScreen(scrn)
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetCharCallback(app.CharCallback)
	MenuScreen = createMenu()
	app.AddScreen(MenuScreen)
	app.MenuScreen(MenuScreen)
	SettingsScreen = createSettings(Settings)
	app.AddScreen(SettingsScreen)
	app.ActivateScreen(MenuScreen)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		Update()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
