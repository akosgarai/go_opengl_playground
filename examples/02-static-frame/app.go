package main

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	app *application.Application

	glWrapper glwrapper.Wrapper

	Builder          *window.WindowBuilder
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowTitle      = "Example - static frame"
	WindowFullScreen = false
	Aspect           = false
)

func init() {
	runtime.LockOSThread()

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
	title := os.Getenv("TITLE")
	if title != "" {
		WindowTitle = title
	}
	Builder = window.NewWindowBuilder()
	fullScreen := os.Getenv("FULL")
	if fullScreen == "1" {
		WindowFullScreen = true
		WindowWidth, WindowHeight = Builder.GetCurrentMonitorResolution()
	}
	aspect := os.Getenv("ASPECT")
	if aspect != "" {
		Aspect = true

	}
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(1.0, 1.0, 0.0, 1.0)
	glWrapper.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
}

func mainScreen() *screen.Screen {
	scrn := screen.New()
	scrn.Setup(setupApp)
	// add shader program
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)
	// models: 4 rectangle.
	// - top and bottom: 2 * 0.5, vertical positions: 0.75, -0.75, color: 0,0,1
	// - left, right: 1 * 0.5 (2 - 2*widthofthetopbottom), horizontal positions: 0.75, -0.75, color: 0,1,1
	widthMul := float32(1.0)
	heightMul := float32(1.0)
	if Aspect {
		if WindowWidth > WindowHeight {
			widthMul = float32(WindowHeight) / float32(WindowWidth)
		}
		if WindowWidth < WindowHeight {
			heightMul = float32(WindowWidth) / float32(WindowHeight)
		}
	}
	mod := model.New()
	verticalFrameWidth := float32(2.0)
	verticalFrameHeight := float32(0.5) * heightMul
	squareBaseVertical := rectangle.NewExact(verticalFrameWidth, verticalFrameHeight)
	verticalFrameColor := []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	v, i, _ := squareBaseVertical.ColoredMeshInput(verticalFrameColor)
	// top mesh
	mshTop := mesh.NewColorMesh(v, i, verticalFrameColor, glWrapper)
	mshTop.SetPosition(mgl32.Vec3{0, 1 - (verticalFrameHeight / 2), 0})
	mod.AddMesh(mshTop)
	// bottom mesh
	mshBottom := mesh.NewColorMesh(v, i, verticalFrameColor, glWrapper)
	mshBottom.SetPosition(mgl32.Vec3{0, -1 + (verticalFrameHeight / 2), 0})
	mod.AddMesh(mshBottom)

	horizontalFrameWidth := float32(0.5) * widthMul
	horizontalFrameHeight := float32(2.0) - 2*verticalFrameHeight
	squareBaseHorizontal := rectangle.NewExact(horizontalFrameWidth, horizontalFrameHeight)
	horizontalFrameColor := []mgl32.Vec3{mgl32.Vec3{0, 1, 1}}
	v, i, _ = squareBaseHorizontal.ColoredMeshInput(horizontalFrameColor)
	// left mesh
	mshLeft := mesh.NewColorMesh(v, i, horizontalFrameColor, glWrapper)
	mshLeft.SetPosition(mgl32.Vec3{-1 + (horizontalFrameWidth / 2), 0, 0})
	mod.AddMesh(mshLeft)
	// right mesh
	mshRight := mesh.NewColorMesh(v, i, horizontalFrameColor, glWrapper)
	mshRight.SetPosition(mgl32.Vec3{1 - (horizontalFrameWidth / 2), 0, 0})
	mod.AddMesh(mshRight)
	mod.RotateX(90)
	scrn.AddModelToShader(mod, shaderProgram)
	return scrn
}

func main() {
	app = application.New(glWrapper)

	Builder.SetFullScreen(WindowFullScreen)
	Builder.SetDecorated(WindowDecorated)
	Builder.SetTitle(WindowTitle)
	Builder.SetWindowSize(WindowWidth, WindowHeight)
	Builder.PrintCurrentMonitorData()

	app.SetWindow(Builder.Build())
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := mainScreen()
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
