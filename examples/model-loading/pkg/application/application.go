package application

import (
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/mesh"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG = glfw.KeyH
)

type Shader interface {
	Use()
	SetUniformMat4(string, mgl32.Mat4)
	GetId() uint32
}

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
}

type Application struct {
	window    Window
	camera    Camera
	cameraSet bool

	shaderMap map[Shader][]*mesh.Mesh
}

type Window interface {
	GetCursorPos() (float64, float64)
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	ShouldClose() bool
	SwapBuffers()
}

// New returns an application instance
func New() *Application {
	return &Application{
		cameraSet: false,
		shaderMap: make(map[Shader][]*mesh.Mesh),
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

// AddMeshToShader attaches the mest to a shader.
func (a *Application) AddMeshToShader(m *mesh.Mesh, s Shader) {
	a.shaderMap[s] = append(a.shaderMap[s], m)
}

// Draw calls Draw function in every drawable item.
func (a *Application) Draw() {
	for s, _ := range a.shaderMap {
		s.Use()
		s.SetUniformMat4("view", a.camera.GetViewMatrix())
		s.SetUniformMat4("projection", a.camera.GetProjectionMatrix())
		for index, _ := range a.shaderMap[s] {
			a.shaderMap[s][index].Draw(s)
		}
		// Set the lighting, viewPosition, etc uniforms.
		// loop meshes, call the draw function.
	}
}
