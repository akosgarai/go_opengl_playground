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
	"github.com/akosgarai/playground_engine/pkg/transformations"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowTitle                      = "Example - real-time editor"
	ForegroundDistanceFromBackground = float32(0.002)
	InputFieldDistanceFromForeground = float32(0.003)
	LEFT_MOUSE_BUTTON                = glfw.MouseButtonLeft
	FormItemsDistanceFromScreen      = float32(0.01)
	CollisionEpsilon                 = float32(0.001)
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
	Buttonframe        = mgl32.Vec2{0.1, 0.2}
	Buttonsurface      = mgl32.Vec2{0.0925, 0.185}
	buttonDefaultColor = []mgl32.Vec3{mgl32.Vec3{0.4, 0.4, 0.4}}
	buttonHoverColor   = []mgl32.Vec3{mgl32.Vec3{0.4, 0.8, 0.4}}
	// Text Input variables
	TextInputframe        = mgl32.Vec2{0.2, 0.8}
	TextInputsurface      = mgl32.Vec2{0.1925, 0.785}
	TextInputField        = mgl32.Vec2{0.09625, 0.785}
	TextInputDefaultColor = []mgl32.Vec3{mgl32.Vec3{0.4, 0.4, 0.4}}
	TextInputHoverColor   = []mgl32.Vec3{mgl32.Vec3{0.4, 0.8, 0.4}}
	TextInputFieldColor   = []mgl32.Vec3{mgl32.Vec3{1.0, 1.0, 1.0}}
	// the states
	AllScreenStates = []string{"Default", "Material", "MaterialAmbientForm", "MaterialDiffuseForm", "MaterialSpecularForm", "MaterialShininessForm"}
	ScreenLabels    = map[string]string{
		"Default":               "Editor application",
		"Material":              "Material",
		"MaterialAmbientForm":   "Material > Ambient",
		"MaterialDiffuseForm":   "Material > Diffuse",
		"MaterialSpecularForm":  "Material > Specualar",
		"MaterialShininessForm": "Material > Shininess",
	}
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

// It is the representation of the slider input ui item.
// The idea about the item: Based on rectangles.
// The background rectangle is responsible for the hover event
// The foreground rectangle is split (horizontal) two half.
// The top half contains the label of the item.
// The bottom half contains a slip bar and also a label
// for the value of the slip bar.
// The ratio between she slider and the label is 3/1
type SliderInput struct {
	*model.BaseModel
	*Label
	defaultColor   []mgl32.Vec3
	hoverColor     []mgl32.Vec3
	frameSize      mgl32.Vec2
	surfaceSize    mgl32.Vec2
	screen         interfaces.Mesh
	positionOnForm mgl32.Vec3
	aspect         float32
	textInputColor []mgl32.Vec3 // The read-only surface color for the value of the slider.
	sliderMin      float32
	sliderMax      float32
	sliderCurrent  float32 // This is the current value of the slider input.
}

// NewSliderInput returns a slider input instance. The following inputs has to be set:
// Size of the frame mesh, size of the surface mesh,
// default and hover color of the button, color of the text field,
// and the min,max values of the slider.
// The size (vec2) inputs, the x component means the length on the horizontal axis,
// the y component means the length on the vertical axis.
func NewSliderInput(sizeFrame, sizeSurface mgl32.Vec2, defaultCol, hoverCol, tiCol []mgl32.Vec3, scrn interfaces.Mesh, pos mgl32.Vec3, aspect, min, max float32) *SliderInput {
	si := &SliderInput{
		BaseModel:      model.New(),
		Label:          nil,
		defaultColor:   defaultCol,
		hoverColor:     hoverCol,
		frameSize:      sizeFrame,
		surfaceSize:    sizeSurface,
		textInputColor: tiCol,
		screen:         scrn,
		positionOnForm: pos,
		aspect:         aspect,
		sliderMin:      min,
		sliderMax:      max,
		sliderCurrent:  (max-min)/2.0 + min,
	}
	si.baseModelToDefaultState()
	return si
}
func (si *SliderInput) SetAspect(aspect float32) {
	si.aspect = aspect
}
func (si *SliderInput) SetLabel(l *Label) {
	si.Label = l
}
func (si *SliderInput) HasLabel() bool {
	return si.Label != nil
}

// PinToScreen sets the parent of the bg mesh to the given one and updates its position.
func (si *SliderInput) PinToScreen(scrn interfaces.Mesh, pos mgl32.Vec3) {
	msh, _ := si.GetMeshByIndex(0)
	m := msh.(*mesh.ColorMesh)
	m.SetParent(scrn)
	m.SetPosition(pos)
}

// Hover changes the color of the surface to the hoverColor.
func (si *SliderInput) Hover() {
	si.baseModelToHoverState()
}

// Clear changes the color of the surface to the defaultColor.
func (si *SliderInput) Clear() {
	si.baseModelToDefaultState()
}

