package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
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
	WindowTitle  = "Example - menu screen with screen"
	FontFile     = "/assets/fonts/Desyrel/desyrel.regular.ttf"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
	LEFT_MOUSE_BUTTON    = glfw.MouseButtonLeft
)

var (
	app        *application.Application
	AppScreen  *screen.Screen
	glWrapper  glwrapper.Wrapper
	lastUpdate int64

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby
)

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
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

func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(1.0, 1.0, 0.0, 1.0)
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
}
func Paper(width, height float32, position mgl32.Vec3) *mesh.TexturedMaterialMesh {
	rect := rectangle.NewExact(width, height)
	v, i, bo := rect.MeshInput()
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)

	msh := mesh.NewTexturedMaterialMesh(v, i, tex, DefaultMaterial, glWrapper)
	msh.SetBoundingObject(bo)
	msh.SetPosition(position)
	return msh
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-0.28, -0.23, 2.4}, mgl32.Vec3{0, -1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func GameScreen() *screen.Screen {
	gameScreen := screen.New()
	gameScreen.SetCamera(CreateCamera())
	gameScreen.SetCameraMovementMap(CameraMovementMap())
	gameScreen.SetRotateOnEdgeDistance(CameraDistance)
	bgShaderApplication := shader.NewMenuBackgroundShader(glWrapper)
	fgShaderApplication := shader.NewFontShader(glWrapper)
	gameScreen.AddShader(bgShaderApplication)
	gameScreen.AddShader(fgShaderApplication)

	StartableModel := model.New()
	Wall := Paper(2, 2, mgl32.Vec3{-0.4, -0.3, -0.0})
	Wall.RotateX(-90)
	StartableModel.AddMesh(Wall)

	gameScreen.AddModelToShader(StartableModel, bgShaderApplication)
	StartableFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	cols1 := []mgl32.Vec3{
		mgl32.Vec3{0.0, 1.0, 0.0},
	}
	cols2 := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 1.0},
	}
	cols3 := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 0.0},
	}
	StartableFonts.PrintTo("How are You?", -0.5, 0.2, -0.01, 3.0/float32(WindowWidth), glWrapper, Wall, cols2)
	StartableFonts.PrintTo("Press Esc for Menu!", -0.7, -0.2, -0.01, 3.0/float32(WindowWidth), glWrapper, Wall, cols3)
	StartableFonts.PrintTo("Ken sent me!", -0.2, -0.75, -0.01, 3.0/float32(WindowWidth), glWrapper, Wall, cols1)
	StartableFonts.SetTransparent(true)
	gameScreen.AddModelToShader(StartableFonts, fgShaderApplication)
	return gameScreen
}

func main() {
	runtime.LockOSThread()
	app = application.New(glWrapper)
	Window := window.InitGlfw(WindowWidth, WindowHeight, WindowTitle)
	app.SetWindow(Window)
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	StartableFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	StartableFonts.SetTransparent(true)
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)
	menu := screen.NewMenuScreen(tex, DefaultMaterial, HighlightMaterial, StartableFonts, mgl32.Vec3{1, 0, 0}, mgl32.Vec3{1, 0, 0}, glWrapper)
	form := screen.NewFormScreen(material.Ruby, "Settings", glWrapper, float32(WindowWidth), float32(WindowHeight))
	form.AddFormItemBool("Bool 01.", glWrapper, true)
	form.AddFormItemBool("Bool 02.", glWrapper, false)
	form.AddFormItemBool("Bool 03.", glWrapper, false)
	form.AddFormItemInt("Int 04.", glWrapper, "0", nil)
	form.AddFormItemInt("Int 05.", glWrapper, "12", nil)
	form.AddFormItemInt("Int 06.", glWrapper, "-3", nil)
	form.AddFormItemBool("Bool 07.", glWrapper, true)
	form.AddFormItemFloat("Float 08.", glWrapper, "0.0", nil)
	form.AddFormItemFloat("Float 09.", glWrapper, "1.876", nil)
	form.AddFormItemFloat("Float 10.", glWrapper, "-0.44", nil)
	form.AddFormItemText("Text 11.", glWrapper, "Some", nil)
	form.AddFormItemText("Text 12.", glWrapper, "sample", nil)
	form.AddFormItemText("Text 13.", glWrapper, "text", nil)
	form.AddFormItemInt64("Int64 14.", glWrapper, "0", nil)
	form.AddFormItemInt64("Int64 15.", glWrapper, "1231234", nil)
	form.AddFormItemInt64("Int64 16.", glWrapper, "-1239876", nil)
	AppScreen = GameScreen()
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
		AppScreen.SetCamera(CreateCamera())
		app.ActivateScreen(AppScreen)
	}
	startEvent := func() {
		app.ActivateScreen(AppScreen)
		menu.SetState("world-started", true)
		menu.BuildScreen(glWrapper, 3/float32(WindowWidth))
	}
	settingsEvent := func() {
		app.ActivateScreen(form)
	}
	continueEvent := func() {
		app.ActivateScreen(AppScreen)
	}
	exitEvent := func() {
		app.GetWindow().SetShouldClose(true)
	}
	cont := screen.NewMenuScreenOption("continue", contS, continueEvent)
	menu.AddOption(*cont) // continue
	start := screen.NewMenuScreenOption("start", contNS, startEvent)
	menu.AddOption(*start) // start
	restart := screen.NewMenuScreenOption("restart", contS, restartEvent)
	menu.AddOption(*restart) // restart
	settings := screen.NewMenuScreenOption("settings", contAll, settingsEvent)
	menu.AddOption(*settings) // settings
	exit := screen.NewMenuScreenOption("exit", contAll, exitEvent)
	menu.AddOption(*exit) // exit
	menu.BuildScreen(glWrapper, 3/float32(WindowWidth))

	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.SetWrapper(glWrapper)
	app.GetWindow().SetCharCallback(app.CharCallback)
	app.AddScreen(AppScreen)
	app.AddScreen(menu)
	app.AddScreen(form)
	app.MenuScreen(menu)
	app.ActivateScreen(menu)

	lastUpdate = time.Now().UnixNano()

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		Update()
		app.GetWindow().SwapBuffers()
	}
}
