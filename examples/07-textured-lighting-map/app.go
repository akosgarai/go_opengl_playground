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
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
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
	WindowTitle  = "Example - textured lighting map"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app            *application.Application
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()

	lastUpdate int64
	startTime  int64

	PointLightSource  *light.Light
	LightSourceSphere *mesh.MaterialMesh
	CubeMesh          *mesh.TexturedMesh

	glWrapper glwrapper.Wrapper
)

func InitSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.3, 0.3, 0.3}, colorValidator)
	// light source
	Settings.AddConfig("LSAmbient", "Light ambient", "The ambient color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSDiffuse", "Light diffuse", "The diffuse color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSSpecular", "Light specular", "The specular color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSPosition", "Light position", "The position vector of the lightsource.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	// sphere
	Settings.AddConfig("MSPosition", "Sphere position", "The position vector of the sphere.", mgl32.Vec3{-3.0, -0.5, -3.0}, nil)
	Settings.AddConfig("MSScale", "Sphere scale", "The scale vector of the sphere.", mgl32.Vec3{0.15, 0.15, 0.15}, nil)
	Settings.AddConfig("MSAmbient", "Sphere mat a.", "The ambient color component of the sphere.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("MSDiffuse", "Sphere mat d.", "The diffuse color component of the sphere.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("MSSpecular", "Sphere mat s.", "The specular color component of the sphere.", mgl32.Vec3{1.0, 1.0, 1.0}, colorValidator)
	Settings.AddConfig("MSShininess", "Sphere mat sh.", "The shininess of the sphere material.", float32(3000.0), nil)
	Settings.AddConfig("MSPrecision", "Sphere precision", "The precision of the sphere.", int(15), nil)
	Settings.AddConfig("MSRoundSpeed", "Sphere r. s.", "The round speed of the sphere. It goes one round in this interval.", float32(3000.0), nil)
	// cylinder
	Settings.AddConfig("CylinderPosition", "Cylinder position", "The position of the cylinder.", mgl32.Vec3{2, 2, 2}, nil)
	Settings.AddConfig("CylinderRad", "Cylinder rad", "The radius of the cylinder mesh.", float32(1.5), nil)
	Settings.AddConfig("CylinderLength", "Cylinder length", "The length of the cylinder mesh.", float32(3.0), nil)
	Settings.AddConfig("CylinderPrec", "Cylinder prec", "The precision of the cylinder mesh.", int(30), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{3.3, -10, 14.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
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
func CreateCameraFromSettings() *camera.Camera {
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

// It generates the lightsource sphere.
func CreateWhiteSphere() *mesh.MaterialMesh {
	mat := material.New(
		Settings["MSAmbient"].GetCurrentValue().(mgl32.Vec3),
		Settings["MSDiffuse"].GetCurrentValue().(mgl32.Vec3),
		Settings["MSSpecular"].GetCurrentValue().(mgl32.Vec3),
		Settings["MSShininess"].GetCurrentValue().(float32))
	sph := sphere.New(Settings["MSPrecision"].GetCurrentValue().(int))
	v, i, _ := sph.MaterialMeshInput()
	LightSourceSphere := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	LightSourceSphere.SetPosition(Settings["MSPosition"].GetCurrentValue().(mgl32.Vec3))
	LightSourceSphere.SetDirection((mgl32.Vec3{9, 0, -3}).Normalize())
	distance := (LightSourceSphere.GetPosition().Sub(CubeMesh.GetPosition())).Len()
	LightSourceSphere.SetSpeed((float32(2) * float32(3.1415) * distance) / Settings["MSRoundSpeed"].GetCurrentValue().(float32))
	LightSourceSphere.SetScale(Settings["MSScale"].GetCurrentValue().(mgl32.Vec3))
	return LightSourceSphere
}

// It generates a cube.
func CreateCubeMesh(t texture.Textures) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	return m
}
func CreateCylinderMesh(t texture.Textures) *mesh.TexturedMesh {
	rad := Settings["CylinderRad"].GetCurrentValue().(float32)
	len := Settings["CylinderLength"].GetCurrentValue().(float32)
	prec := Settings["CylinderPrec"].GetCurrentValue().(int)
	c := cylinder.New(rad, prec, len)
	v, i, _ := c.TexturedMeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(Settings["CylinderPosition"].GetCurrentValue().(mgl32.Vec3))
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
	lightSourceRotationAngleRadian := mgl32.DegToRad(360.0 / Settings["MSRoundSpeed"].GetCurrentValue().(float32) * float32(delta))
	lightDirectionRotationMatrix := mgl32.HomogRotate3D(lightSourceRotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentLightSourceDirection := LightSourceSphere.GetDirection()
	LightSourceSphere.SetDirection(mgl32.TransformNormal(currentLightSourceDirection, lightDirectionRotationMatrix))
	PointLightSource.SetPosition(LightSourceSphere.GetPosition())

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

		"MSPosition",
		"MSScale",
		"MSAmbient",
		"MSDiffuse",
		"MSSpecular",
		"MSShininess", "MSPrecision",
		"MSRoundSpeed",

		"CylinderPosition",
		"CylinderRad", "CylinderLength",
		"CylinderPrec",

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
	scrn.SetCamera(CreateCameraFromSettings())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(Settings["CameraRotationEdge"].GetCurrentValue().(float32))

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

	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/colored-image-for-texture-testing-diffuse.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	tex.AddTexture(baseDir()+"/assets/colored-image-for-texture-testing-specular.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	TexModel := model.New()
	MatModel := model.New()

	CubeMesh = CreateCubeMesh(tex)
	TexModel.AddMesh(CubeMesh)
	TexModel.AddMesh(CreateCylinderMesh(tex))
	scrn.AddModelToShader(TexModel, shaderProgramTexture)

	shaderProgramWhite := shader.NewShader(baseDir()+"/shaders/lightsource.vert", baseDir()+"/shaders/lightsource.frag", glWrapper)
	scrn.AddShader(shaderProgramWhite)

	LightSourceSphere = CreateWhiteSphere()
	MatModel.AddMesh(LightSourceSphere)
	scrn.AddModelToShader(MatModel, shaderProgramWhite)
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
	runtime.LockOSThread()
	InitSettings()

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
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