// The black line for the slider. It represents the min-max interval.
// The length of it: 3/4 of the foreground width - width of the slider.
// The width of the slider is interval length / 10
func (si *SliderInput) sliderAreaWidth() float32 {
	return 3 * si.surfaceSize.Y() / 4 / si.aspect
}
func (si *SliderInput) sliderSquareWidth() float32 {
	return si.sliderAreaWidth() / 10.0
}
func (si *SliderInput) sliderSpaceWidth() float32 {
	return si.sliderAreaWidth() * 0.9
}
func (si *SliderInput) baseModelToHoverState() {
	si.baseModelToState("hover")
}
func (si *SliderInput) baseModelToDefaultState() {
	si.baseModelToState("default")
}
func (si *SliderInput) baseModelToState(state string) {
	var bgColor []mgl32.Vec3
	switch state {
	case "default":
		bgColor = si.defaultColor
		break
	case "hover":
		bgColor = si.hoverColor
		break
	default:
		panic(fmt.Sprintf("Unknown state name: %s.", state))
	}
	bgRect := rectangle.NewExact(si.frameSize.Y()/si.aspect, si.frameSize.X()/si.aspect)
	V, I, BO := bgRect.ColoredMeshInput(bgColor)
	bg := mesh.NewColorMesh(V, I, bgColor, glWrapper)
	bg.SetBoundingObject(BO)
	bg.RotateY(-90)
	fgRect := rectangle.NewExact(si.surfaceSize.Y()/si.aspect, si.surfaceSize.X()/si.aspect)
	V, I, _ = fgRect.ColoredMeshInput(si.defaultColor)
	fg := mesh.NewColorMesh(V, I, si.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, 0.0})
	fg.SetParent(bg)
	// text input field
	// it has to be on the bottom half of the screen -> height: fg.height/2.
	// width: fg.width/4.
	tiRect := rectangle.NewExact(si.surfaceSize.Y()/4/si.aspect, si.surfaceSize.X()/2/si.aspect)
	V, I, _ = tiRect.ColoredMeshInput(si.textInputColor)
	tif := mesh.NewColorMesh(V, I, si.textInputColor, glWrapper)
	tif.SetPosition(mgl32.Vec3{si.surfaceSize.X() / 4.0 / si.aspect, -InputFieldDistanceFromForeground, si.sliderAreaWidth() / 2.0})
	tif.SetParent(fg)
	blackLineColor := []mgl32.Vec3{mgl32.Vec3{0, 0, 0}}
	blackLineRect := rectangle.NewExact(si.sliderSpaceWidth(), 0.005/si.aspect)
	V, I, _ = blackLineRect.ColoredMeshInput(blackLineColor)
	blackLine := mesh.NewColorMesh(V, I, blackLineColor, glWrapper)
	blackLine.SetParent(fg)
	blackLine.SetPosition(mgl32.Vec3{si.surfaceSize.X() / 4.0 / si.aspect, -InputFieldDistanceFromForeground, -si.surfaceSize.Y() / 8 / si.aspect})
	sliderRect := rectangle.NewExact(si.sliderSquareWidth(), si.sliderSquareWidth())
	sliderColor := []mgl32.Vec3{mgl32.Vec3{0.3, 0.5, 0.7}}
	V, I, BO = sliderRect.ColoredMeshInput(sliderColor)
	slider := mesh.NewColorMesh(V, I, sliderColor, glWrapper)
	slider.SetBoundingObject(BO)
	slider.SetParent(blackLine)
	// The slider is a scrollable item, so that we need to redraw it with its current position.
	origSlider, err := si.GetMeshByIndex(4)
	if err != nil {
		slider.SetPosition(si.calculatePositionBasedOnCurrent())
	} else {
		slider.SetPosition(origSlider.GetPosition())
	}
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	m.AddMesh(tif)
	m.AddMesh(blackLine)
	m.AddMesh(slider)
	si.BaseModel = m
	pos := mgl32.Vec3{si.positionOnForm.X() / si.aspect, si.positionOnForm.Y(), si.positionOnForm.Z() / si.aspect}
	si.PinToScreen(si.screen, pos)
	if si.HasLabel() {
		si.SetLabelSurface(fg)
	}
}
func (si *SliderInput) SliderCollision(coords mgl32.Vec3) bool {
	slider, err := si.GetMeshByIndex(4)
	if err != nil {
		fmt.Printf("Slider issue. %#v\n", err)
		return false
	}
	coordsInSliderPlane := mgl32.Vec3{coords.X(), coords.Y() - (2 * ForegroundDistanceFromBackground) - InputFieldDistanceFromForeground, coords.Z()}
	msh, dist := si.ClosestMeshTo(coordsInSliderPlane)
	if dist > CollisionEpsilon {
		return false
	}
	return msh == slider
}
func (si *SliderInput) SetCurrentValue(current float32) {
	si.sliderCurrent = current
	msh, err := si.GetMeshByIndex(4)
	if err != nil {
		fmt.Println("Mesh is missing. Skipping position update.")
		return
	}
	msh.SetPosition(si.calculatePositionBasedOnCurrent())

}
func (si *SliderInput) MoveSliderWith(x float32) {
	slider, err := si.GetMeshByIndex(4)
	if err != nil {
		return
	}
	currentPosition := slider.GetPosition()
	// The vertical position has to be kept above the black line, so that it might be updated with the max/min value.
	sliderSpaceWidth := si.sliderSpaceWidth()
	newVerticalCoordinateValue := currentPosition.Z() - x
	if newVerticalCoordinateValue > sliderSpaceWidth/2 {
		newVerticalCoordinateValue = sliderSpaceWidth / 2
	}
	if newVerticalCoordinateValue < -sliderSpaceWidth/2 {
		newVerticalCoordinateValue = -sliderSpaceWidth / 2
	}
	newPosition := mgl32.Vec3{currentPosition.X(), currentPosition.Y(), newVerticalCoordinateValue}
	slider.SetPosition(newPosition)
	si.updateCurrentFromPosition(newPosition)
}
func (si *SliderInput) calculatePositionBasedOnCurrent() mgl32.Vec3 {
	intervalLength := si.sliderMax - si.sliderMin
	ratio := si.sliderCurrent / intervalLength
	sliderSpaceWidth := si.sliderSpaceWidth()
	diffFromLeft := sliderSpaceWidth * ratio
	return mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, -sliderSpaceWidth/2 + diffFromLeft}
}
func (si *SliderInput) updateCurrentFromPosition(currentPosition mgl32.Vec3) {
	sliderSpaceWidth := si.sliderSpaceWidth()
	diffFromMinimum := currentPosition.Z() + sliderSpaceWidth/2
	ratio := diffFromMinimum / sliderSpaceWidth
	value := (si.sliderMax-si.sliderMin)*ratio + si.sliderMin
	si.sliderCurrent = value
}

