package application

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/interfaces"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG  = glfw.KeyH
	EXPORT = glfw.KeyP
)

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
	GetVelocity() float32
}

type Window interface {
	GetCursorPos() (float64, float64)
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	ShouldClose() bool
	SwapBuffers()
}

type Application struct {
	window    Window
	camera    Camera
	cameraSet bool

	shaderMap  map[interfaces.Shader][]interfaces.Model
	mouseDowns map[glfw.MouseButton]bool
	MousePosX  float64
	MousePosY  float64

	directionalLightSources []DirectionalLightSource
	pointLightSources       []PointLightSource
	spotLightSources        []SpotLightSource

	keyDowns map[glfw.Key]bool

	// it holds the keyMaps for the camera movement
	cameraKeyboardMovementMap map[string]glfw.Key
}

// New returns an application instance
func New() *Application {
	return &Application{
		cameraSet:                 false,
		shaderMap:                 make(map[interfaces.Shader][]interfaces.Model),
		mouseDowns:                make(map[glfw.MouseButton]bool),
		directionalLightSources:   []DirectionalLightSource{},
		pointLightSources:         []PointLightSource{},
		spotLightSources:          []SpotLightSource{},
		keyDowns:                  make(map[glfw.Key]bool),
		cameraKeyboardMovementMap: make(map[string]glfw.Key),
	}
}

// Log returns the string representation of this object.
func (a *Application) Log() string {
	logString := "Application:\n"
	if a.cameraSet {
		logString += " - camera : " + a.camera.Log() + "\n"
	}
	return logString
}

// SetCameraMovementMap sets the cameraKeyboardMovementMap variable.
// Currently the following values are supported: 'forward', 'back',
// 'left', 'right', 'up', 'down'
func (a *Application) SetCameraMovementMap(m map[string]glfw.Key) {
	a.cameraKeyboardMovementMap = m
}

// SetWindow updates the window with the new one.
func (a *Application) SetWindow(w Window) {
	a.window = w
}

// GetWindow returns the current window of the application.
func (a *Application) GetWindow() Window {
	return a.window
}

// SetCamera updates the camera with the new one.
func (a *Application) SetCamera(c Camera) {
	a.cameraSet = true
	a.camera = c
}

// GetCamera returns the current camera of the application.
func (a *Application) GetCamera() Camera {
	return a.camera
}

// AddShader method inserts the new shader to the shaderMap
func (a *Application) AddShader(s interfaces.Shader) {
	a.shaderMap[s] = []interfaces.Model{}
}

// AddModelToShader attaches the model to a shader.
func (a *Application) AddModelToShader(m interfaces.Model, s interfaces.Shader) {
	a.shaderMap[s] = append(a.shaderMap[s], m)
}

