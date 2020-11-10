package main

import (
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

type Label struct {
	text     string     // This value is printed to the surface mesh
	color    mgl32.Vec3 // The value will be printed with this color.
	position mgl32.Vec3 // The position of the text (relative from the surface).
	size     float32    // The size of the text
	surface  interfaces.Mesh
}

// NewLabel returns the label instance
func NewLabel(text string, color, position mgl32.Vec3, size float32, surface interfaces.Mesh) *Label {
	return &Label{
		text:     text,
		color:    color,
		position: position,
		size:     size,
		surface:  surface,
	}
}

// GetLabelText returns the text of the label.
func (l *Label) GetLabelText() string {
	return l.text
}

// GetLabelColor returns the color of the label.
func (l *Label) GetLabelColor() mgl32.Vec3 {
	return l.color
}

// GetLabelPosition returns the position of the label.
func (l *Label) GetLabelPosition() mgl32.Vec3 {
	return l.position
}

// GetLabelSize returns the size of the label.
func (l *Label) GetLabelSize() float32 {
	return l.size
}

// GetLabelSurface returns the surface of the label.
func (l *Label) GetLabelSurface() interfaces.Mesh {
	return l.surface
}

// SetLabelSurface updatess the surface of the label.
func (l *Label) SetLabelSurface(surface interfaces.Mesh) {
	l.surface = surface
}

/*
 * Now i want to extend the application with the button model. It supposed to contain 2 meshes (bg, surface), should have the update options.
 * The menuModel should be extended with the menu items. It could be an array of models with some support functions.
 */

// It is the representation of the button ui item.
type Button struct {
	*model.BaseModel
	*Label
	defaultColor   []mgl32.Vec3
	hoverColor     []mgl32.Vec3
	frameSize      mgl32.Vec2
	surfaceSize    mgl32.Vec2
	screen         interfaces.Mesh
	positionOnForm mgl32.Vec3
	aspect         float32
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
	bgRect := rectangle.NewExact(b.frameSize.Y()/b.aspect, b.frameSize.X()/b.aspect)
	V, I, BO := bgRect.ColoredMeshInput(b.hoverColor)
	bg := mesh.NewColorMesh(V, I, b.hoverColor, glWrapper)
	bg.SetBoundingObject(BO)
	fgRect := rectangle.NewExact(b.surfaceSize.Y()/b.aspect, b.surfaceSize.X()/b.aspect)
	V, I, _ = fgRect.ColoredMeshInput(b.defaultColor)
	fg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -0.001, 0.0})
	fg.SetParent(bg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	b.BaseModel = m
	pos := mgl32.Vec3{b.positionOnForm.X() / b.aspect, b.positionOnForm.Y(), b.positionOnForm.Z() / b.aspect}
	b.PinToScreen(b.screen, pos)
	if b.HasLabel() {
		b.SetLabelSurface(fg)
	}
}
func (b *Button) baseModelToDefaultState() {
	bgRect := rectangle.NewExact(b.frameSize.Y()/b.aspect, b.frameSize.X()/b.aspect)
	V, I, BO := bgRect.ColoredMeshInput(b.defaultColor)
	bg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	bg.SetBoundingObject(BO)
	fgRect := rectangle.NewExact(b.surfaceSize.Y()/b.aspect, b.surfaceSize.X()/b.aspect)
	V, I, _ = fgRect.ColoredMeshInput(b.defaultColor)
	fg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -0.001, 0.0})
	fg.SetParent(bg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	b.BaseModel = m
	pos := mgl32.Vec3{b.positionOnForm.X() / b.aspect, b.positionOnForm.Y(), b.positionOnForm.Z() / b.aspect}
	b.PinToScreen(b.screen, pos)
	if b.HasLabel() {
		b.SetLabelSurface(fg)
	}
}
func (b *Button) SetAspect(aspect float32) {
	b.aspect = aspect
}
func (b *Button) SetLabel(l *Label) {
	b.Label = l
}
func (b *Button) HasLabel() bool {
	return b.Label != nil
}