// It is the representation of the text input ui item.
// The idea about the item: Based on rectangles.
// The background rectangle is responsible for the hover event
// The foreground rectangle is split (horizontal) two half.
// The top half contains the label of the item.
// The bottom half contains the current value of the input.
type TextInput struct {
	*model.BaseModel
	*Label
	defaultColor   []mgl32.Vec3
	hoverColor     []mgl32.Vec3
	frameSize      mgl32.Vec2
	surfaceSize    mgl32.Vec2
	textInputSize  mgl32.Vec2
	textInputColor []mgl32.Vec3
	screen         interfaces.Mesh
	positionOnForm mgl32.Vec3
	aspect         float32
}

// NewTextInput returns a text input instance. The following inputs has to be set:
// Size of the frame mesh, size of the surface mesh, size of the textInput,
// default and hover color of the button, color of the text input field.
// The size (vec2) inputs, the x component means the length on the horizontal axis,
// the y component means the length on the vertical axis.
func NewTextInput(sizeFrame, sizeSurface, textInputSize mgl32.Vec2, defaultCol, hoverCol, tiCol []mgl32.Vec3, scrn interfaces.Mesh, pos mgl32.Vec3, aspect float32) *TextInput {
	ti := &TextInput{
		BaseModel:      model.New(),
		Label:          nil,
		defaultColor:   defaultCol,
		hoverColor:     hoverCol,
		frameSize:      sizeFrame,
		surfaceSize:    sizeSurface,
		textInputSize:  textInputSize,
		textInputColor: tiCol,
		screen:         scrn,
		positionOnForm: pos,
		aspect:         aspect,
	}
	ti.baseModelToDefaultState()
	return ti
}
func (ti *TextInput) SetAspect(aspect float32) {
	ti.aspect = aspect
}
func (ti *TextInput) SetLabel(l *Label) {
	ti.Label = l
}
func (ti *TextInput) HasLabel() bool {
	return ti.Label != nil
}

// PinToScreen sets the parent of the bg mesh to the given one and updates its position.
func (ti *TextInput) PinToScreen(scrn interfaces.Mesh, pos mgl32.Vec3) {
	msh, _ := ti.GetMeshByIndex(0)
	m := msh.(*mesh.ColorMesh)
	m.SetParent(scrn)
	m.SetPosition(pos)
}

// Hover changes the color of the surface to the hoverColor.
func (ti *TextInput) Hover() {
	ti.baseModelToHoverState()
}

