package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/triangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

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
	Model              = model.New()

	lastUpdate int64

	glWrapper glwrapper.Wrapper
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0, 1, 0}, 0.0, 90.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

// It generates a bunch of triangles and sets their color to static blue.
func GenerateTriangles() {
	triang := triangle.New(90, 45, 45)
	v1, indices1, _ := triang.ColoredMeshInput(triangleColorFront)
	v2, indices2, _ := triang.ColoredMeshInput(triangleColorBack)
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j * length)
			topY := float32(i * length)
			topZ := float32(0.0)

			m := mesh.NewColorMesh(v1, indices1, triangleColorFront, glWrapper)
			m.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m.SetScale(mgl32.Vec3{length, length, length})
			m.SetDirection(mgl32.Vec3{0, 0, 1})
			m.SetSpeed(float32(1.0) / float32(1000000000.0))

			Model.AddMesh(m)

			m2 := mesh.NewColorMesh(v2, indices2, triangleColorBack, glWrapper)
			m2.SetPosition(mgl32.Vec3{topX, topY, topZ})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetScale(mgl32.Vec3{length, length, length})
			m2.SetDirection(mgl32.Vec3{0, 0, 1})
			m2.SetSpeed(float32(1.0) / float32(1000000000.0))

			Model.AddMesh(m2)
		}
	}
	Model.RotateX(45)
}

// Update the z coordinates of the vectors.
func Update() {
	now := time.Now().UnixNano()
	delta := float64(now - lastUpdate)
	app.Update(delta)
	lastUpdate = now

}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	scrn.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)

	GenerateTriangles()
	scrn.AddModelToShader(Model, shaderProgram)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	lastUpdate = time.Now().UnixNano()

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
