package main

import (
	"fmt"
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
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowTitle     = "Example - Cubemap menu screen"
	CAMERA_POSITION = glfw.KeyP
)

var (
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

	lastUpdate int64
	startTime  int64

	glWrapper glwrapper.Wrapper
	// config
	Settings = config.New()
	// This flag is active until the screen animation is running.
	// The animation is the movement to the computers.
	AnimationIsrunning = false
	// Camera log
	CameraLogIsPrinted = false
	// true if the current screen is the app screen.
	AppScreenIsActive = false
	// middle monitor & screen position
	MiddleMonitorPosition = mgl32.Vec3{2.0, 0, 0}
	MiddleScreenPosition  = mgl32.Vec3{-0.02, 0.0, 0.07}
)

func init() {
	// lock thread
	runtime.LockOSThread()
	// init the WindowBuilder
	setupWindowBuilder()
	// Camera configuration with initial values
	addCameraConfigToSettings()
	// Directional light configuration with initial values
	addDirectionalLightConfigToSettings()
	// Application config - camera start & stop position.
	addAnimationConfigToSettings()
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
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
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0, 0, -1}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(0.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.001), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(10.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.5), nil)
}
func addDirectionalLightConfigToSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("DLAmbient", "DL ambient", "The ambient color component of the directional lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("DLDiffuse", "DL diffuse", "The diffuse color component of the directional lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("DLSpecular", "DL specular", "The specular color component of the directional lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("DLDirection", "DL direction", "The direction vector of the directional lightsource.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
}
func addAnimationConfigToSettings() {
	Settings.AddConfig("CameraInitPos", "Cam init position", "The initial position of the camera.", mgl32.Vec3{-3.5, 0.0, 1.25}, nil)
	Settings.AddConfig("CameraFiniPos", "Cam finish position", "The initial position of the camera.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
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
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = CreateApplicationScreen()
		app.ActivateScreen(AppScreen)
		AppScreenIsActive = true
		AnimationIsrunning = true
	}
	startEvent := func() {
		lastUpdate = time.Now().UnixNano()
		startTime = lastUpdate
		AppScreen = CreateApplicationScreen()
		app.ActivateScreen(AppScreen)
		AppScreenIsActive = true
		AnimationIsrunning = true
		MenuScreen.SetState("world-started", true)
		MenuScreen.BuildScreen()
	}
	settingsEvent := func() {
		app.ActivateScreen(SettingsScreen)
		AppScreenIsActive = false
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

// It creates a new fps camera with the necessary setup from settings screen
func CreateCameraFromSettings() interfaces.Camera {
	cameraPosition := Settings["CameraInitPos"].GetCurrentValue().(mgl32.Vec3)
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
	cm["mode"] = "default"
	cm["rotateOnEdgeDistance"] = float32(0.0)
	cm["forward"] = []glfw.Key{glfw.KeyW}
	cm["back"] = []glfw.Key{glfw.KeyS}
	cm["up"] = []glfw.Key{glfw.KeyQ}
	cm["down"] = []glfw.Key{glfw.KeyE}
	cm["left"] = []glfw.Key{glfw.KeyA}
	cm["right"] = []glfw.Key{glfw.KeyD}
	return cm
}

// Create a textured rectangle, that represents the monitors.
func CreateTexturedRectangle(t texture.Textures) *mesh.TexturedMesh {
	r := rectangle.NewExact(2, 2)
	V, I, _ := r.MeshInput()
	rect := mesh.NewTexturedMesh(V, I, t, glWrapper)
	return rect
}

// Create material cube, that represents the table and the screens.
func CreateMaterialCube(m *material.Material, size mgl32.Vec3) *mesh.MaterialMesh {
	r := cuboid.New(size.X(), size.Y(), size.Z())
	V, I, _ := r.MaterialMeshInput()
	rect := mesh.NewMaterialMesh(V, I, m, glWrapper)
	return rect
}

// Create a textured rectangle, that represents the walls.
func CreateTexturedCube(t texture.Textures) *mesh.TexturedMesh {
	r := cuboid.New(8, 8, 8)
	V, I, _ := r.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	return mesh.NewTexturedMesh(V, I, t, glWrapper)
}
func CreateApplicationScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCameraFromSettings(), CameraMovementOptions())
	monitors := model.New()
	desk := model.New()
	rustySurface := model.New()
	screens := model.New()

	var monitorTexture texture.Textures
	monitorTexture.AddTexture(baseDir()+"/assets/crt_monitor_1280.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	monitorTexture.AddTexture(baseDir()+"/assets/crt_monitor_1280.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", glWrapper)

	middleMonitor := CreateTexturedRectangle(monitorTexture)
	middleMonitor.SetPosition(MiddleMonitorPosition)
	middleMonitor.RotateZ(90)
	monitors.AddMesh(middleMonitor)
	middleMonitorScreen := CreateMaterialCube(material.Emerald, mgl32.Vec3{1.5, 1.3, 0})
	middleMonitorScreen.SetParent(middleMonitor)
	middleMonitorScreen.SetPosition(MiddleScreenPosition)
	screens.AddMesh(middleMonitorScreen)

	rightMonitor := CreateTexturedRectangle(monitorTexture)
	rightMonitorPosition := mgl32.TransformCoordinate(MiddleMonitorPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(-60)))
	rightMonitor.SetPosition(rightMonitorPosition)
	rightMonitor.RotateZ(45)
	monitors.AddMesh(rightMonitor)
	rightMonitorScreen := CreateMaterialCube(material.Emerald, mgl32.Vec3{1.5, 1.3, 0})
	rightMonitorScreen.SetParent(rightMonitor)
	rightMonitorScreen.SetPosition(mgl32.TransformCoordinate(MiddleScreenPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(-60))))
	screens.AddMesh(rightMonitorScreen)

	leftMonitor := CreateTexturedRectangle(monitorTexture)
	leftMonitorPosition := mgl32.TransformCoordinate(MiddleMonitorPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(60)))
	leftMonitor.SetPosition(leftMonitorPosition)
	leftMonitor.RotateZ(135)
	monitors.AddMesh(leftMonitor)
	leftMonitorScreen := CreateMaterialCube(material.Emerald, mgl32.Vec3{1.5, 1.3, 0})
	leftMonitorScreen.SetParent(leftMonitor)
	leftMonitorScreen.SetPosition(mgl32.TransformCoordinate(MiddleScreenPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(60))))
	screens.AddMesh(leftMonitorScreen)

	monitors.SetTransparent(true)
	shaderAppTextureBlending := shader.NewTextureShaderBlending(glWrapper)
	scrn.AddShader(shaderAppTextureBlending)
	scrn.AddModelToShader(monitors, shaderAppTextureBlending)

	tableSurfaceMaterial := material.Chrome
	tableSurface := CreateMaterialCube(tableSurfaceMaterial, mgl32.Vec3{2, 6, 0.05})
	tableSurface.SetPosition(mgl32.Vec3{1.5, 0, -1})
	tableSurface.RotateX(90)
	desk.AddMesh(tableSurface)

	shaderAppMaterial := shader.NewMaterialShader(glWrapper)
	scrn.AddShader(shaderAppMaterial)
	scrn.AddModelToShader(desk, shaderAppMaterial)
	scrn.AddModelToShader(screens, shaderAppMaterial)

	shaderAppTexture := shader.NewTextureShader(glWrapper)
	scrn.AddShader(shaderAppTexture)
	var rustyTexture texture.Textures
	rustyTexture.AddTexture(baseDir()+"/assets/rusty_iron_1280.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	rustyTexture.AddTexture(baseDir()+"/assets/rusty_iron_1280.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", glWrapper)
	rustySurface.AddMesh(CreateTexturedCube(rustyTexture))
	scrn.AddModelToShader(rustySurface, shaderAppTexture)

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
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
}