// Clear changes the color of the surface to the defaultColor.
func (ti *TextInput) Clear() {
	ti.baseModelToDefaultState()
}
func (ti *TextInput) baseModelToHoverState() {
	bgRect := rectangle.NewExact(ti.frameSize.Y()/ti.aspect, ti.frameSize.X()/ti.aspect)
	V, I, BO := bgRect.ColoredMeshInput(ti.hoverColor)
	bg := mesh.NewColorMesh(V, I, ti.hoverColor, glWrapper)
	bg.SetBoundingObject(BO)
	bg.RotateY(-90)
	fgRect := rectangle.NewExact(ti.surfaceSize.Y()/ti.aspect, ti.surfaceSize.X()/ti.aspect)
	V, I, _ = fgRect.ColoredMeshInput(ti.defaultColor)
	fg := mesh.NewColorMesh(V, I, ti.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, 0.0})
	fg.SetParent(bg)
	// text input field
	tiRect := rectangle.NewExact(ti.textInputSize.Y()/ti.aspect, ti.textInputSize.X()/ti.aspect)
	V, I, _ = tiRect.ColoredMeshInput(ti.textInputColor)
	tif := mesh.NewColorMesh(V, I, ti.textInputColor, glWrapper)
	tif.SetPosition(mgl32.Vec3{ti.textInputSize.X() / 2.0 / ti.aspect, -InputFieldDistanceFromForeground, 0.0})
	tif.SetParent(fg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	m.AddMesh(tif)
	ti.BaseModel = m
	pos := mgl32.Vec3{ti.positionOnForm.X() / ti.aspect, ti.positionOnForm.Y(), ti.positionOnForm.Z() / ti.aspect}
	ti.PinToScreen(ti.screen, pos)
	if ti.HasLabel() {
		ti.SetLabelSurface(fg)
	}
}
func (ti *TextInput) baseModelToDefaultState() {
	bgRect := rectangle.NewExact(ti.frameSize.Y()/ti.aspect, ti.frameSize.X()/ti.aspect)
	V, I, BO := bgRect.ColoredMeshInput(ti.defaultColor)
	bg := mesh.NewColorMesh(V, I, ti.defaultColor, glWrapper)
	bg.SetBoundingObject(BO)
	bg.RotateY(-90)
	fgRect := rectangle.NewExact(ti.surfaceSize.Y()/ti.aspect, ti.surfaceSize.X()/ti.aspect)
	V, I, _ = fgRect.ColoredMeshInput(ti.defaultColor)
	fg := mesh.NewColorMesh(V, I, ti.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, 0.0})
	fg.SetParent(bg)
	// text input field
	tiRect := rectangle.NewExact(ti.textInputSize.Y()/ti.aspect, ti.textInputSize.X()/ti.aspect)
	V, I, _ = tiRect.ColoredMeshInput(ti.textInputColor)
	tif := mesh.NewColorMesh(V, I, ti.textInputColor, glWrapper)
	tif.SetPosition(mgl32.Vec3{ti.textInputSize.X() / 2.0 / ti.aspect, -InputFieldDistanceFromForeground, 0.0})
	tif.SetParent(fg)
	m := model.New()
	m.AddMesh(bg)
	m.AddMesh(fg)
	m.AddMesh(tif)
	ti.BaseModel = m
	pos := mgl32.Vec3{ti.positionOnForm.X() / ti.aspect, ti.positionOnForm.Y(), ti.positionOnForm.Z() / ti.aspect}
	ti.PinToScreen(ti.screen, pos)
	if ti.HasLabel() {
		ti.SetLabelSurface(fg)
	}
}

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
	clicked        bool
	clickCallback  func()
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
	bg.RotateY(-90)
	fgRect := rectangle.NewExact(b.surfaceSize.Y()/b.aspect, b.surfaceSize.X()/b.aspect)
	V, I, _ = fgRect.ColoredMeshInput(b.defaultColor)
	fg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, 0.0})
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
	bg.RotateY(-90)
	fgRect := rectangle.NewExact(b.surfaceSize.Y()/b.aspect, b.surfaceSize.X()/b.aspect)
	V, I, _ = fgRect.ColoredMeshInput(b.defaultColor)
	fg := mesh.NewColorMesh(V, I, b.defaultColor, glWrapper)
	fg.SetPosition(mgl32.Vec3{0.0, -ForegroundDistanceFromBackground, 0.0})
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
		clicked:        false,
	}
	btn.baseModelToDefaultState()
	return btn
}

