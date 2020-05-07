package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	WindowTitle  = "Example - mesh deformer"

	rows   = 10
	cols   = 10
	length = 10
)

var (
	app *application.Application

	triangleColorFront = mgl32.Vec3{0, 0, 1}
	triangleColorBack  = mgl32.Vec3{0, 0.5, 1}

	lastUpdate int64
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{4.0, 9.7, -0.3}, mgl32.Vec3{0, 0, 1}, 13.0, 58.5)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

// It generates a bunch of triangles and sets their color to static blue.
func GenerateTriangles(shaderProgram *shader.Shader) {
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float32(j * length)
			topY := float32(i * length)
			topZ := float32(0.0)

			coords := [3]mgl32.Vec3{
				mgl32.Vec3{topX, topY, topZ},
				mgl32.Vec3{topX, topY - float32(length), topZ},
				mgl32.Vec3{topX - float32(length), topY - float32(length), topZ},
			}
			colors := [3]mgl32.Vec3{triangleColorFront, triangleColorFront, triangleColorFront}
			item := triangle.New(coords, colors, shaderProgram)
			item.SetDirection(mgl32.Vec3{0, 0, 1})
			item.SetSpeed(float32(1.0) / float32(1000000000.0))
			app.AddItem(item)

			coords = [3]mgl32.Vec3{
				mgl32.Vec3{topX, topY, topZ},
				mgl32.Vec3{topX - float32(length), topY - float32(length), topZ},
				mgl32.Vec3{topX - float32(length), topY, topZ},
			}
			colors = [3]mgl32.Vec3{triangleColorBack, triangleColorBack, triangleColorBack}
			item = triangle.New(coords, colors, shaderProgram)
			item.SetDirection(mgl32.Vec3{0, 0, 1})
			item.SetSpeed(float32(1.0) / float32(1000000000.0))
			app.AddItem(item)
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
	shader.InitOpenGL()

	app.SetCamera(CreateCamera())

	shaderProgram := shader.NewShader("examples/04-mesh-deformer/vertexshader.vert", "examples/04-mesh-deformer/fragmentshader.frag")

	GenerateTriangles(shaderProgram)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
