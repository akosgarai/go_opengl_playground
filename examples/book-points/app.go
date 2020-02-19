package main

import (
	"fmt"
	"runtime"

	C "github.com/akosgarai/opengl_playground/pkg/camera"
	M "github.com/akosgarai/opengl_playground/pkg/matrix"
	"github.com/akosgarai/opengl_playground/pkg/primitives"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	V "github.com/akosgarai/opengl_playground/pkg/vector"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - draw points from mouse inputs - graph book based solution"
)

var (
	DebugPrint         = false
	mouseButtonPressed = false
	mousePositionX     = 0.0
	mousePositionY     = 0.0
	colorR             = 0.0
	colorG             = 0.0
	colorB             = 0.0
	modelSize          = 100
)

type Application struct {
	Points []primitives.PointBB

	KeyDowns map[string]bool
}

// AddPoint inserts a new point to the points.
func (a *Application) AddPoint(point primitives.PointBB) {
	a.Points = append(a.Points, point)
}

var app Application

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

	vertexShader, err := shader.CompileShader(shader.VertexShaderPointWithColorMVPSource, gl.VERTEX_SHADER)
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
func mouseHandler(window *glfw.Window) {
	x, y := window.GetCursorPos()

	if window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press {
		if !mouseButtonPressed {
			mousePositionX = x
			mousePositionY = y
			mouseButtonPressed = true
		}
	} else {
		if mouseButtonPressed {
			mouseButtonPressed = false
			x, y := convertMouseCoordinates()
			positionVector := V.Vector{x, y, 0.0}
			scalarPart := float64(modelSize / 2)
			modelPositionVector := (positionVector.MultiplyScalar(scalarPart)).Add(V.Vector{scalarPart, scalarPart, scalarPart})
			app.AddPoint(
				primitives.PointBB{
					modelPositionVector,
					V.Vector{colorR, colorG, colorB},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyH) == glfw.Press {
		if !DebugPrint {
			DebugPrint = true
			fmt.Println(app.Points)
		}
	} else {
		DebugPrint = false
	}
	if window.GetKey(glfw.KeyR) == glfw.Press {
		colorR = 1
	} else {
		colorR = 0
	}
	if window.GetKey(glfw.KeyG) == glfw.Press {
		colorG = 1
	} else {
		colorG = 0
	}
	if window.GetKey(glfw.KeyB) == glfw.Press {
		colorB = 1
	} else {
		colorB = 0
	}
}
func convertMouseCoordinates() (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (mousePositionX - halfWidth) / (halfWidth)
	y := (halfHeight - mousePositionY) / (halfHeight)
	return x, y
}
func buildVAO() []float32 {
	var vao []float32
	for _, item := range app.Points {
		vao = append(vao, float32(item.Coordinate.X))
		vao = append(vao, float32(item.Coordinate.Y))
		vao = append(vao, float32(item.Coordinate.Z))
	}
	return vao
}
func Draw() {
	if len(app.Points) < 1 {
		return
	}
	points := buildVAO()
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*3, gl.PtrOffset(4*3*len(app.Points)))

	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.POINTS, 0, int32(len(app.Points)))
}

func main() {
	app.AddPoint(primitives.PointBB{V.Vector{50, 50, 0}, V.Vector{1, 1, 1}})
	app.KeyDowns = make(map[string]bool)
	app.KeyDowns["W"] = false
	app.KeyDowns["A"] = false
	app.KeyDowns["S"] = false
	app.KeyDowns["D"] = false

	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	// - Projection matrix
	angleOfView := float64(45)
	near := float64(0.1)
	far := float64(1000)
	aspect := float64(windowWidth / windowHeight)
	ProjectionMatrix := M.Perspective(angleOfView, aspect, near, far)
	fmt.Println(ProjectionMatrix)
	// - Camera matrix
	cameraPosition := V.Vector{50, 50, 50}
	cameraLookAt := (V.Vector{50, 50, 0}).Normalize()
	worldUpDirection := V.Vector{0, 1, 0}
	camera := C.New(cameraPosition, cameraLookAt, worldUpDirection)
	CameraTransformationMatrix := camera.GetTransformation()
	fmt.Println(CameraTransformationMatrix.GetPoints())
	// - Model matrix
	translationMatrix := M.Translation(V.Vector{-50, -50, -50})
	scaleMatrix := M.Scale(V.Vector{1 / 50.0, 1 / 50.0, 1 / 50.0})
	ModelMatrix := translationMatrix.Dot(*scaleMatrix)
	fmt.Println(ModelMatrix)

	// - calculate MVP
	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	MVPMatrix := ProjectionMatrix.Dot(*(CameraTransformationMatrix.Dot(*ModelMatrix)))
	MVP := MVPMatrix.GetPoints()
	gl.UniformMatrix4fv(mvpLocation, 1, false, &MVP[0])
	fmt.Println(MVP)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		mouseHandler(window)
		keyHandler(window)
		CameraTransformationMatrix = camera.GetTransformation()
		MVPMatrix = ProjectionMatrix.Dot(*(CameraTransformationMatrix.Dot(*ModelMatrix)))
		MVP = MVPMatrix.GetPoints()
		gl.UniformMatrix4fv(mvpLocation, 1, false, &MVP[0])
		Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