// It represents our editor.
type EditorScreen struct {
	*screen.Screen
	menuShader *shader.Shader
	menuModels map[string][]interfaces.Model
	charset    *model.Charset
	// We draw the UI items based on the state.
	// Default state - only the Material button is shown.
	// Material state - the Color and the back buttons are shown. The Material
	// text is printed to the top of the screen mesh.
	// Material[Amibent|Diffuse|Specular|Shininess]Form state - the back button is printed,
	// the Material > Color comp. text is printed to the top of the screen mesh.
	state string
	// the sphere model for easy access.
	sphereModel *model.BaseModel
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
	MenuModels := make(map[string][]interfaces.Model)
	aspectRatio := scrn.GetAspectRatio()
	screenMesh := CreateMenuRectangle(aspectRatio)
	screenLabelMesh := CreateMenuLabelRectangle(aspectRatio)
	screenLabelMesh.SetParent(screenMesh)
	ModelMenu.AddMesh(screenMesh)
	ModelMenu.AddMesh(screenLabelMesh)
	es := &EditorScreen{
		Screen:      scrn,
		menuShader:  shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper),
		charset:     nil,
		state:       "Default",
		sphereModel: ModelSphere,
	}
	MenuModels["Default"] = append(MenuModels["Default"], ModelMenu)
	MenuModels["Material"] = append(MenuModels["Material"], ModelMenu)
	MenuModels["MaterialAmbientForm"] = append(MenuModels["MaterialAmbientForm"], ModelMenu)
	MenuModels["MaterialDiffuseForm"] = append(MenuModels["MaterialDiffuseForm"], ModelMenu)
	MenuModels["MaterialSpecularForm"] = append(MenuModels["MaterialSpecularForm"], ModelMenu)
	MenuModels["MaterialShininessForm"] = append(MenuModels["MaterialShininessForm"], ModelMenu)
	btn := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, -0.35}, aspectRatio)
	s, err := btn.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btn.SetLabel(NewLabel("Material", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btn.clickCallback = es.SetStateMaterial
	MenuModels["Default"] = append(MenuModels["Default"], btn)
	// Ambient button Material State
	btnAmbient := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, -0.35}, aspectRatio)
	s, err = btnAmbient.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnAmbient.SetLabel(NewLabel("Ambient", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnAmbient.clickCallback = es.SetStateMaterialAmbientForm
	MenuModels["Material"] = append(MenuModels["Material"], btnAmbient)
	// Diffuse button Material State
	btnDiffuse := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, -0.15}, aspectRatio)
	s, err = btnDiffuse.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnDiffuse.SetLabel(NewLabel("Diffuse", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnDiffuse.clickCallback = es.SetStateMaterialDiffuseForm
	MenuModels["Material"] = append(MenuModels["Material"], btnDiffuse)
	// Specular button Material State
	btnSpecular := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, 0.05}, aspectRatio)
	s, err = btnSpecular.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnSpecular.SetLabel(NewLabel("Specular", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnSpecular.clickCallback = es.SetStateMaterialSpecularForm
	MenuModels["Material"] = append(MenuModels["Material"], btnSpecular)
	// Shininess button Material State
	btnShininess := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.8, -FormItemsDistanceFromScreen, -0.15}, aspectRatio)
	s, err = btnShininess.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnShininess.SetLabel(NewLabel("Shininess", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnShininess.clickCallback = es.SetStateMaterialShininessForm
	MenuModels["Material"] = append(MenuModels["Material"], btnShininess)
	// back button Material State
	btnBack := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, 0.35}, aspectRatio)
	s, err = btnBack.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnBack.SetLabel(NewLabel("Back", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnBack.clickCallback = es.SetStateDefault
	MenuModels["Material"] = append(MenuModels["Material"], btnBack)
	// back button To Material State from color forms.
	btnBackForm := NewButton(Buttonframe, Buttonsurface, buttonDefaultColor, buttonHoverColor, screenMesh, mgl32.Vec3{0.9, -FormItemsDistanceFromScreen, 0.35}, aspectRatio)
	s, err = btnBackForm.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on btn branch.")
		panic(err)
	}
	btnBackForm.SetLabel(NewLabel("Back", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, 0, -FormItemsDistanceFromScreen}, 0.0005, s))
	btnBackForm.clickCallback = es.SetStateMaterial
	MenuModels["MaterialAmbientForm"] = append(MenuModels["MaterialAmbientForm"], btnBackForm)
	MenuModels["MaterialDiffuseForm"] = append(MenuModels["MaterialDiffuseForm"], btnBackForm)
	MenuModels["MaterialSpecularForm"] = append(MenuModels["MaterialSpecularForm"], btnBackForm)
	MenuModels["MaterialShininessForm"] = append(MenuModels["MaterialShininessForm"], btnBackForm)
	// slider inputs for the color components.
	siRed := NewSliderInput(TextInputframe, TextInputsurface, TextInputDefaultColor, TextInputHoverColor, TextInputFieldColor, screenMesh, mgl32.Vec3{-0.5, -FormItemsDistanceFromScreen, 0.0}, aspectRatio, 0.0, 1.0)
	s, err = siRed.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on si branch.")
		panic(err)
	}
	siRed.SetLabel(NewLabel("Red", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, TextInputField.X() / aspectRatio / 2, -0.01}, 0.0005, s))
	MenuModels["MaterialAmbientForm"] = append(MenuModels["MaterialAmbientForm"], siRed)
	MenuModels["MaterialDiffuseForm"] = append(MenuModels["MaterialDiffuseForm"], siRed)
	MenuModels["MaterialSpecularForm"] = append(MenuModels["MaterialSpecularForm"], siRed)
	siGreen := NewSliderInput(TextInputframe, TextInputsurface, TextInputDefaultColor, TextInputHoverColor, TextInputFieldColor, screenMesh, mgl32.Vec3{-0.2, -FormItemsDistanceFromScreen, 0.0}, aspectRatio, 0.0, 1.0)
	s, err = siGreen.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on si branch.")
		panic(err)
	}
	siGreen.SetLabel(NewLabel("Green", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, TextInputField.X() / aspectRatio / 2, -0.01}, 0.0005, s))
	MenuModels["MaterialAmbientForm"] = append(MenuModels["MaterialAmbientForm"], siGreen)
	MenuModels["MaterialDiffuseForm"] = append(MenuModels["MaterialDiffuseForm"], siGreen)
	MenuModels["MaterialSpecularForm"] = append(MenuModels["MaterialSpecularForm"], siGreen)
	siBlue := NewSliderInput(TextInputframe, TextInputsurface, TextInputDefaultColor, TextInputHoverColor, TextInputFieldColor, screenMesh, mgl32.Vec3{0.1, -FormItemsDistanceFromScreen, 0.0}, aspectRatio, 0.0, 1.0)
	s, err = siBlue.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on si branch.")
		panic(err)
	}
	siBlue.SetLabel(NewLabel("Blue", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, TextInputField.X() / aspectRatio / 2, -0.01}, 0.0005, s))
	MenuModels["MaterialAmbientForm"] = append(MenuModels["MaterialAmbientForm"], siBlue)
	MenuModels["MaterialDiffuseForm"] = append(MenuModels["MaterialDiffuseForm"], siBlue)
	MenuModels["MaterialSpecularForm"] = append(MenuModels["MaterialSpecularForm"], siBlue)
	// shininess component
	siShininess := NewSliderInput(TextInputframe, TextInputsurface, TextInputDefaultColor, TextInputHoverColor, TextInputFieldColor, screenMesh, mgl32.Vec3{0.1, -FormItemsDistanceFromScreen, 0.0}, aspectRatio, 0.0, 128.0)
	s, err = siShininess.GetMeshByIndex(1)
	if err != nil {
		fmt.Println("Something terrible happened on si branch.")
		panic(err)
	}
	siShininess.SetLabel(NewLabel("Shininess", mgl32.Vec3{0, 0, 0.05}, mgl32.Vec3{0, TextInputField.X() / aspectRatio / 2, -0.01}, 0.0005, s))
	MenuModels["MaterialShininessForm"] = append(MenuModels["MaterialShininessForm"], siShininess)
	es.menuModels = MenuModels
	es.AddShader(es.menuShader)
	es.Setup(es.setupApp)
	es.defaultCharset()
	// font shader
	fontShader := shader.NewShader(baseDir()+"/shaders/font.vert", baseDir()+"/shaders/font.frag", es.GetWrapper())
	es.AddShader(fontShader)
	es.AddModelToShader(es.charset, fontShader)
	return es
}

