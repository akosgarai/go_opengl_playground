package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

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

	glWrapper wrapper.Wrapper
)

func GenerateTriangleMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	triang := triangle.New(60, 60, 60)
	v, i := triang.ColoredMeshInput(col)
	return mesh.NewColorMesh(v, i, col, glWrapper)
}
func GenerateSquareMesh(col []mgl32.Vec3) *mesh.ColorMesh {
	square := rectangle.NewSquare()
	v, i := square.ColoredMeshInput(col)
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

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	shaderProgram := shader.NewShader("examples/03-button-handler/shaders/vertexshader.vert", "examples/03-button-handler/shaders/fragmentshader.frag", glWrapper)
	app.AddShader(shaderProgram)

	TriangMesh = GenerateTriangleMesh(triangleColors)
	TriangMesh.SetRotationAngle(mgl32.DegToRad(90))
	TriangMesh.SetRotationAxis(mgl32.Vec3{1, 0, 0})
	TriangMesh.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	TriangMesh.SetPosition(mgl32.Vec3{-0.4, 0, 0.3})
	TriangMesh.SetSpeed(speed)
	mod.AddMesh(TriangMesh)

	SquareMesh = GenerateSquareMesh(squareColors)
	SquareMesh.SetRotationAngle(mgl32.DegToRad(90))
	SquareMesh.SetRotationAxis(mgl32.Vec3{1, 0, 0})
	SquareMesh.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	SquareMesh.SetPosition(mgl32.Vec3{0.4, 0, -0.3})
	SquareMesh.SetSpeed(speed)
	mod.AddMesh(SquareMesh)
	app.AddModelToShader(mod, shaderProgram)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
