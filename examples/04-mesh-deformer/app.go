package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - mesh deformer"

	rows   = 10
	cols   = 10
	length = 10
)

var (
	app *application.Application

	triangleColorFront = []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	triangleColorBack  = []mgl32.Vec3{mgl32.Vec3{0, 0.5, 1}}

	lastUpdate int64
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0, 1, 0}, 0.0, 90.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

// It generates a bunch of triangles and sets their color to static blue.
func GenerateTriangles(shaderProgram *shader.Shader) {
	triang := triangle.New(90, 45, 45)
	v1, indicies1 := triang.ColoredMeshInput(triangleColorFront)
	v2, indicies2 := triang.ColoredMeshInput(triangleColorBack)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j * length)
			topY := float32(i * length)
			topZ := float32(0.0)

			m := mesh.NewColorMesh(v1, indicies1)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetScale(mgl32.Vec3{length, length, length})
			m.SetDirection(mgl32.Vec3{0, 0, 1})
			m.SetSpeed(float32(1.0) / float32(1000000000.0))

			app.AddMeshToShader(m, shaderProgram)

			m2 := mesh.NewColorMesh(v2, indicies2)
			m2.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetRotationAngle(mgl32.DegToRad(45))
			m2.SetRotationAxis(mgl32.Vec3{1, -1, 0})
			m2.SetDirection(mgl32.Vec3{0, 0, 1})
			m2.SetSpeed(float32(1.0) / float32(1000000000.0))

			app.AddMeshToShader(m2, shaderProgram)
		}
	}
}

// Update the z coordinates of the vectors.
func Update() {
	now := time.Now().UnixNano()
	delta := float64(now - lastUpdate)
	app.Update(delta)
	lastUpdate = now

}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	wrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/04-mesh-deformer/shaders/vertexshader.vert", "examples/04-mesh-deformer/shaders/fragmentshader.frag")
	app.AddShader(shaderProgram)

	GenerateTriangles(shaderProgram)

	lastUpdate = time.Now().UnixNano()

	wrapper.Enable(wrapper.DEPTH_TEST)
	wrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		wrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