// AddMenuPanel activates the menu form on the screen.
func (scrn *EditorScreen) AddMenuPanel() {
	for index, _ := range scrn.menuModels[scrn.state] {
		scrn.AddModelToShader(scrn.menuModels[scrn.state][index], scrn.menuShader)
	}
}
func (scrn *EditorScreen) SetStateDefault() {
	scrn.setState("Default")
}
func (scrn *EditorScreen) SetStateMaterial() {
	scrn.setState("Material")
}
func (scrn *EditorScreen) SetStateMaterialAmbientForm() {
	msh, err := scrn.sphereModel.GetMeshByIndex(0)
	if err != nil {
		fmt.Println("Mesh is missing. Skipping update.")
		return
	}
	colorComponent := msh.(*mesh.MaterialMesh).Material.GetAmbient()
	scrn.menuModels["MaterialAmbientForm"][2].(*SliderInput).SetCurrentValue(colorComponent.X())
	scrn.menuModels["MaterialAmbientForm"][3].(*SliderInput).SetCurrentValue(colorComponent.Y())
	scrn.menuModels["MaterialAmbientForm"][4].(*SliderInput).SetCurrentValue(colorComponent.Z())
	scrn.setState("MaterialAmbientForm")
}
func (scrn *EditorScreen) SetStateMaterialDiffuseForm() {
	msh, err := scrn.sphereModel.GetMeshByIndex(0)
	if err != nil {
		fmt.Println("Mesh is missing. Skipping update.")
		return
	}
	colorComponent := msh.(*mesh.MaterialMesh).Material.GetDiffuse()
	scrn.menuModels["MaterialDiffuseForm"][2].(*SliderInput).SetCurrentValue(colorComponent.X())
	scrn.menuModels["MaterialDiffuseForm"][3].(*SliderInput).SetCurrentValue(colorComponent.Y())
	scrn.menuModels["MaterialDiffuseForm"][4].(*SliderInput).SetCurrentValue(colorComponent.Z())
	scrn.setState("MaterialDiffuseForm")
}
func (scrn *EditorScreen) SetStateMaterialSpecularForm() {
	msh, err := scrn.sphereModel.GetMeshByIndex(0)
	if err != nil {
		fmt.Println("Mesh is missing. Skipping update.")
		return
	}
	colorComponent := msh.(*mesh.MaterialMesh).Material.GetSpecular()
	scrn.menuModels["MaterialSpecularForm"][2].(*SliderInput).SetCurrentValue(colorComponent.X())
	scrn.menuModels["MaterialSpecularForm"][3].(*SliderInput).SetCurrentValue(colorComponent.Y())
	scrn.menuModels["MaterialSpecularForm"][4].(*SliderInput).SetCurrentValue(colorComponent.Z())
	scrn.setState("MaterialSpecularForm")
}
func (scrn *EditorScreen) SetStateMaterialShininessForm() {
	msh, err := scrn.sphereModel.GetMeshByIndex(0)
	if err != nil {
		fmt.Println("Mesh is missing. Skipping update.")
		return
	}
	colorComponent := msh.(*mesh.MaterialMesh).Material.GetShininess()
	scrn.menuModels["MaterialShininessForm"][2].(*SliderInput).SetCurrentValue(colorComponent)
	scrn.setState("MaterialShininessForm")
}
func (scrn *EditorScreen) setState(newState string) {
	scrn.RemoveMenuPanel()
	scrn.state = newState
	scrn.AddMenuPanel()
}

