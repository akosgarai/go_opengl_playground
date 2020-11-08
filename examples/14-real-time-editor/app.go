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
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
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
	app       *application.Application
	AppScreen *EditorScreen
	// camera & button
	lastUpdate int64
	lastToggle int64

	// window related variables
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowFullScreen = false

	glWrapper glwrapper.Wrapper
	// menu screen flag
	MenuScreenEnabled = true
	// Button metric
	frame              = mgl32.Vec2{0.2, 0.1}
	surface            = mgl32.Vec2{0.185, 0.0925}
	buttonDefaultColor = []mgl32.Vec3{mgl32.Vec3{0.4, 0.4, 0.4}}
	buttonHoverColor   = []mgl32.Vec3{mgl32.Vec3{0.4, 0.8, 0.4}}
)

/*
 * Now i want to extend the application with the button model. It supposed to contain 2 meshes (bg, surface), should have the update options.
 * The menuModel should be extended with the menu items. It could be an array of models with some support functions.
 */

// It is the representation of the button ui item.
type Button struct {
	*model.BaseModel
	defaultColor []mgl32.Vec3
	hoverColor   []mgl32.Vec3
	frameSize    mgl32.Vec2
	surfaceSize  mgl32.Vec2
	screen       interfaces.Mesh
	position     mgl32.Vec3
}

// PinToScreen sets the parent of the bg mesh to the given one and updates its position.
func (b *Button) PinToScreen(scrn interfaces.Mesh, pos mgl32.Vec3) {
	msh, _ := b.GetMeshByIndex(0)
	m := msh.(*mesh.ColorMesh)
	m.SetParent(scrn)
	m.SetPosition(pos)
}

// Hover changes the color of the surface to the hoverColor.
func (b *Button) Hover() {
	b.baseModelToHoverState()
}

// Clear changes the color of the surface to the defaultColor.
func (b *Button) Clear() {
	b.baseModelToDefaultState()
}
func (b *Button) baseModelToHoverState() {
	bgRect := rectangle.NewExact(b.frameSize.Y(), b.frameSize.X())
	V, I, BO := bgRect.ColoredMeshInput(b.defaultColor)
	bg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	bg.SetBoundingObject(BO)
	fgRect := rectangle.NewExact(b.surfaceSize.Y(), b.surfaceSize.X())
	V, I, _ = fgRect.ColoredMeshInput(b.hoverColor)
	fg := mesh.NewColorMesh(V, I, b.hoverColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, 0.001, 0.0})
	fg.SetParent(bg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	b.BaseModel = m
	b.PinToScreen(b.screen, b.position)
}
func (b *Button) baseModelToDefaultState() {
	bgRect := rectangle.NewExact(b.frameSize.Y(), b.frameSize.X())
	V, I, BO := bgRect.ColoredMeshInput(b.defaultColor)
	bg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	bg.SetBoundingObject(BO)
	fgRect := rectangle.NewExact(b.surfaceSize.Y(), b.surfaceSize.X())
	V, I, _ = fgRect.ColoredMeshInput(b.defaultColor)
	fg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, 0.001, 0.0})
	fg.SetParent(bg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	b.BaseModel = m
	b.PinToScreen(b.screen, b.position)
}

// NewButton returns a button instance. The following inputs has to be set:
// Size of the frame mesh, size of the surface mesh, default and hover color of the button.
// The size (vec2) inputs, the x component means the length on the horizontal axis,
// the y component means the length on the vertical axis.
func NewButton(sizeFrame, sizeSurface mgl32.Vec2, defaultCol, hoverCol []mgl32.Vec3, scrn interfaces.Mesh, pos mgl32.Vec3) *Button {
	btn := &Button{
		BaseModel:    model.New(),
		defaultColor: defaultCol,
		hoverColor:   hoverCol,
		frameSize:    sizeFrame,
		surfaceSize:  sizeSurface,
		screen:       scrn,
		position:     pos,
	}
	btn.baseModelToDefaultState()
	return btn
}

// It represents our editor.
type EditorScreen struct {
	*screen.Screen
	menuShader *shader.Shader
	menuModels []interfaces.Model
}

