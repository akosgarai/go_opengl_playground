package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/callbacks/primitives"
	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth          = 800
	windowHeight         = 800
	windowTitle          = "Example - shapes with camera"
	moveSpeed            = 1.0 / 1000.0
	epsilon              = 50.0
	cameraDirection      = 0.1
	cameraDirectionSpeed = 5
	// buttons
	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE
)

var (
	cameraLastUpdate int64
	app              *application.Application
)

// It creates the shader application for the items.
func CreateShader() uint32 {
	vertexShader, err := shader.CompileShader(shader.VertexShaderModelViewProjectionSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := shader.CompileShader(shader.FragmentShaderBasicSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	return program
}

// It generates a cube.
func CreateCube() *primitives.Cube {
	return primitives.NewCubeByVectorAndLength(mgl32.Vec3{0, 0, -1.5}, 2.5)
}

// It generates a Sphere.
func CreateSphere() *primitives.Sphere {
	sphere := primitives.NewSphere()
	sphere.SetCenter(mgl32.Vec3{3, 3, 5})
	sphere.SetColor(mgl32.Vec3{0, 0, 1})
	sphere.SetRadius(2.0)
	return sphere
}

func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-3, -5, 18.0}, mgl32.Vec3{0, 1, 0}, -90.0, 0.0)
	camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 100.0)
	return camera
}

func SetupKeyMap() map[glfw.Key]bool {
	keyDowns := make(map[glfw.Key]bool)
	keyDowns[FORWARD] = false
	keyDowns[LEFT] = false
	keyDowns[RIGHT] = false
	keyDowns[BACKWARD] = false
	keyDowns[UP] = false
	keyDowns[DOWN] = false

	return keyDowns
}

func Update() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))

	if epsilon > moveTime {
		return
	}
	cameraLastUpdate = nowUnix
	// Move camera
	forward := 0.0

	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = moveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -moveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	horisontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horisontal = -moveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horisontal = moveSpeed * moveTime
	}
	if horisontal != 0 {
		app.GetCamera().Strafe(float32(horisontal))
	}
	vertical := 0.0
	if app.GetKeyState(UP) && !app.GetKeyState(DOWN) {
		vertical = moveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = -moveSpeed * moveTime
	}
	if vertical != 0 {
		app.GetCamera().Lift(float32(vertical))
	}
}
func main() {
	runtime.LockOSThread()

	app = application.New()

	app.SetWindow(application.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	application.InitOpenGL()

	shaderProgramId := CreateShader()

	app.SetCamera(CreateCamera())
	cameraLastUpdate = time.Now().UnixNano()

	app.SetKeys(SetupKeyMap())

	cube := CreateCube()
	cube.SetShaderProgram(shaderProgramId)
	app.AddItem(cube)

	sphere := CreateSphere()
	sphere.SetShaderProgram(shaderProgramId)
	app.AddItem(sphere)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.GetWindow().SetMouseButtonCallback(app.DummyMouseButtonCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.DrawWithUniforms()
		app.GetWindow().SwapBuffers()
	}
}
