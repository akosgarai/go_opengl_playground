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
	WindowTitle  = "Example - textured spheres"

	FORM_ENV_NAME = "SETTINGS"
	ON_VALUE      = "on"
)

var (
	app            *application.Application
	Sun            *mesh.TexturedMesh
	Earth          *mesh.TexturedMesh
	MatPlanet      *mesh.TexturedMaterialMesh
	SettingsScreen *screen.FormScreen
	MenuScreen     *screen.MenuScreen
	AppScreen      *screen.Screen
	Settings       = config.New()
	lastUpdate     int64
	startTime      int64

	rotationAngle   = float32(0.0)
	spherePrimitive = sphere.New(20)

	glWrapper glwrapper.Wrapper
)

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
func InitSettings() {
	var colorValidator model.FloatValidator
	colorValidator = func(f float32) bool { return f >= 0 && f <= 1 }
	Settings.AddConfig("ClearCol", "BG color", "The clear color of the window. It is used as the color of the background.", mgl32.Vec3{0.0, 0.0, 0.0}, colorValidator)
	// light source
	Settings.AddConfig("LSAmbient", "Light ambient", "The ambient color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSDiffuse", "Light diffuse", "The diffuse color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSSpecular", "Light specular", "The specular color component of the lightsource.", mgl32.Vec3{0.5, 0.5, 0.5}, colorValidator)
	Settings.AddConfig("LSPosition", "Light position", "The position vector of the lightsource.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("LSConstantTerm", "Constant term", "The constant term of the light equasion.", float32(1.0), nil)
	Settings.AddConfig("LSLinearTerm", "Linear term", "The linear term of the light equasion.", float32(0.14), nil)
	Settings.AddConfig("LSQuadraticTerm", "Quadratic term", "The quadratic term of the light equasion.", float32(0.07), nil)
	// sun
	Settings.AddConfig("SunPosition", "Sun position", "The center point of the sun.", mgl32.Vec3{0.0, 0.0, 0.0}, nil)
	Settings.AddConfig("SunRadius", "Sun radius", "The radius of the sun. This value is used for scaling.", float32(0.1), nil)
	Settings.AddConfig("SunRoundSpeed", "Sun r. s.", "The round speed of the sun.", float32(0.01), nil)
	// earth
	Settings.AddConfig("EarthPosition", "Earth position", "The center point of the earth.", mgl32.Vec3{0.4, 0.0, 0.0}, nil)
	Settings.AddConfig("EarthRadius", "Earth radius", "The radius of the earth. This value is used for scaling.", float32(0.01), nil)
	Settings.AddConfig("EarthRoundSpeed", "Earth r. s.", "The round speed of the earth planet. It goes one round in this interval.", float32(3000.0), nil)
	// material planet
	Settings.AddConfig("MaterialPosition", "Planet position", "The center point of the material planet.", mgl32.Vec3{0.2, 0.0, 0.0}, nil)
	Settings.AddConfig("MaterialRadius", "Planet radius", "The radius of the material planet. This value is used for scaling.", float32(0.02), nil)
	Settings.AddConfig("MaterialRoundSpeed", "Planet r. s.", "The round speed of the material planet. It goes one round in this interval.", float32(7000.0), nil)
	// skybox distance
	Settings.AddConfig("SkyboxDistance", "Sky distance", "The distance of the skybox from the origo. This value is used for scaling.", float32(20.0), nil)
	// camera options:
	// - position
	Settings.AddConfig("CameraPos", "Cam position", "The initial position of the camera.", mgl32.Vec3{0.0, 0.0, -2.0}, nil)
	// - up direction
	Settings.AddConfig("WorldUp", "World up dir", "The up direction in the world.", mgl32.Vec3{0.0, 1.0, 0.0}, nil)
	// - pitch, yaw
	Settings.AddConfig("CameraYaw", "Cam Yaw", "The yaw (angle) of the camera. Rotation on the Z axis.", float32(90.0), nil)
	Settings.AddConfig("CameraPitch", "Cam Pitch", "The pitch (angle) of the camera. Rotation on the Y axis.", float32(0.0), nil)
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

func TexturedSphere(t texture.Textures, position mgl32.Vec3, scale float32) *mesh.TexturedMesh {
	v, i, _ := spherePrimitive.TexturedMeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(position)
	m.SetScale(mgl32.Vec3{scale, scale, scale})
	return m
}
func TexturedMaterialSphere(t texture.Textures, mat *material.Material, position mgl32.Vec3, scale float32) *mesh.TexturedMaterialMesh {
	v, i, _ := spherePrimitive.TexturedMeshInput()
	m := mesh.NewTexturedMaterialMesh(v, i, t, mat, glWrapper)
	m.SetPosition(position)
	m.SetScale(mgl32.Vec3{scale, scale, scale})
	return m
}

// It generates a cube map.
func CubeMap(t texture.Textures) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i, _ := cube.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	return m
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
func updateSun(moveTime float64) {
	Sun.RotateY(float32(moveTime) * Settings["SunRoundSpeed"].GetCurrentValue().(float32))
}
func updatePlanets(moveTime float64) {
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	rotationAngleRadian := mgl32.DegToRad(360.0 / Settings["EarthRoundSpeed"].GetCurrentValue().(float32) * float32(moveTime))
	rotationMatrix := mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentDirection := Earth.GetDirection()
	Earth.SetDirection(mgl32.TransformNormal(currentDirection, rotationMatrix))

	rotationAngleRadian = mgl32.DegToRad(360 / Settings["MaterialRoundSpeed"].GetCurrentValue().(float32) * float32(moveTime))
	rotationMatrix = mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentDirection = MatPlanet.GetDirection()
	MatPlanet.SetDirection(mgl32.TransformNormal(currentDirection, rotationMatrix))
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	updatePlanets(delta)
	updateSun(delta)

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
		"SunPosition",
		"SunRadius", "SunRoundSpeed",
		"EarthPosition",
		"EarthRadius", "EarthRoundSpeed",
		"MaterialPosition",
		"MaterialRadius", "MaterialRoundSpeed",
		"SkyboxDistance",

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

	PointLightSource := light.NewPointLight([4]mgl32.Vec3{
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
	scrn.AddPointLightSource(PointLightSource, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader(baseDir()+"/shaders/texture.vert", baseDir()+"/shaders/texture.frag", glWrapper)
	scrn.AddShader(shaderProgramTexture)

	TexModel := model.New()
	TexMatModel := model.New()
	CmModel := model.New()
	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture(baseDir()+"/assets/sun.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	Sun = TexturedSphere(sunTexture, Settings["SunPosition"].GetCurrentValue().(mgl32.Vec3), Settings["SunRadius"].GetCurrentValue().(float32))
	TexModel.AddMesh(Sun)
	// sun texture
	var earthTexture texture.Textures
	earthTexture.AddTexture(baseDir()+"/assets/earth.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	earthTexture.AddTexture(baseDir()+"/assets/earth.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	Earth = TexturedSphere(sunTexture, Settings["EarthPosition"].GetCurrentValue().(mgl32.Vec3), Settings["EarthRadius"].GetCurrentValue().(float32))
	distance := Earth.GetPosition().Len()
	Earth.SetSpeed((float32(2) * float32(3.1415) * distance) / Settings["EarthRoundSpeed"].GetCurrentValue().(float32))
	Earth.SetDirection((mgl32.Vec3{0, 0, 1}).Normalize())
	TexModel.AddMesh(Earth)
	scrn.AddModelToShader(TexModel, shaderProgramTexture)
	// other planet texture
	shaderProgramTextureMaterial := shader.NewShader(baseDir()+"/shaders/texturemat.vert", baseDir()+"/shaders/texturemat.frag", glWrapper)
	scrn.AddShader(shaderProgramTextureMaterial)
	var materialTexture texture.Textures
	materialTexture.AddTexture(baseDir()+"/assets/venus.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	materialTexture.AddTexture(baseDir()+"/assets/venus.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", glWrapper)
	MatPlanet = TexturedMaterialSphere(sunTexture, material.Gold, Settings["MaterialPosition"].GetCurrentValue().(mgl32.Vec3), Settings["MaterialRadius"].GetCurrentValue().(float32))
	distance = MatPlanet.GetPosition().Len()
	MatPlanet.SetSpeed((float32(2) * float32(3.1415) * distance) / Settings["MaterialRoundSpeed"].GetCurrentValue().(float32))
	MatPlanet.SetDirection((mgl32.Vec3{0, 0, 1}).Normalize())
	TexMatModel.AddMesh(MatPlanet)
	scrn.AddModelToShader(TexMatModel, shaderProgramTextureMaterial)

	shaderProgramCubeMap := shader.NewShader(baseDir()+"/shaders/cubeMap.vert", baseDir()+"/shaders/cubeMap.frag", glWrapper)
	scrn.AddShader(shaderProgramCubeMap)
	var cubeMapTexture texture.Textures
	cubeMapTexture.AddCubeMapTexture(baseDir()+"/assets", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "skybox", glWrapper)
	cubeMap := CubeMap(cubeMapTexture)
	d := Settings["SkyboxDistance"].GetCurrentValue().(float32)
	cubeMap.SetScale(mgl32.Vec3{d, d, d})
	CmModel.AddMesh(cubeMap)
	scrn.AddModelToShader(CmModel, shaderProgramCubeMap)
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
