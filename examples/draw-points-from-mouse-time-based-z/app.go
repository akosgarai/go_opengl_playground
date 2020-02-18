package main

import (
	"fmt"
	"runtime"

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
	modelSize          = 100
)

type Application struct {
	Points []primitives.Point
}

// AddPoint inserts a new point to the points.
func (a *Application) AddPoint(point primitives.Point) {
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
			positionVector := primitives.Vector{x, y, 1.0}
			scalarPart := float64(modelSize / 2)
			modelPositionVector := (positionVector.MultiplyScalar(scalarPart)).AddScalar(scalarPart)
			app.AddPoint(
				primitives.Point{
					modelPositionVector,
					primitives.Vector{colorR, colorG, colorB},
				})
		}
	}
}

// Key handler function. it supports the debug option. (print out the points of the app)
func keyHandler(window *glfw.Window) {
	if window.GetKey(glfw.KeyD) == glfw.Press {
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
	for _, item := range app.Points {
		vao = append(vao, float32(item.Color.X))
		vao = append(vao, float32(item.Color.Y))
		vao = append(vao, float32(item.Color.Z))
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
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()

	gl.UseProgram(program)

	gl.Enable(gl.PROGRAM_POINT_SIZE)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Viewport(0, 0, windowWidth, windowHeight)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	// based on the following tutorial: http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
	// I have to define and multiply the following matrixes:
	// - Projection 'glm::mat4 Projection = glm::perspective(glm::radians(45.0f), (float) width / (float)height, 0.1f, 100.0f);'
	angleOfView := float64(180)
	near := float64(0.1)
	far := float64(100)
	ProjectionMatrix := primitives.ProjectionMatrix4x4(angleOfView, near, far)
	fmt.Println(ProjectionMatrix)
	// - Camera matrix
	/* based on the following doc: https://community.khronos.org/t/about-glm-lookat-function/74506
	 */
	cameraPosition := primitives.Vector{50, 50, 50}
	cameraLookAt := primitives.Vector{50, 50, 0}
	worldUpDirection := primitives.Vector{0, 1, 0}
	zAxis := (cameraPosition.Subtract(cameraLookAt)).Normalize()
	xAxis := ((worldUpDirection.Normalize()).Cross(zAxis)).Normalize()
	yAxis := zAxis.Cross(xAxis)
	cameraTranslationMatrix := primitives.TranslationMatrixT4x4(
		-float32(cameraPosition.X),
		-float32(cameraPosition.Y),
		-float32(cameraPosition.Z))
	cameraRotationMatrix := primitives.UnitMatrix4x4()
	cameraRotationMatrix.Points[0] = float32(xAxis.X)
	cameraRotationMatrix.Points[4] = float32(xAxis.Y)
	cameraRotationMatrix.Points[8] = float32(xAxis.Z)
	cameraRotationMatrix.Points[1] = float32(yAxis.X)
	cameraRotationMatrix.Points[5] = float32(yAxis.Y)
	cameraRotationMatrix.Points[9] = float32(yAxis.Z)
	cameraRotationMatrix.Points[2] = float32(zAxis.X)
	cameraRotationMatrix.Points[6] = float32(zAxis.Y)
	cameraRotationMatrix.Points[10] = float32(zAxis.Z)
	CameraMatrix := cameraRotationMatrix.Dot(cameraTranslationMatrix)
	fmt.Println(CameraMatrix)
	// - Model matrix (the scale)
	translationMatrix := primitives.TranslationMatrixT4x4(-50, -50, -50)
	scaleMatrix := primitives.ScaleMatrix4x4(1/50.0, 1/50.0, 1/50.0)
	ModelMatrix := translationMatrix.Dot(scaleMatrix)
	// mvp = projection * camera * model
	MVP := ((ProjectionMatrix.Dot(CameraMatrix)).Dot(ModelMatrix)).Points
	mvpLocation := gl.GetUniformLocation(program, gl.Str("MVP\x00"))
	gl.UniformMatrix4fv(mvpLocation, 1, false, &MVP[0])
	fmt.Println(MVP)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		keyHandler(window)
		mouseHandler(window)
		Draw()
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