func NewEditorScreen() *EditorScreen {
	scrn := screen.New()
	scrn.SetupCamera(CreateCamera(), CameraMovementOptions())
	shaderProgram := shader.NewMaterialShader(glWrapper)
	scrn.AddShader(shaderProgram)
	ModelSphere := model.New()
	ModelSphere.AddMesh(CreateJadeSphere())
	scrn.AddModelToShader(ModelSphere, shaderProgram)
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
		mgl32.Vec3{1.0, 1.0, 1.0},
	})
	// Add the lightources to the application
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.Setup(setupApp)
	ModelMenu := model.New()
	var MenuModels []interfaces.Model
	screenMesh := CreateMenuRectangle()
	ModelMenu.AddMesh(screenMesh)
	MenuModels = append(MenuModels, ModelMenu)
	btn := NewButton(frame, surface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{-0.01, 0.7, 0.35})
	MenuModels = append(MenuModels, btn)
	return &EditorScreen{
		Screen:     scrn,
		menuShader: shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper),
		menuModels: MenuModels,
	}
}

// AddMenuPanel activates the menu form on the screen.
func (scrn *EditorScreen) AddMenuPanel() {
	for index, _ := range scrn.menuModels {
		scrn.AddModelToShader(scrn.menuModels[index], scrn.menuShader)
	}
}

// RemoveMenuPanel removes the menu form from the screen.
func (scrn *EditorScreen) RemoveMenuPanel() {
	for index, _ := range scrn.menuModels {
		scrn.RemoveModelFromShader(scrn.menuModels[index], scrn.menuShader)
	}
}

// MenuItemsDefaultState rebuilds the menu models in their default state.
func (scrn *EditorScreen) MenuItemsDefaultState() {
	for index, _ := range scrn.menuModels {
		item := scrn.menuModels[index].(*Button)
		item.Clear()
	}
}
func init() {
	// lock thread
	runtime.LockOSThread()
	app = application.New(glWrapper)
	// Setup the window and the WindowBuilder
	app.SetWindow(setupWindowBuilder().Build())
	glWrapper.InitOpenGL()
	// application screen
	AppScreen = CreateApplicationScreen()
	app.AddScreen(AppScreen)
	app.ActivateScreen(AppScreen)
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	lastUpdate = time.Now().UnixNano() / int64(time.Millisecond)
	lastToggle = lastUpdate
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
	JadeSphere.SetPosition(mgl32.Vec3{5.0, 0.0, 1.4142})
	return JadeSphere
}
func CreateMenuRectangle() *mesh.ColorMesh {
	rect := rectangle.NewExact(2.0, 1.0)
	colors := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 1.0}}
	v, i, _ := rect.ColoredMeshInput(colors)
	menu := mesh.NewColorMesh(v, i, colors, glWrapper)
	menu.SetPosition(mgl32.Vec3{1.4142, 0.0, -0.5})
	menu.RotateZ(90)
	return menu
}
func CreateApplicationScreen() *EditorScreen {
	scrn := NewEditorScreen()
	if MenuScreenEnabled {
		scrn.AddMenuPanel()
	}
	return scrn
}
func Update() {
	current := time.Now().UnixNano() / int64(time.Millisecond)
	delta := current - lastUpdate
	lastUpdate = current
	app.Update(float64(delta))
	if app.GetKeyState(glfw.KeyS) && lastUpdate-lastToggle > 200 {
		fmt.Printf("Key pressed. LastUpdate: '%d', lastToggle: '%d', diff: %d\n", lastUpdate, lastToggle, lastUpdate-lastToggle)
		MenuScreenEnabled = !MenuScreenEnabled
		lastToggle = lastUpdate
		if MenuScreenEnabled {
			AppScreen.AddMenuPanel()
		} else {
			AppScreen.RemoveMenuPanel()
		}
	}
	mod, _, dist := AppScreen.GetClosestModelMeshDistance()
	switch mod.(type) {
	case (*Button):
		if dist < 0.1 {
			btn := mod.(*Button)
			btn.Hover()
		}
		fmt.Printf("Distance: %f\n", dist)
		break
	}
}

func main() {
	defer glfw.Terminate()

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