// It creates the Settings screen.
func CreateSettingsScreen(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{}
	return app.BuildFormScreen(defaults, formItemOrders, "Settings")
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
	if app.GetKeyState(CAMERA_POSITION) {
		if !CameraLogIsPrinted {
			fmt.Println(AppScreen.GetCamera().Log())
			CameraLogIsPrinted = true
		}
	} else {
		CameraLogIsPrinted = false
	}
	// Camera animation for the AppScreen:
	// It starts from the (-3.5,0,1.2) position and goes to (0,0,0)
	// First it goest to the +x direction until it reaches 0. Then it
	// changes the direction to -z and goes until it reaches 0.
	if AppScreenIsActive {
		if !AnimationIsrunning {
			return
		}
		camPos := AppScreen.GetCamera().GetPosition()
		finish := Settings["CameraFiniPos"].GetCurrentValue().(mgl32.Vec3)
		step := float32(delta) * AppScreen.GetCamera().GetVelocity()
		if camPos.X() < finish.X() {
			AppScreen.GetCamera().Walk(step)
		} else if camPos.Z() > finish.Z() {
			AppScreen.GetCamera().Lift(-step)
		} else {
			AnimationIsrunning = false
		}
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
	// application screen
	AppScreen = CreateApplicationScreen()
	app.AddScreen(AppScreen)
	// menu screen
	MenuScreen = CreateMenuScreen()
	app.AddScreen(MenuScreen)
	app.MenuScreen(MenuScreen)
	// settings screen
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
