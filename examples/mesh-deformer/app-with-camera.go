package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	mat "github.com/akosgarai/opengl_playground/pkg/primitives/matrix"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	tr "github.com/akosgarai/opengl_playground/pkg/primitives/triangle"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - mesh deformer"
	moveSpeed    = 5.0
)

var (
	cameraRotate = false
	DebugPrint   = false
)

type Application struct {
	triangles             []tr.Triangle
	defaultTriangleLength int
	triangleColorFront    vec.Vector
	triangleColorBack     vec.Vector

	cameraLastUpdate     int64
	camera               *camera.Camera
	cameraDirection      float64
	cameraDirectionSpeed float64
	moveSpeed            float64
	epsilon              float64

	worldWidth       int
	worldHeight      int
	worldDepth       int
	worldUpDirection *vec.Vector

	window  *glfw.Window
	program uint32

	KeyDowns map[string]bool
}

func NewApplication() *Application {
	var app Application
	app.worldWidth = 1000
	app.worldHeight = 1000
	app.worldDepth = 1000
	app.defaultTriangleLength = 10
	app.triangleColorFront = vec.Vector{0, 0, 1}
	app.triangleColorBack = vec.Vector{0, 0.5, 1}
	app.GenerateTriangles()

	app.moveSpeed = 10.0 / 1000.0
	app.epsilon = 50.0

	app.camera = camera.NewCamera(mgl32.Vec3{100.0, 100.0, 100.0}, mgl32.Vec3{0, 1, 0}, -180.0, 0.0)
	app.cameraDirection = 0.1
	app.cameraDirectionSpeed = 50
	fmt.Println("Camera state after new function")
	fmt.Println(app.camera.Log())
	app.camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 150.0)
	fmt.Println("Camera state after setupProjection function")
	fmt.Println(app.camera.Log())
	app.KeyDowns = make(map[string]bool)
	app.KeyDowns["W"] = false
	app.KeyDowns["A"] = false
	app.KeyDowns["S"] = false
	app.KeyDowns["D"] = false
	app.KeyDowns["Q"] = false
	app.KeyDowns["E"] = false
	app.KeyDowns["dLeft"] = false
	app.KeyDowns["dRight"] = false
	app.KeyDowns["dUp"] = false
	app.KeyDowns["dDown"] = false
	return &app
}

func (a *Application) AddTriangle(triangle tr.Triangle) {
	a.triangles = append(a.triangles, triangle)
}

// It generates a bunch of triangles and sets their color to static blue.
func (a *Application) GenerateTriangles() {
	rows := a.worldDepth / a.defaultTriangleLength
	cols := a.worldWidth / a.defaultTriangleLength
	length := a.defaultTriangleLength
	for i := 0; i <= rows; i++ {
		for j := 0; j <= cols; j++ {
			topX := float64(j * length)
			topY := float64(i * length)
			topZ := 0.0

			triangle := *tr.NewTriangle(
				vec.Vector{topX, topY, topZ},
				vec.Vector{topX, topY - float64(length), topZ},
				vec.Vector{topX - float64(length), topY - float64(length), topZ},
			)
			triangle.SetColor(a.triangleColorFront)
			a.AddTriangle(triangle)
			triangle = *tr.NewTriangle(
				vec.Vector{topX, topY, topZ},
				vec.Vector{topX - float64(length), topY - float64(length), topZ},
				vec.Vector{topX - float64(length), topY, topZ},
			)
			triangle.SetColor(a.triangleColorBack)
			a.AddTriangle(triangle)
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func (a *Application) KeyHandler() {
	if a.window.GetKey(glfw.KeyH) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Printf("app.camera: %s\n", a.camera.Log())
		}
	} else {
		DebugPrint = false
	}
	if a.window.GetKey(glfw.KeyW) == glfw.Press {
		a.KeyDowns["W"] = true
	} else {
		a.KeyDowns["W"] = false
	}
	if a.window.GetKey(glfw.KeyA) == glfw.Press {
		a.KeyDowns["A"] = true
	} else {
		a.KeyDowns["A"] = false
	}
	if a.window.GetKey(glfw.KeyS) == glfw.Press {
		a.KeyDowns["S"] = true
	} else {
		a.KeyDowns["S"] = false
	}
	if a.window.GetKey(glfw.KeyD) == glfw.Press {
		a.KeyDowns["D"] = true
	} else {
		a.KeyDowns["D"] = false
	}
	if a.window.GetKey(glfw.KeyQ) == glfw.Press {
		a.KeyDowns["Q"] = true
	} else {
		a.KeyDowns["Q"] = false
	}
	if a.window.GetKey(glfw.KeyE) == glfw.Press {
		a.KeyDowns["E"] = true
	} else {
		a.KeyDowns["E"] = false
	}
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - a.cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))
	// if the camera has been updated recently, we can skip now
	if a.epsilon > moveTime {
		return
	}
	a.cameraLastUpdate = nowUnix
	// Move camera
	forward := 0.0
	if a.KeyDowns["W"] && !a.KeyDowns["S"] {
		forward = a.moveSpeed * moveTime
	} else if a.KeyDowns["S"] && !a.KeyDowns["W"] {
		forward = -a.moveSpeed * moveTime
	}
	if forward != 0 {
		a.camera.Walk(float32(forward))
	}
	horisontal := 0.0
	if a.KeyDowns["A"] && !a.KeyDowns["D"] {
		horisontal = -a.moveSpeed * moveTime
	} else if a.KeyDowns["D"] && !a.KeyDowns["A"] {
		horisontal = a.moveSpeed * moveTime
	}
	if horisontal != 0 {
		a.camera.Strafe(float32(horisontal))
	}
	vertical := 0.0
	if a.KeyDowns["Q"] && !a.KeyDowns["E"] {
		vertical = -a.moveSpeed * moveTime
	} else if a.KeyDowns["E"] && !a.KeyDowns["Q"] {
		vertical = a.moveSpeed * moveTime
	}
	if vertical != 0 {
		a.camera.Lift(float32(vertical))
	}
}