// NewButton returns a button instance. The following inputs has to be set:
// Size of the frame mesh, size of the surface mesh, default and hover color of the button.
// The size (vec2) inputs, the x component means the length on the horizontal axis,
// the y component means the length on the vertical axis.
func NewButton(sizeFrame, sizeSurface mgl32.Vec2, defaultCol, hoverCol []mgl32.Vec3, scrn interfaces.Mesh, pos mgl32.Vec3, aspect float32) *Button {
	btn := &Button{
		BaseModel:      model.New(),
		Label:          nil,
		defaultColor:   defaultCol,
		hoverColor:     hoverCol,
		frameSize:      sizeFrame,
		surfaceSize:    sizeSurface,
		screen:         scrn,
		positionOnForm: pos,
		aspect:         aspect,
	}
	btn.baseModelToDefaultState()
	return btn
}

/*
type FormSurfaceModel struct {
	*model.BaseModel
	surfaceColor []mgl32.Vec3
	surfaceSize  mgl32.Vec2
	position     mgl32.Vec3
	aspect       float32
}
*/

// It represents our editor.
type EditorScreen struct {
	*screen.Screen
	menuShader *shader.Shader
	menuModels []interfaces.Model
	charset    *model.Charset
	screenMesh *mesh.ColorMesh
}

func NewEditorScreen() *EditorScreen {
	scrn := screen.New()
	scrn.SetWrapper(glWrapper)
	wX, wY := app.GetWindow().GetSize()
	scrn.SetWindowSize(float32(wX), float32(wY))
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
	ModelMenu := model.New()
	var MenuModels []interfaces.Model
	aspectRatio := scrn.GetAspectRatio()
	screenMesh := CreateMenuRectangle(aspectRatio)
	ModelMenu.AddMesh(screenMesh)
	MenuModels = append(MenuModels, ModelMenu)
	btn := NewButton(frame, surface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -0.01, -0.35}, aspectRatio)
	s, _ := btn.GetMeshByIndex(1)
	btn.SetLabel(NewLabel("Test", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, 0}, 3/800, s))
	MenuModels = append(MenuModels, btn)
	es := &EditorScreen{
		Screen:     scrn,
		menuShader: shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper),
		menuModels: MenuModels,
		charset:    nil,
		screenMesh: screenMesh,
	}
	es.AddShader(es.menuShader)
	es.Setup(es.setupApp)
	es.defaultCharset()
	// font shader
	fontShader := shader.NewShader(baseDir()+"/shaders/font.vert", baseDir()+"/shaders/font.frag", es.GetWrapper())
	es.AddShader(fontShader)
	es.AddModelToShader(es.charset, fontShader)
	// Try to write something to the surface of the blue form container.
	es.charset.PrintTo("Something", 0.0, 0.0, 0.01, 3/800, scrn.GetWrapper(), es.screenMesh, []mgl32.Vec3{mgl32.Vec3{1, 0, 1}})
	return es
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
		switch scrn.menuModels[index].(type) {
		case *Button:
			item := scrn.menuModels[index].(*Button)
			item.Clear()
			break
		}
	}
}
func (scrn *EditorScreen) Update(dt float64, p interfaces.Pointer, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	posX, posY := p.GetCurrent()
	mCoords := mgl32.Vec3{float32(-posY) / scrn.GetAspectRatio(), -0.01, float32(posX)}
	scrn.UpdateWithDistance(dt, mCoords)
	closestModel, _, dist := scrn.GetClosestModelMeshDistance()
	scrn.MenuItemsDefaultState()
	switch closestModel.(type) {
	case (*Button):
		if dist < 0.001 {
			btn := closestModel.(*Button)
			btn.Hover()
		}
		break
	}
	for index, _ := range scrn.menuModels {
		switch scrn.menuModels[index].(type) {
		case *Button:
			item := scrn.menuModels[index].(*Button)
			if item.HasLabel() {
				pos := item.GetLabelPosition()
				scrn.charset.PrintTo(item.GetLabelText(), pos.X(), pos.Y(), pos.Z(), item.GetLabelSize(), scrn.GetWrapper(), item.GetLabelSurface(), []mgl32.Vec3{item.GetLabelColor()})
			}
			break
		}
	}
}
func (scrn *EditorScreen) defaultCharset() {
	cs, err := model.LoadCharset("assets/fonts/Desyrel/desyrel.regular.ttf", 32, 127, 40.0, 72, scrn.GetWrapper())
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	scrn.charset = cs
}
func (scrn *EditorScreen) setupApp(w interfaces.GLWrapper) {
	scrn.GetWrapper().Enable(glwrapper.DEPTH_TEST)
	scrn.GetWrapper().DepthFunc(glwrapper.LESS)
	scrn.GetWrapper().ClearColor(0.3, 0.3, 0.3, 1.0)
	scrn.GetWrapper().Enable(glwrapper.BLEND)
	scrn.GetWrapper().BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	wW, wH := scrn.GetWindowSize()
	scrn.GetWrapper().Viewport(0, 0, int32(wW), int32(wH))
}
func (scrn *EditorScreen) ResizeEvent(wW, wH float32) {
	scrn.SetWindowSize(wW, wH)
	for index, _ := range scrn.menuModels {
		switch scrn.menuModels[index].(type) {
		case *Button:
			item := scrn.menuModels[index].(*Button)
			item.SetAspect(wW / wH)
			break
		}
	}
}