// RemoveMenuPanel removes the menu form from the screen.
func (scrn *EditorScreen) RemoveMenuPanel() {
	for index, _ := range scrn.menuModels[scrn.state] {
		switch scrn.menuModels[scrn.state][index].(type) {
		case *Button:
			item := scrn.menuModels[scrn.state][index].(*Button)
			if item.HasLabel() {
				fmt.Println("Removing the button.")
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			break
		case *TextInput:
			item := scrn.menuModels[scrn.state][index].(*TextInput)
			if item.HasLabel() {
				fmt.Println("Removing the TextInput.")
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			break
		case *SliderInput:
			item := scrn.menuModels[scrn.state][index].(*SliderInput)
			if item.HasLabel() {
				fmt.Println("Removing the SliderInput.")
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			// clean the field where the curren value is printed out.
			msh, err := item.GetMeshByIndex(2)
			if err == nil {
				scrn.charset.CleanSurface(msh)
			}
			break
		default:
			// default case for the menu panel. We have to clean it.
			msh, err := scrn.menuModels[scrn.state][index].(*model.BaseModel).GetMeshByIndex(1)
			if err == nil {
				scrn.charset.CleanSurface(msh)
			}
			break
		}
		scrn.RemoveModelFromShader(scrn.menuModels[scrn.state][index], scrn.menuShader)
	}
}

// MenuItemsDefaultState rebuilds the menu models in their default state.
func (scrn *EditorScreen) MenuItemsDefaultState() {
	for index, _ := range scrn.menuModels[scrn.state] {
		switch scrn.menuModels[scrn.state][index].(type) {
		case *Button:
			item := scrn.menuModels[scrn.state][index].(*Button)
			if item.HasLabel() {
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			item.Clear()
			break
		case *TextInput:
			item := scrn.menuModels[scrn.state][index].(*TextInput)
			if item.HasLabel() {
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			item.Clear()
			break
		case *SliderInput:
			item := scrn.menuModels[scrn.state][index].(*SliderInput)
			if item.HasLabel() {
				scrn.charset.CleanSurface(item.GetLabelSurface())
			}
			// clean the field where the curren value is printed out.
			msh, err := item.GetMeshByIndex(2)
			if err == nil {
				scrn.charset.CleanSurface(msh)
			}
			item.Clear()
			break
		}
	}
}
func (scrn *EditorScreen) releaseButtons() {
	for _, name := range AllScreenStates {
		if _, ok := scrn.menuModels[name]; !ok {
			continue
		}
		for index, _ := range scrn.menuModels[name] {
			switch scrn.menuModels[name][index].(type) {
			case *Button:
				item := scrn.menuModels[name][index].(*Button)
				item.clicked = false
			}
		}
	}
}
func (scrn *EditorScreen) Update(dt float64, p interfaces.Pointer, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	posX, posY := p.GetCurrent()
	mCoords := mgl32.Vec3{float32(-posY) / scrn.GetAspectRatio(), -FormItemsDistanceFromScreen, float32(posX)}
	scrn.UpdateWithDistance(dt, mCoords)
	closestModel, _, dist := scrn.GetClosestModelMeshDistance()
	scrn.MenuItemsDefaultState()
	switch closestModel.(type) {
	case (*Button):
		if dist < CollisionEpsilon {
			btn := closestModel.(*Button)
			if btn.HasLabel() {
				scrn.charset.CleanSurface(btn.GetLabelSurface())
			}
			btn.Hover()
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				btn.clicked = true
			}
			if btn.clicked && !buttonStore.Get(LEFT_MOUSE_BUTTON) {
				btn.clickCallback()
			}

		} else {
			scrn.releaseButtons()
		}
		break
	case (*TextInput):
		if dist < CollisionEpsilon {
			ti := closestModel.(*TextInput)
			if ti.HasLabel() {
				scrn.charset.CleanSurface(ti.GetLabelSurface())
			}
			ti.Hover()
		}
		scrn.releaseButtons()
		break
	case (*SliderInput):
		if dist < CollisionEpsilon {
			si := closestModel.(*SliderInput)
			if si.HasLabel() {
				scrn.charset.CleanSurface(si.GetLabelSurface())
			}
			// clean the field where the curren value is printed out.
			msh, err := si.GetMeshByIndex(2)
			if err == nil {
				scrn.charset.CleanSurface(msh)
			}
			si.Hover()
			// If the left mouse button is pressed, we have to detect collision against the slider.
			if buttonStore.Get(LEFT_MOUSE_BUTTON) && si.SliderCollision(mCoords) {
				dX, _ := p.GetDelta()
				si.MoveSliderWith(float32(dX))
				scrn.updateMaterialColorComponent()
			}
		}
		scrn.releaseButtons()
		break
	default:
		scrn.releaseButtons()
	}
	if MenuScreenEnabled {
		for index, _ := range scrn.menuModels[scrn.state] {
			switch scrn.menuModels[scrn.state][index].(type) {
			case *Button:
				item := scrn.menuModels[scrn.state][index].(*Button)
				if item.HasLabel() {
					pos := item.GetLabelPosition()
					scrn.charset.PrintTo(item.GetLabelText(), pos.X(), pos.Y(), pos.Z(), item.GetLabelSize()/item.aspect, scrn.GetWrapper(), item.GetLabelSurface(), []mgl32.Vec3{item.GetLabelColor()})
				}
				break
			case *TextInput:
				item := scrn.menuModels[scrn.state][index].(*TextInput)
				if item.HasLabel() {
					pos := item.GetLabelPosition()
					scrn.charset.PrintTo(item.GetLabelText(), pos.X(), pos.Y(), pos.Z(), item.GetLabelSize()/item.aspect, scrn.GetWrapper(), item.GetLabelSurface(), []mgl32.Vec3{item.GetLabelColor()})
				}
				break
			case *SliderInput:
				item := scrn.menuModels[scrn.state][index].(*SliderInput)
				if item.HasLabel() {
					pos := item.GetLabelPosition()
					scrn.charset.PrintTo(item.GetLabelText(), pos.X(), pos.Y(), pos.Z(), item.GetLabelSize()/item.aspect, scrn.GetWrapper(), item.GetLabelSurface(), []mgl32.Vec3{item.GetLabelColor()})
					// print the current value
					msh, err := item.GetMeshByIndex(2)
					if err == nil {
						valueText := transformations.Float32ToStringExact(item.sliderCurrent)
						w, _ := scrn.charset.TextContainerSize(valueText, item.GetLabelSize()/item.aspect)
						scrn.charset.PrintTo(transformations.Float32ToStringExact(item.sliderCurrent), -w/2, 0, pos.Z(), item.GetLabelSize()/item.aspect, scrn.GetWrapper(), msh, []mgl32.Vec3{item.GetLabelColor()})
					}
				}
				break
			default:
				// default case for the state label.
				// Print the state label to the top.
				msh, err := scrn.menuModels[scrn.state][index].(*model.BaseModel).GetMeshByIndex(1)
				if err == nil {
					labelSize := 0.001 / scrn.GetAspectRatio()
					w, h := scrn.charset.TextContainerSize(ScreenLabels[scrn.state], labelSize)
					scrn.charset.PrintTo(ScreenLabels[scrn.state], -w/2.0, -h, -0.01, labelSize, scrn.GetWrapper(), msh, []mgl32.Vec3{mgl32.Vec3{0, 0, 0}})
				}
				break
			}
		}
	}
}
func (scrn *EditorScreen) updateMaterialColorComponent() {
	sphereMesh, err := scrn.sphereModel.GetMeshByIndex(0)
	if err != nil {
		fmt.Println("Mesh is missing, skip material update.")
		return
	}
	origMaterial := sphereMesh.(*mesh.MaterialMesh).Material
	switch scrn.state {
	case "MaterialAmbientForm":
		red := scrn.menuModels[scrn.state][2].(*SliderInput).sliderCurrent
		green := scrn.menuModels[scrn.state][3].(*SliderInput).sliderCurrent
		blue := scrn.menuModels[scrn.state][4].(*SliderInput).sliderCurrent
		newMaterial := material.New(mgl32.Vec3{red, green, blue}, origMaterial.GetDiffuse(), origMaterial.GetSpecular(), origMaterial.GetShininess())
		sphereMesh.(*mesh.MaterialMesh).Material = newMaterial
		break
	case "MaterialDiffuseForm":
		red := scrn.menuModels[scrn.state][2].(*SliderInput).sliderCurrent
		green := scrn.menuModels[scrn.state][3].(*SliderInput).sliderCurrent
		blue := scrn.menuModels[scrn.state][4].(*SliderInput).sliderCurrent
		newMaterial := material.New(origMaterial.GetAmbient(), mgl32.Vec3{red, green, blue}, origMaterial.GetSpecular(), origMaterial.GetShininess())
		sphereMesh.(*mesh.MaterialMesh).Material = newMaterial
		break
	case "MaterialSpecularForm":
		red := scrn.menuModels[scrn.state][2].(*SliderInput).sliderCurrent
		green := scrn.menuModels[scrn.state][3].(*SliderInput).sliderCurrent
		blue := scrn.menuModels[scrn.state][4].(*SliderInput).sliderCurrent
		newMaterial := material.New(origMaterial.GetAmbient(), origMaterial.GetDiffuse(), mgl32.Vec3{red, green, blue}, origMaterial.GetShininess())
		sphereMesh.(*mesh.MaterialMesh).Material = newMaterial
		break
	case "MaterialShininessForm":
		shininess := scrn.menuModels[scrn.state][2].(*SliderInput).sliderCurrent
		newMaterial := material.New(origMaterial.GetAmbient(), origMaterial.GetDiffuse(), origMaterial.GetSpecular(), shininess)
		sphereMesh.(*mesh.MaterialMesh).Material = newMaterial
		break
	default:
		fmt.Println("Irrelevant state, skipping update.")
		break
	}
}
func (scrn *EditorScreen) defaultCharset() {
	cs, err := model.LoadCharset("assets/fonts/Desyrel/desyrel.regular.ttf", 32, 127, 20.0, 300, scrn.GetWrapper())
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	scrn.charset = cs
	// Update the position of the labels. It depends on the charset setup.
	for _, name := range AllScreenStates {
		if _, ok := scrn.menuModels[name]; !ok {
			continue
		}
		for index, _ := range scrn.menuModels[name] {
			switch scrn.menuModels[name][index].(type) {
			case *Button:
				item := scrn.menuModels[name][index].(*Button)
				if item.HasLabel() {
					w, h := scrn.charset.TextContainerSize(item.GetLabelText(), item.GetLabelSize())
					pos := item.GetLabelPosition()
					item.SetLabel(NewLabel(item.GetLabelText(), item.GetLabelColor(), mgl32.Vec3{-w / 2 / item.aspect, -h / 4, pos.Z()}, item.GetLabelSize(), item.GetLabelSurface()))
				}
				break
			case *TextInput:
				item := scrn.menuModels[name][index].(*TextInput)
				if item.HasLabel() {
					w, _ := scrn.charset.TextContainerSize(item.GetLabelText(), item.GetLabelSize())
					pos := item.GetLabelPosition()
					item.SetLabel(NewLabel(item.GetLabelText(), item.GetLabelColor(), mgl32.Vec3{-w / 2 / item.aspect, pos.Y(), pos.Z()}, item.GetLabelSize(), item.GetLabelSurface()))
				}
				break
			case *SliderInput:
				item := scrn.menuModels[name][index].(*SliderInput)
				if item.HasLabel() {
					w, _ := scrn.charset.TextContainerSize(item.GetLabelText(), item.GetLabelSize())
					pos := item.GetLabelPosition()
					item.SetLabel(NewLabel(item.GetLabelText(), item.GetLabelColor(), mgl32.Vec3{-w / 2 / item.aspect, pos.Y(), pos.Z()}, item.GetLabelSize(), item.GetLabelSurface()))
				}
				break
			}
		}
	}
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
	for index, _ := range scrn.menuModels[scrn.state] {
		switch scrn.menuModels[scrn.state][index].(type) {
		case *Button:
			item := scrn.menuModels[scrn.state][index].(*Button)
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
	app.GetWindow().SetMouseButtonCallback(app.MouseButtonCallback)
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
	fmt.Println(camera.Log())
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
func CreateMenuLabelRectangle(aspect float32) *mesh.ColorMesh {
	rect := rectangle.NewExact(1.0/aspect, 0.2/aspect)
	colors := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 1.0}}
	v, i, _ := rect.ColoredMeshInput(colors)
	menu := mesh.NewColorMesh(v, i, colors, glWrapper)
	menu.SetPosition(mgl32.Vec3{-0.9 / aspect, -0.1 / aspect, 0.0})
	menu.RotateY(-90)
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