/*
* Mouse click handler logic:
* - if the mouse moved - call the function with the delta values.
 */
func (a *Application) MouseHandler() {
	currX, currY := a.window.GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, windowWidth, windowHeight)
	// dUp
	if y > 1.0-a.cameraDirection && y < 1.0 {
		a.KeyDowns["dUp"] = true
	} else {
		a.KeyDowns["dUp"] = false
	}
	// dDown
	if y < -1.0+a.cameraDirection && y > -1.0 {
		a.KeyDowns["dDown"] = true
	} else {
		a.KeyDowns["dDown"] = false
	}
	// dLeft
	if x < -1.0+a.cameraDirection && x > -1.0 {
		a.KeyDowns["dLeft"] = true
	} else {
		a.KeyDowns["dLeft"] = false
	}
	// dRight
	if x > 1.0-a.cameraDirection && x < 1.0 {
		a.KeyDowns["dRight"] = true
	} else {
		a.KeyDowns["dRight"] = false
	}

	dX := 0.0
	dY := 0.0
	if a.KeyDowns["dUp"] && !a.KeyDowns["dDown"] {
		dY = 0.01 * a.cameraDirectionSpeed
	} else if a.KeyDowns["dDown"] && !a.KeyDowns["dUp"] {
		dY = -0.01 * a.cameraDirectionSpeed
	}
	if a.KeyDowns["dLeft"] && !a.KeyDowns["dRight"] {
		dX = -0.01 * a.cameraDirectionSpeed
	} else if a.KeyDowns["dRight"] && !a.KeyDowns["dLeft"] {
		dX = 0.01 * a.cameraDirectionSpeed
	}
	a.camera.UpdateDirection(float32(dX), float32(dY))
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertexShader, err := shader.CompileShader(shader.VertexShaderDeformVertexPositionModelViewProjectionSource, gl.VERTEX_SHADER)
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

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	// Configure global settings
	gl.UseProgram(app.program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	nowUnix := time.Now().UnixNano()

	modelLocation := gl.GetUniformLocation(app.program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(app.program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(app.program, gl.Str("projection\x00"))
	timeLocation := gl.GetUniformLocation(app.program, gl.Str("time\x00"))

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// time
		elapsedTimeNano := time.Now().UnixNano() - nowUnix
		gl.Uniform1f(timeLocation, float32(elapsedTimeNano/10000000))

		M := mat.UnitMatrix().GetMatrix()
		gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
		V := app.camera.GetViewMatrix()
		gl.UniformMatrix4fv(viewLocation, 1, false, &V[0])
		P := app.camera.GetProjectionMatrix()
		gl.UniformMatrix4fv(projectionLocation, 1, false, &P[0])

		for _, item := range app.triangles {
			item.Draw()
		}
		app.KeyHandler()
		app.MouseHandler()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
