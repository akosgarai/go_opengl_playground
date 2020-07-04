package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
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
	WindowTitle  = "Example - static button handler"
	epsilon      = 30
	speed        = float32(0.15) / float32(1000.0)

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	SWAP     = glfw.KeyT
	DEBUG    = glfw.KeyH
)

var (
	triangleColors = []mgl32.Vec3{
		mgl32.Vec3{0, 1, 0},
	}
	squareColors = []mgl32.Vec3{
		mgl32.Vec3{0, 0, 1},
	}
	lastUpdate = time.Now().UnixNano() / 1000000
	mod        = model.New()

	app *application.Application

	TriangMesh *mesh.ColorMesh
	SquareMesh *mesh.ColorMesh

	glWrapper glwrapper.Wrapper
)

func GenerateTriangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	triang := triangle.New(60, 60, 60)
	v, i, _ := triang.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}
func GenerateSquareMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := rectangle.NewSquare()
	v, i, _ := square.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}

func Update() {
	// now in milisec.
	nowUnixM := time.Now().UnixNano() / 1000000
	delta := nowUnixM - lastUpdate
	sqDir := SquareMesh.GetDirection()
	trDir := TriangMesh.GetDirection()
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		SquareMesh.SetDirection(mgl32.Vec3{sqDir.X(), 1, sqDir.Y()})
		TriangMesh.SetDirection(mgl32.Vec3{trDir.X(), -1, trDir.Y()})
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		SquareMesh.SetDirection(mgl32.Vec3{sqDir.X(), -1, sqDir.Y()})
		TriangMesh.SetDirection(mgl32.Vec3{trDir.X(), 1, trDir.Y()})
	} else {
		SquareMesh.SetDirection(mgl32.Vec3{sqDir.X(), 0, sqDir.Y()})
		TriangMesh.SetDirection(mgl32.Vec3{trDir.X(), 0, trDir.Y()})
	}
	sqDir = SquareMesh.GetDirection()
	trDir = TriangMesh.GetDirection()
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		SquareMesh.SetDirection(mgl32.Vec3{-1, sqDir.Y(), sqDir.Z()})
		TriangMesh.SetDirection(mgl32.Vec3{1, trDir.Y(), trDir.Z()})
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		SquareMesh.SetDirection(mgl32.Vec3{1, sqDir.Y(), sqDir.Z()})
		TriangMesh.SetDirection(mgl32.Vec3{-1, trDir.Y(), trDir.Z()})
	} else {
		SquareMesh.SetDirection(mgl32.Vec3{0, sqDir.Y(), sqDir.Z()})
		TriangMesh.SetDirection(mgl32.Vec3{0, trDir.Y(), trDir.Z()})
	}
	if epsilon > delta {
		return
	}
	lastUpdate = nowUnixM
	app.Update(float64(delta))
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	shaderProgram := shader.NewShader(baseDir()+"/shaders/vertexshader.vert", baseDir()+"/shaders/fragmentshader.frag", glWrapper)
	scrn.AddShader(shaderProgram)

	TriangMesh = GenerateTriangleMesh(triangleColors)
	TriangMesh.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	TriangMesh.SetPosition(mgl32.Vec3{-0.4, 0, 0.3})
	TriangMesh.SetSpeed(speed)
	mod.AddMesh(TriangMesh)

	SquareMesh = GenerateSquareMesh(squareColors)
	SquareMesh.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	SquareMesh.SetPosition(mgl32.Vec3{0.4, 0, -0.3})
	SquareMesh.SetSpeed(speed)
	mod.AddMesh(SquareMesh)
	mod.RotateX(90)
	scrn.AddModelToShader(mod, shaderProgram)
	scrn.Setup(setupApp)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