// type SizeCallback func(w *Window, width int, height int)
func ResizeCallback(w *glfw.Window, width, height int) {
	AppScreen.ResizeEvent(float32(width), float32(height))
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
	app.GetWindow().SetSizeCallback(ResizeCallback)
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

// It creates a new camera with the necessary setup
// cameraPosition: Vector{X : 0.0000000000, Y : -1.7319999933, Z : 0.0000000000}
// worldUp: Vector{X : 1.0000000000, Y : 0.0000000000, Z : 0.0000000000}
// cameraFrontDirection: Vector{X : -0.0000000437, Y : 1.0000000000, Z : -0.0000000000}
// cameraUpDirection: Vector{X : -1.0000000000, Y : -0.0000000437, Z : 0.0000000000}
// cameraRightDirection: Vector{X : -0.0000000000, Y : 0.0000000000, Z : 1.0000000000}
// yaw : 0.0000000000
// pitch : 90.0000000000
// velocity : 0.0049999999
// rotationStep : 0.0000000000
// ProjectionOptions:
//  - fov : 45.0000000000
//  - aspectRatio : 1.0000000000
//  - far : 10.0000000000
//  - near : 0.0010000000
func CreateCamera() *camera.DefaultCamera {
	mat := mgl32.Perspective(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 10.0)
	// 1.732 ~ sqrt(3)
	camera := camera.NewCamera(mgl32.Vec3{0.0, -mat[0], 0.0}, mgl32.Vec3{1, 0, 0}, 0.0, 90.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.0001, 10.0)
	camera.SetVelocity(float32(0.005))
	return camera
}

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["mode"] = "default"
	cm["rotateOnEdgeDistance"] = float32(0.0)
	cm["forward"] = []glfw.Key{glfw.KeyY}
	cm["back"] = []glfw.Key{glfw.KeyX}
	return cm
}

// It generates the Jade sphere.
func CreateJadeSphere() *mesh.MaterialMesh {
	sph := sphere.New(20)
	v, i, _ := sph.MaterialMeshInput()
	JadeSphere := mesh.NewMaterialMesh(v, i, material.Jade, glWrapper)
	JadeSphere.SetPosition(mgl32.Vec3{0.0, 3.5858, -1.4142})
	return JadeSphere
}
func CreateMenuRectangle(aspect float32) *mesh.ColorMesh {
	rect := rectangle.NewExact(2.0/aspect, 1.0/aspect)
	colors := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 1.0}}
	v, i, _ := rect.ColoredMeshInput(colors)
	menu := mesh.NewColorMesh(v, i, colors, glWrapper)
	menu.SetPosition(mgl32.Vec3{0.0, 0.0, 0.5 / aspect})
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
		MenuScreenEnabled = !MenuScreenEnabled
		lastToggle = lastUpdate
		if MenuScreenEnabled {
			AppScreen.AddMenuPanel()
		} else {
			AppScreen.RemoveMenuPanel()
		}
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
