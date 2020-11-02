package main

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowTitle = "Example - real-time editor"
)

var (
	app *application.Application

	// window related variables
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowFullScreen = false

	glWrapper glwrapper.Wrapper
)

func init() {
	// lock thread
	runtime.LockOSThread()
	app = application.New(glWrapper)
	// Setup the window and the WindowBuilder
	app.SetWindow(setupWindowBuilder().Build())
	glWrapper.InitOpenGL()
	// application screen
	scrn := CreateApplicationScreen()
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)
}
func setupWindowBuilder() *window.WindowBuilder {
	Builder := window.NewWindowBuilder()
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
	return Builder
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.3, 0.3, 0.3, 1.0)
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.DefaultCamera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0, 1, 0}, 0.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 10.0)
	return camera
}

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["mode"] = "default"
	cm["rotateOnEdgeDistance"] = float32(0.0)
	return cm
}

// It generates the Jade sphere.
func CreateJadeSphere() *mesh.MaterialMesh {
	sph := sphere.New(20)
	v, i, _ := sph.MaterialMeshInput()
	JadeSphere := mesh.NewMaterialMesh(v, i, material.Jade, glWrapper)
	JadeSphere.SetPosition(mgl32.Vec3{3.0, 0.0, 0.0})
	return JadeSphere
}
func CreateApplicationScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCamera(), CameraMovementOptions())
	shaderProgram := shader.NewMaterialShader(glWrapper)
	scrn.AddShader(shaderProgram)
	Model := model.New()
	Model.AddMesh(CreateJadeSphere())
	scrn.AddModelToShader(Model, shaderProgram)
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
	})
	// Add the lightources to the application
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.Setup(setupApp)
	return scrn
}

func main() {
	defer glfw.Terminate()

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
