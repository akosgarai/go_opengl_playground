package main

import (
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
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
	WindowTitle  = "Example - menu screen with model"
	FontFile     = "/assets/fonts/Desyrel/desyrel.regular.ttf"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
	LEFT_MOUSE_BUTTON    = glfw.MouseButtonLeft
)

var (
	app        *application.Application
	MenuScreen *screen.Screen
	AppScreen  *screen.Screen
	glWrapper  glwrapper.Wrapper

	StartButton *mesh.TexturedMaterialMesh
	ExitButton  *mesh.TexturedMaterialMesh
	Wall        *mesh.TexturedMaterialMesh

	lastUpdate int64

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby
)

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

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-0.28, -0.23, 2.4}, mgl32.Vec3{0, -1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
	_, msh, distance := app.GetClosestModelMeshDistance()
	switch msh.(type) {
	case *mesh.TexturedMaterialMesh:
		tmMesh := msh.(*mesh.TexturedMaterialMesh)
		if distance < 0.01 {
			tmMesh.Material = HighlightMaterial
			if app.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
				if tmMesh == ExitButton {
					fmt.Println("Exit button has been pressed.\n")
					app.GetWindow().SetShouldClose(true)
				} else if tmMesh == StartButton {
					fmt.Println("Start button has been pressed.\n")
					app.ActivateScreen(AppScreen)
					glWrapper.ClearColor(1.0, 1.0, 0.0, 1.0)
				}
			}
		} else if distance < 1.8 && tmMesh == Wall {
			tmMesh.Material = HighlightMaterial
		} else {
			tmMesh.Material = DefaultMaterial
		}
		break
	}

}

func main() {
	runtime.LockOSThread()
	app = application.New()
	Window := window.InitGlfw(WindowWidth, WindowHeight, WindowTitle)
	app.SetWindow(Window)
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	MenuScreen = screen.New()
	AppScreen = screen.New()

	AppScreen.SetCamera(CreateCamera())
	AppScreen.SetCameraMovementMap(CameraMovementMap())
	AppScreen.SetRotateOnEdgeDistance(CameraDistance)

	fontShader := shader.NewShader(baseDir()+"/shaders/font.vert", baseDir()+"/shaders/font.frag", glWrapper)
	MenuScreen.AddShader(fontShader)
	AppScreen.AddShader(fontShader)
	paperShader := shader.NewShader(baseDir()+"/shaders/paper.vert", baseDir()+"/shaders/paper.frag", glWrapper)
	MenuScreen.AddShader(paperShader)
	AppScreen.AddShader(paperShader)

	paperModel := model.New()
	StartButton = Paper(1, 0.2, mgl32.Vec3{-0.0, 0.3, -0.0})
	StartButton.RotateX(-90)
	paperModel.AddMesh(StartButton)
	ExitButton = Paper(1, 0.2, mgl32.Vec3{-0.0, -0.3, -0.0})
	ExitButton.RotateX(-90)
	paperModel.AddMesh(ExitButton)
	MenuScreen.AddModelToShader(paperModel, paperShader)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
	StartableModel := model.New()
	Wall = Paper(2, 2, mgl32.Vec3{-0.4, -0.3, -0.0})
	Wall.RotateX(-90)
	StartableModel.AddMesh(Wall)

	AppScreen.AddModelToShader(StartableModel, paperShader)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	MenuFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
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
	AppScreen.AddModelToShader(StartableFonts, fontShader)
	MenuFonts.PrintTo(" - Start - ", -0.4, -0.03, 0.01, 3.0/float32(WindowWidth), glWrapper, StartButton, cols1)
	MenuFonts.PrintTo(" - Exit - ", -0.4, -0.03, 0.01, 3.0/float32(WindowWidth), glWrapper, ExitButton, cols2)
	MenuFonts.SetTransparent(true)
	MenuScreen.AddModelToShader(MenuFonts, fontShader)
	lastUpdate = time.Now().UnixNano()
	app.AddScreen(MenuScreen)
	app.AddScreen(AppScreen)
	app.MenuScreen(MenuScreen)
	app.ActivateScreen(MenuScreen)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