// cameraKeyboardMovement is responsible for handling a movement for a specific direction.
// The direction is described by the key strings. The handler function name is also added
// as input to be able to call it. For the movement we also need to know the delta time,
// that is also added as function input. In case of invalid function name,
// it prints out some message to the console.
func (a *Application) cameraKeyboardMovement(directionKey, oppositeKey, handlerName string, delta float64) {
	keyStateDirection := false
	keyStateOpposite := false
	if val, ok := a.cameraKeyboardMovementMap[directionKey]; ok {
		keyStateDirection = a.GetKeyState(val)
	}
	if val, ok := a.cameraKeyboardMovementMap[oppositeKey]; ok {
		keyStateOpposite = a.GetKeyState(val)
	}
	step := float32(0.0)
	if keyStateDirection && !keyStateOpposite {
		step = float32(delta) * a.camera.GetVelocity()
	} else if keyStateOpposite && !keyStateDirection {
		step = -float32(delta) * a.camera.GetVelocity()
	}
	if step != 0 {
		method := reflect.ValueOf(a.camera).MethodByName(handlerName)
		if method.IsZero() {
			fmt.Printf("Invalid method name '%s' was given for camera movement.\n", handlerName)
			return
		}
		var inputParams []reflect.Value
		inputParams = append(inputParams, reflect.ValueOf(step))
		method.Call(inputParams)
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement (the rotation is unhandled yet), if
// the camera is set.
func (a *Application) Update(dt float64) {
	if a.cameraSet {
		a.cameraKeyboardMovement("forward", "back", "Walk", dt)
		a.cameraKeyboardMovement("right", "left", "Strafe", dt)
		a.cameraKeyboardMovement("up", "down", "Lift", dt)
	}
	for s, _ := range a.shaderMap {
		for index, _ := range a.shaderMap[s] {
			a.shaderMap[s][index].Update(dt)
		}
	}
}

// Draw calls Draw function in every drawable item. It loops on the shaderMap (shaders).
// For each shader, first set it to used state, setup camera realted uniforms,
// then setup light related uniforms. Then we can pass the shader to the Model for drawing.
func (a *Application) Draw() {
	for s, _ := range a.shaderMap {
		s.Use()
		if a.cameraSet {
			s.SetUniformMat4("view", a.camera.GetViewMatrix())
			s.SetUniformMat4("projection", a.camera.GetProjectionMatrix())
			cameraPos := a.camera.GetPosition()
			s.SetUniform3f("viewPosition", cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		} else {
			s.SetUniformMat4("view", mgl32.Ident4())
			s.SetUniformMat4("projection", mgl32.Ident4())
		}
		a.lightHandler(s)
		for index, _ := range a.shaderMap[s] {
			a.shaderMap[s][index].Draw(s)
		}
	}
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case DEBUG:
		if action != glfw.Release {
			fmt.Printf("%s\n", a.Log())
		}
		break
	case EXPORT:
		if action != glfw.Release {
			a.export()
		}
		break
	default:
		a.SetKeyState(key, action)
		break
	}
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	a.MousePosX, a.MousePosY = w.GetCursorPos()
	switch button {
	default:
		a.SetButtonState(button, action)
		break
	}
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetKeyState(key glfw.Key, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.keyDowns[key] = isButtonPressed
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetButtonState(button glfw.MouseButton, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.mouseDowns[button] = isButtonPressed
}

// GetMouseButtonState returns the state of the given button
func (a *Application) GetMouseButtonState(button glfw.MouseButton) bool {
	return a.mouseDowns[button]
}

// GetKeyState returns the state of the given key
func (a *Application) GetKeyState(key glfw.Key) bool {
	return a.keyDowns[key]
}

// Setup light related uniforms.
func (a *Application) lightHandler(s interfaces.Shader) {
	a.setupDirectionalLightForShader(s)
	a.setupPointLightForShader(s)
	a.setupSpotLightForShader(s)
}

// Setup directional light related uniforms. It iterates over the directional sources
// and setups each uniform, where the name is not empty.
func (a *Application) setupDirectionalLightForShader(s interfaces.Shader) {
	for _, source := range a.directionalLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
	}
	s.SetUniform1i("NumberOfDirectionalLightSources", int32(len(a.directionalLightSources)))

}

// Setup point light relates uniforms. It iterates over the point light sources and sets
// up every uniform, where the name is not empty.
func (a *Application) setupPointLightForShader(s interfaces.Shader) {
	for _, source := range a.pointLightSources {
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
	}
	s.SetUniform1i("NumberOfPointLightSources", int32(len(a.pointLightSources)))
}

// Setup spot light related uniforms. It iterates over the spot light sources and sets up
// every uniform, where the name is not empty.
func (a *Application) setupSpotLightForShader(s interfaces.Shader) {
	for _, source := range a.spotLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
		if source.CutoffUniformName != "" {
			s.SetUniform1f(source.CutoffUniformName, source.LightSource.GetCutoff())
		}
		if source.OuterCutoffUniformName != "" {
			s.SetUniform1f(source.OuterCutoffUniformName, source.LightSource.GetOuterCutoff())
		}
	}
	s.SetUniform1i("NumberOfSpotLightSources", int32(len(a.spotLightSources)))
}

// AddDirectionalLightSource sets up a directional light source.
// It takes a DirectionalLight input that contains the model related info,
// and it also takes a [4]string, with the uniform names that are used in the shader applications
// the 'DirectionUniformName', 'AmbientUniformName', 'DiffuseUniformName', 'SpecularUniformName'.
// They has to be in this order.
func (a *Application) AddDirectionalLightSource(lightSource interfaces.DirectionalLight, uniformNames [4]string) {
	var dSource DirectionalLightSource
	dSource.LightSource = lightSource
	dSource.DirectionUniformName = uniformNames[0]
	dSource.AmbientUniformName = uniformNames[1]
	dSource.DiffuseUniformName = uniformNames[2]
	dSource.SpecularUniformName = uniformNames[3]

	a.directionalLightSources = append(a.directionalLightSources, dSource)
}

// AddPointLightSource sets up a point light source. It takes a PointLight
// input that contains the model related info, and it also containt the uniform names in [7]string format.
// The order has to be the following: 'PositionUniformName', 'AmbientUniformName', 'DiffuseUniformName',
// 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName', 'QuadraticTermUniformName'.
func (a *Application) AddPointLightSource(lightSource interfaces.PointLight, uniformNames [7]string) {
	var pSource PointLightSource
	pSource.LightSource = lightSource
	pSource.PositionUniformName = uniformNames[0]
	pSource.AmbientUniformName = uniformNames[1]
	pSource.DiffuseUniformName = uniformNames[2]
	pSource.SpecularUniformName = uniformNames[3]
	pSource.ConstantTermUniformName = uniformNames[4]
	pSource.LinearTermUniformName = uniformNames[5]
	pSource.QuadraticTermUniformName = uniformNames[6]

	a.pointLightSources = append(a.pointLightSources, pSource)
}

// AddSpotLightSource sets up a spot light source. It takes a SpotLight input
// that contains the model related info, and it also contains the uniform names in [10]string format.
// The order has to be the following: 'PositionUniformName', 'DirectionUniformName', 'AmbientUniformName',
// 'DiffuseUniformName', 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName',
// 'QuadraticTermUniformName', 'CutoffUniformName'.
func (a *Application) AddSpotLightSource(lightSource interfaces.SpotLight, uniformNames [10]string) {
	var sSource SpotLightSource
	sSource.LightSource = lightSource
	sSource.PositionUniformName = uniformNames[0]
	sSource.DirectionUniformName = uniformNames[1]
	sSource.AmbientUniformName = uniformNames[2]
	sSource.DiffuseUniformName = uniformNames[3]
	sSource.SpecularUniformName = uniformNames[4]
	sSource.ConstantTermUniformName = uniformNames[5]
	sSource.LinearTermUniformName = uniformNames[6]
	sSource.QuadraticTermUniformName = uniformNames[7]
	sSource.CutoffUniformName = uniformNames[8]
	sSource.OuterCutoffUniformName = uniformNames[8]

	a.spotLightSources = append(a.spotLightSources, sSource)
}

// This function is called for starting the export process. It is attached to a key callback.
func (a *Application) export() {
	ExportBaseDir := "./exports"
	Directory := time.Now().Format("20060102150405")
	err := os.Mkdir(ExportBaseDir+"/"+Directory, os.ModeDir|os.ModePerm)
	if err != nil {
		fmt.Printf("Cannot create export directory. '%s'\n", err.Error())
	}
	i := 0
	for s, _ := range a.shaderMap {
		modelDir := strconv.Itoa(i)
		err := os.Mkdir(ExportBaseDir+"/"+Directory+"/"+modelDir, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Cannot create model directory. '%s'\n", err.Error())
		}
		for index, _ := range a.shaderMap[s] {
			a.shaderMap[s][index].Export(ExportBaseDir + "/" + Directory + "/" + modelDir)
		}
		i++
	}
}
