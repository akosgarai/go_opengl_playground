package main

import (
	"os"
	"path"
	"runtime"
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
	"github.com/akosgarai/playground_engine/pkg/primitives/cylinder"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - material light - with rotation"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app *application.Application

	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	lastUpdate int64
	startTime  int64

	PointLightSource *light.Light
	LightSourceCube  *mesh.MaterialMesh
	JadeCube         *mesh.MaterialMesh

	glWrapper glwrapper.Wrapper
)

func init() {
	runtime.LockOSThread()

	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	// light source
	Settings.AddConfig("LSAmbient", "Light ambient", "The ambient color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSDiffuse", "Light diffuse", "The diffuse color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSSpecular", "Light specular", "The specular color component of the lightsource.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("LSPosition", "Light position", "The position vector of the lightsource.", mgl32.Vec3{-3.0, 0.0, -3.0}, nil)
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	// cubes
	Settings.AddConfig("JadePosition", "Jade cube pos", "The position vector of the jade cube.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("JadeScale", "Jade cube scale", "The scale vector of the jade cube.", mgl32.Vec3{1.0, 1.0, 1.0}, nil)
	Settings.AddConfig("RPPosition", "RP cube pos", "The position vector of the red plastic cube.", mgl32.Vec3{-6.5, -3.5, -4.5}, nil)
	Settings.AddConfig("RPScale", "RP cube scale", "The scale vector of the red plastic cube.", mgl32.Vec3{2.0, 2.0, 2.0}, nil)
	Settings.AddConfig("ObsidianPosition", "Obs cube pos", "The position vector of the obsidian cube.", mgl32.Vec3{-7.5, -4.5, -0.5}, nil)
	Settings.AddConfig("ObsidianScale", "Obs cube scale", "The scale vector of the obsidian cube.", mgl32.Vec3{1.0, 1.0, 1.0}, nil)
	Settings.AddConfig("CopperPosition", "Cop cube pos", "The position vector of the copper cube.", mgl32.Vec3{2.0, -4.5, -0.5}, nil)
	Settings.AddConfig("CopperScale", "Cop cube scale", "The scale vector of the copper cube.", mgl32.Vec3{1.0, 1.0, 1.0}, nil)
	Settings.AddConfig("SilverPosition", "Sil cube pos", "The position vector of the silver cube.", mgl32.Vec3{2.0, -2.5, -1.5}, nil)
	Settings.AddConfig("SilverScale", "Sil cube scale", "The scale vector of the silver cube.", mgl32.Vec3{1.0, 1.0, 1.0}, nil)
	Settings.AddConfig("TurquoisePosition", "Tur cyl pos", "The position vector of the turquoise cylinder.", mgl32.Vec3{4, 3, -3}, nil)
	Settings.AddConfig("TurquoiseScale", "Tur cyl scale", "The scale vector of the turquoise cylinder.", mgl32.Vec3{1.0, 1.0, 1.0}, nil)
	Settings.AddConfig("LSPosition", "Ls cube pos", "The position vector of the lightsource cube.", mgl32.Vec3{-3.0, -1.5, -3.0}, nil)
	Settings.AddConfig("LSRoundSpeed", "ls round sp", "The round speed of the lightsource. It goes one round in this interval.", float32(3000.0), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{3.3, -10, 14.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, -1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(-101.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(21.5), nil)
	// - fov, far, near clip
	Settings.AddConfig("CameraNear", "Cam Near", "The near clip plane of the camera.", float32(0.001), nil)
	Settings.AddConfig("CameraFar", "Cam Far", "The far clip plane of the camera.", float32(50.0), nil)
	Settings.AddConfig("CameraFov", "Cam Fov", "The field of view (angle) of the camera.", float32(45.0), nil)
	// - move speed
	Settings.AddConfig("CameraVelocity", "Cam Speed", "The movement velocity of the camera. If it moves, it moves with this speed.", float32(0.005), nil)
	// - direction speed
	Settings.AddConfig("CameraRotation", "Cam Rotate", "The rotation velocity of the camera. If it rotates, it rotates with this speed.", float32(0.005), nil)
	// - rotate on edge distance.
	Settings.AddConfig("CameraRotationEdge", "Cam Edge", "The rotation cam be triggered if the mouse is near to the edge of the screen.", float32(0.1), nil)
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

// It creates a new camera with the necessary setup from settings screen
func CreateCameraFromSettings() *camera.DefaultCamera {
	cameraPosition := Settings["CameraPos"].GetCurrentValue().(mgl32.Vec3)
	worldUp := Settings["WorldUp"].GetCurrentValue().(mgl32.Vec3)
	yawAngle := Settings["CameraYaw"].GetCurrentValue().(float32)
	pitchAngle := Settings["CameraPitch"].GetCurrentValue().(float32)
	fov := Settings["CameraFov"].GetCurrentValue().(float32)
	near := Settings["CameraNear"].GetCurrentValue().(float32)
	far := Settings["CameraFar"].GetCurrentValue().(float32)
	moveSpeed := Settings["CameraVelocity"].GetCurrentValue().(float32)
	directionSpeed := Settings["CameraRotation"].GetCurrentValue().(float32)
	camera := camera.NewCamera(cameraPosition, worldUp, yawAngle, pitchAngle)
	camera.SetupProjection(fov, float32(WindowWidth)/float32(WindowHeight), near, far)
	camera.SetVelocity(moveSpeed)
	camera.SetRotationStep(directionSpeed)
	return camera
}

func CreateMaterialCube(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	cube := cuboid.NewCube()
	v, i, bo := cube.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	m.SetBoundingObject(bo)
	return m
}
func CreateMaterialCylinder(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	c := cylinder.New(0.75, 30, 3)
	v, i, bo := c.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	m.SetBoundingObject(bo)
	return m
}

func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	lightSourceRotationAngleRadian := mgl32.DegToRad((360 / Settings["LSRoundSpeed"].GetCurrentValue().(float32)) * float32(delta))
	lightDirectionRotationMatrix := mgl32.HomogRotate3D(lightSourceRotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentLightSourceDirection := LightSourceCube.GetDirection()
	LightSourceCube.SetDirection(mgl32.TransformNormal(currentLightSourceDirection, lightDirectionRotationMatrix))
	PointLightSource.SetPosition(LightSourceCube.GetPosition())

	app.Update(delta)
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
func createSettings(defaults config.Config) *screen.FormScreen {
	formItemOrders := []string{
		"ClearCol",
		"LSAmbient",
		"LSDiffuse",
		"LSSpecular",
		"LSPosition",
		"LSConstantTerm", "LSLinearTerm",
		"LSQuadraticTerm",

		"JadePosition",
		"JadeScale",
		"RPPosition",
		"RPScale",
		"ObsidianPosition",
		"ObsidianScale",
		"CopperPosition",
		"CopperScale",
		"SilverPosition",
		"SilverScale",
		"TurquoisePosition",
		"TurquoiseScale",
		"LSPosition",
		"LSRoundSpeed",

		"CameraPos",
		"WorldUp",
		"CameraYaw", "CameraPitch",
		"CameraNear", "CameraFar",
		"CameraFov", "CameraVelocity",
		"CameraRotation", "CameraRotationEdge",
	}
	return app.BuildFormScreen(defaults, formItemOrders, "Settings")
}
func createMenu() *screen.MenuScreen {
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
	options := []screen.Option{
		*screen.NewMenuScreenOption("continue", contS, continueEvent),
		*screen.NewMenuScreenOption("start", contNS, startEvent),
		*screen.NewMenuScreenOption("restart", contS, restartEvent),
		*screen.NewMenuScreenOption("settings", contAll, settingsEvent),
		*screen.NewMenuScreenOption("exit", contAll, exitEvent),
	}
	return app.BuildMenuScreen(options)
}

func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCameraFromSettings(), CameraMovementOptions())

	PointLightSource = light.NewPointLight([4]mgl32.Vec3{
		Settings["LSPosition"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["LSSpecular"].GetCurrentValue().(mgl32.Vec3)},
		[3]float32{
			Settings["LSConstantTerm"].GetCurrentValue().(float32),
			Settings["LSLinearTerm"].GetCurrentValue().(float32),
			Settings["LSQuadraticTerm"].GetCurrentValue().(float32),
		})

	// Add the lightources to the application
	scrn.AddPointLightSource(PointLightSource, [7]string{"light.position", "light.ambient", "light.diffuse", "light.specular", "light.constant", "light.linear", "light.quadratic"})

	materialShader := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(materialShader)

	Model := model.New()
	JadeCube = CreateMaterialCube(material.Jade, Settings["JadePosition"].GetCurrentValue().(mgl32.Vec3))
	JadeCube.SetScale(Settings["JadeScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(JadeCube)

	rpCube := CreateMaterialCube(material.Redplastic, Settings["RPPosition"].GetCurrentValue().(mgl32.Vec3))
	rpCube.SetScale(Settings["RPScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(rpCube)

	obsidianCube := CreateMaterialCube(material.Obsidian, Settings["ObsidianPosition"].GetCurrentValue().(mgl32.Vec3))
	obsidianCube.SetScale(Settings["ObsidianScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(obsidianCube)

	copperCube := CreateMaterialCube(material.Copper, Settings["CopperPosition"].GetCurrentValue().(mgl32.Vec3))
	copperCube.SetScale(Settings["CopperScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(copperCube)

	silverCube := CreateMaterialCube(material.Silver, Settings["SilverPosition"].GetCurrentValue().(mgl32.Vec3))
	silverCube.SetScale(Settings["SilverScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(silverCube)

	turqCylinder := CreateMaterialCylinder(material.Turquoise, Settings["TurquoisePosition"].GetCurrentValue().(mgl32.Vec3))
	turqCylinder.SetScale(Settings["TurquoiseScale"].GetCurrentValue().(mgl32.Vec3))
	Model.AddMesh(turqCylinder)

	mat := material.New(mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, mgl32.Vec3{1, 1, 1}, 144.0)
	LightSourceCube = CreateMaterialCube(mat, Settings["LSPosition"].GetCurrentValue().(mgl32.Vec3))
	LightSourceCube.SetDirection((mgl32.Vec3{9, 0, -3}).Normalize())

	distance := (LightSourceCube.GetPosition().Sub(JadeCube.GetPosition())).Len()
	LightSourceCube.SetSpeed((float32(2) * float32(3.1415) * distance) / Settings["LSRoundSpeed"].GetCurrentValue().(float32))
	Model.AddMesh(LightSourceCube)
	scrn.AddModelToShader(Model, materialShader)
	scrn.Setup(setupApp)
	return scrn
}
func AddFormScreen() bool {
	val := os.Getenv(FORM_ENV_NAME)
	if val == ON_VALUE {
		return true
	}
	return false
}
func main() {
	app = application.New(glWrapper)
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	AppScreen = mainScreen()
	app.AddScreen(AppScreen)
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	if AddFormScreen() {
		app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
		app.GetWindow().SetCharCallback(app.CharCallback)
		MenuScreen = createMenu()
		app.AddScreen(MenuScreen)
		app.MenuScreen(MenuScreen)
		SettingsScreen = createSettings(Settings)
		app.AddScreen(SettingsScreen)
		app.ActivateScreen(SettingsScreen)
	} else {
		app.ActivateScreen(AppScreen)
	}
	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
