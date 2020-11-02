package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	WindowTitle = "Example cube form screen"
)

var (
	app        *application.Application
	lastUpdate int64
	// window related variables
	Builder          *window.WindowBuilder
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowFullScreen = false
	glWrapper        glwrapper.Wrapper
)

func init() {
	runtime.LockOSThread()
	setupWindowBuilder()
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
func buildScreen() *screen.CubeFormScreen {
	b := screen.NewCubeFormScreenBuilder()
	b.SetWrapper(glWrapper)
	b.SetAssetsDirectory("./examples/13-cube-menu/assets")
	b.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
	return b.Build()
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.Update(delta)
}
func main() {

	app = application.New(glWrapper)
	// Setup the window
	app.SetWindow(Builder.Build())
	// Terminate window at the end.
	defer glfw.Terminate()
	// Init opengl.
	glWrapper.InitOpenGL()

	scrn := buildScreen()
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
	app.GetWindow().SetCharCallback(app.CharCallback)
	lastUpdate = time.Now().UnixNano()

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
