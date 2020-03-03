package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs and keyboard colors"
)

var (
	DebugPrint         = false
	mouseButtonPressed = false
	mousePositionX     = 0.0
	mousePositionY     = 0.0
	colorR             = 0.0
	colorG             = 0.0
	colorB             = 0.0
)

type Application struct {
	Points []primitives.Point

	camera           *primitives.Camera
	cameraLastUpdate int64
	worldUpDirection *primitives.Vector
	moveSpeed        float64
	epsilon          float64
	modelSize        int

	window  *glfw.Window
	program uint32

	KeyDowns map[string]bool
}

// NewApplication initializes the application.
func NewApplication() *Application {
	var app Application
	app.modelSize = 10
	app.moveSpeed = 1.0 / 1000.0
	app.epsilon = 50.0
	app.camera = primitives.NewCamera(primitives.Vector{0, 0, 10.0}, primitives.Vector{0, 1, 0}, -90.0, 0.0)
	app.camera.SetupProjection(45, float64(windowWidth/windowHeight), 0.1, 1000.0)

	app.cameraLastUpdate = time.Now().UnixNano()
	app.KeyDowns = make(map[string]bool)
	app.KeyDowns["W"] = false
	app.KeyDowns["A"] = false
	app.KeyDowns["S"] = false
	app.KeyDowns["D"] = false
	app.KeyDowns["Q"] = false
	app.KeyDowns["E"] = false
	return &app
}

// AddPoint inserts a new point to the points.
func (a *Application) AddPoint(point primitives.Point) {
	a.Points = append(a.Points, point)
}

// Basic function for glfw initialization.
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

	vertexShader, err := shader.CompileShader(shader.VertexShaderPointWithColorModelViewProjectionSource, gl.VERTEX_SHADER)
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

/*
* Mouse click handler logic:
* - if the left mouse button is not pressed, and the button is just released, App.AddPoint(), clean up the temp.point.
* - if the button is just pressed, set the point that needs to be added.
 */
func (a *Application) MouseHandler() {
	x, y := a.window.GetCursorPos()

	if a.window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = x
			mousePositionY = y
			mouseButtonPressed = true
		}
	} else {
		if mouseButtonPressed {
			mouseButtonPressed = false
			x, y := convertMouseCoordinates()
			positionVector := primitives.Vector{x, y, -1.0}
			scalarPart := float64(a.modelSize / 2)
			modelPositionVector := (positionVector.MultiplyScalar(scalarPart)).AddScalar(scalarPart)
			a.AddPoint(
				primitives.Point{
					modelPositionVector,
					primitives.Vector{colorR, colorG, colorB},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func (a *Application) KeyHandler() {
	if a.window.GetKey(glfw.KeyH) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Print("app.Points:\n")
			fmt.Println(a.Points)
			fmt.Printf("app.camera: %s\n", a.camera.Log())
		}
	} else {
		DebugPrint = false
	}
	if a.window.GetKey(glfw.KeyR) == glfw.Press {
		colorR = 1
	} else {
		colorR = 0
	}
	if a.window.GetKey(glfw.KeyG) == glfw.Press {
		colorG = 1
	} else {
		colorG = 0
	}
	if a.window.GetKey(glfw.KeyB) == glfw.Press {
		colorB = 1
	} else {
		colorB = 0
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
		a.camera.Walk(forward)
	}
	horisontal := 0.0
	if a.KeyDowns["A"] && !a.KeyDowns["D"] {
		horisontal = -a.moveSpeed * moveTime
	} else if a.KeyDowns["D"] && !a.KeyDowns["A"] {
		horisontal = a.moveSpeed * moveTime
	}
	if horisontal != 0 {
		a.camera.Strafe(horisontal)
	}
	vertical := 0.0
	if a.KeyDowns["Q"] && !a.KeyDowns["E"] {
		vertical = -a.moveSpeed * moveTime
	} else if a.KeyDowns["E"] && !a.KeyDowns["Q"] {
		vertical = a.moveSpeed * moveTime
	}
	if vertical != 0 {
		a.camera.Lift(vertical)
	}
}
func convertMouseCoordinates() (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (mousePositionX - halfWidth) / (halfWidth)
	y := (halfHeight - mousePositionY) / (halfHeight)
	return x, y
}
func (a *Application) buildVAO() []float32 {
	var vao []float32
	for _, item := range a.Points {
		vao = append(vao, float32(item.Coordinate.X))
		vao = append(vao, float32(item.Coordinate.Y))
		vao = append(vao, float32(item.Coordinate.Z))
		vao = append(vao, float32(item.Color.X))
		vao = append(vao, float32(item.Color.Y))
		vao = append(vao, float32(item.Color.Z))
	}
	return vao
}
func (a *Application) Draw() {
	if len(a.Points) < 1 {
		return
	}
	points := a.buildVAO()
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))

	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.POINTS, 0, int32(len(a.Points)))
}

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()

	gl.UseProgram(app.program)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	modelLocation := gl.GetUniformLocation(app.program, gl.Str("model\x00"))
	viewLocation := gl.GetUniformLocation(app.program, gl.Str("view\x00"))
	projectionLocation := gl.GetUniformLocation(app.program, gl.Str("projection\x00"))

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		app.KeyHandler()
		app.MouseHandler()

		M := primitives.UnitMatrix4x4().GetMatrix()
		gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])
		V := app.camera.GetViewMatrix().GetMatrix()
		gl.UniformMatrix4fv(viewLocation, 1, false, &V[0])
		P := app.camera.GetProjectionMatrix().GetMatrix()
		gl.UniformMatrix4fv(projectionLocation, 1, false, &P[0])

		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
