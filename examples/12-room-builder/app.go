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
	WindowTitle  = "Example - room builder tool."

	TOGGLE_DOOR_BUTTON = glfw.KeyC
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	glWrapper glwrapper.Wrapper

	lastUpdate int64
	startTime  int64

	DirectionalLightDirection = (mgl32.Vec3{0.5, 0.5, -0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightDiffuse   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightSpecular  = mgl32.Vec3{0.3, 0.3, 0.3}
)

func init() {
	runtime.LockOSThread()

	Settings.AddConfig("RoomPosition", "Position", "The center point of the floor.", mgl32.Vec3{0, 0, 0}, nil)
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
		"RoomPosition",
		"RoomWidth", "RoomLength",
		"RoomHeight", "RoomWallWidth",
		"RoomDoorWidth", "RoomDoorHeight",
		"RoomWindowWidth", "RoomWindowHeight",
		"RoomXDeg", "RoomYDeg",
		"RoomZDeg", "RoomTextured",
		"RoomFrontWindow", "RoomDoorOpened",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
		"CameraFPS",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Room editor")
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
	mdl, _, _ := app.GetClosestModelMeshDistance()
	switch mdl.(type) {
	case *model.Room:
		if app.GetKeyState(TOGGLE_DOOR_BUTTON) {
			room := mdl.(*model.Room)
			room.PushDoorState()
		}
		break
	}
}
func CreateGround() *model.Terrain {
	gb := model.NewTerrainBuilder()
	gb.SetWidth(4)
	gb.SetLength(4)
	gb.SetIterations(10)
	gb.SetScale(mgl32.Vec3{5, 1, 5})
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
func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCameraFromSettings(), CameraMovementOptions())
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShader(glWrapper)
	scrn.AddShader(shaderProgramTexture)

	scrn.AddModelToShader(CreateGround(), shaderProgramTexture)
	if Settings["RoomTextured"].GetCurrentValue().(bool) {
		// Shader application for the texture + blending
		shaderProgramRoom := shader.NewTextureShaderBlending(glWrapper)
		scrn.AddShader(shaderProgramRoom)
		scrn.AddModelToShader(GenerateRoom(), shaderProgramRoom)
	} else {
		// Shader application for the material objects
		shaderProgramRoom := shader.NewMaterialShader(glWrapper)
		scrn.AddShader(shaderProgramRoom)
		scrn.AddModelToShader(GenerateRoom(), shaderProgramRoom)
	}

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
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

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["forward"] = []glfw.Key{glfw.KeyW}
	cm["back"] = []glfw.Key{glfw.KeyS}
	cm["up"] = []glfw.Key{glfw.KeyQ}
	cm["down"] = []glfw.Key{glfw.KeyE}
	cm["left"] = []glfw.Key{glfw.KeyA}
	cm["right"] = []glfw.Key{glfw.KeyD}
	cm["rotateOnEdgeDistance"] = Settings["CameraRotationEdge"].GetCurrentValue().(float32)
	cm["mode"] = "default"
	return cm
}

func main() {
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
