package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/callbacks/primitives"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/shader"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - plane with ball"
)

type Application struct {
	window  *glfw.Window
	program uint32
	// camera related parameters
	camera               *camera.Camera
	cameraDirection      float64
	cameraDirectionSpeed float32
	cameraLastUpdate     int64
	moveSpeed            float64
	epsilon              float64

	square                *primitives.Square
	sphere                *primitives.Sphere
	sphereDirectionVector mgl32.Vec3 // currently it only has to be up or down.

	KeyDowns map[string]bool
}

// It generates a Sphere.
func (a *Application) GenerateSphere() {
	a.sphere = primitives.NewSphere()
	a.sphere.SetCenter(mgl32.Vec3{0, -5, 0})
	a.sphere.SetColor(mgl32.Vec3{1, 0, 0})
	a.sphere.SetRadius(2.0)
}

// MoveSphere updates the sphere position.
func (a *Application) MoveSphere(delta float64) {
	moveVector := a.sphereDirectionVector.Mul(float32(delta * a.moveSpeed / 5.0))
	newCenter := a.sphere.GetCenter().Add(moveVector)
	if a.sphereDirectionVector.Y() == -1.0 {
		if newCenter.Y() < -5 {
			newCenter = mgl32.Vec3{newCenter.X(), -5, newCenter.Z()}
			a.sphereDirectionVector = a.sphereDirectionVector.Mul(-1.0)
		}
	} else {
		if newCenter.Y() > -2 {
			newCenter = mgl32.Vec3{newCenter.X(), -2, newCenter.Z()}
			a.sphereDirectionVector = a.sphereDirectionVector.Mul(-1.0)
		}
	}
	a.sphere.SetCenter(newCenter)
}

// It generates a square.
func (a *Application) GenerateSquare() {
	squareColor := mgl32.Vec3{0, 1, 0}
	a.square = primitives.NewSquare(
		primitives.Point{mgl32.Vec3{-20, 0, -20}, squareColor},
		primitives.Point{mgl32.Vec3{20, 0, -20}, squareColor},
		primitives.Point{mgl32.Vec3{20, 0, 20}, squareColor},
		primitives.Point{mgl32.Vec3{-20, 0, 20}, squareColor})
	a.square.SetPrecision(10)
}

func NewApplication() *Application {
	var app Application

	app.moveSpeed = 1.0 / 1000.0
	app.epsilon = 50.0

	app.camera = camera.NewCamera(
		mgl32.Vec3{-10, -4, 22.0},
		mgl32.Vec3{0, 1, 0},
		300.0,
		16.0)
	app.camera.SetupProjection(
		45,
		float32(windowWidth)/float32(windowHeight),
		0.1,
		100.0)
	app.cameraDirection = 0.1
	app.cameraDirectionSpeed = 5
	app.cameraLastUpdate = time.Now().UnixNano()

	// objects
	app.GenerateSquare()
	app.GenerateSphere()
	app.sphereDirectionVector = mgl32.Vec3{0, -1, 0}

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

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case glfw.KeyH:
		if action != glfw.Release {
			fmt.Printf("app.camera: %s\n", a.camera.Log())
			fmt.Printf("app.square: %v\n", a.square)
			fmt.Printf("app.sphere: %v\n", a.sphere)
		}
		break
	case glfw.KeyW:
		if action != glfw.Release {
			a.KeyDowns["W"] = true
		} else {
			a.KeyDowns["W"] = false
		}
		break
	case glfw.KeyA:
		if action != glfw.Release {
			a.KeyDowns["A"] = true
		} else {
			a.KeyDowns["A"] = false
		}
		break
	case glfw.KeyS:
		if action != glfw.Release {
			a.KeyDowns["S"] = true
		} else {
			a.KeyDowns["S"] = false
		}
		break
	case glfw.KeyD:
		if action != glfw.Release {
			a.KeyDowns["D"] = true
		} else {
			a.KeyDowns["D"] = false
		}
		break
	case glfw.KeyQ:
		if action != glfw.Release {
			a.KeyDowns["Q"] = true
		} else {
			a.KeyDowns["Q"] = false
		}
		break
	case glfw.KeyE:
		if action != glfw.Release {
			a.KeyDowns["E"] = true
		} else {
			a.KeyDowns["E"] = false
		}
		break
	}
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("MouseButtonCallback has been called with the following options: button: '%d', action: '%d'!, mods: '%d'\n", button, action, mods)
}

// Update is responsible for the camera movement.
func (a *Application) Update() {
	//calculate delta
	nowUnix := time.Now().UnixNano()
	delta := nowUnix - a.cameraLastUpdate
	moveTime := float64(delta / int64(time.Millisecond))
	// Update the ball's position. It's center has to be < -2, > -5. It's current move direction also has to be known.
	a.MoveSphere(moveTime)
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

// Draw is responsible for drawing the screen.
func (a *Application) Draw() {
	a.Update()
	V := a.camera.GetViewMatrix()
	P := a.camera.GetProjectionMatrix()
	a.sphere.DrawWithUniforms(V, P)
	a.square.DrawWithUniforms(V, P)
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
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
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	gl.Viewport(0, 0, windowWidth, windowHeight)
	return program
}

func main() {
	runtime.LockOSThread()

	app := NewApplication()

	app.window = initGlfw()
	defer glfw.Terminate()
	app.program = initOpenGL()
	app.square.SetShaderProgram(app.program)
	app.sphere.SetShaderProgram(app.program)

	gl.UseProgram(app.program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// register keyboard button callback
	app.window.SetKeyCallback(app.KeyCallback)
	// register mouse button callback
	app.window.SetMouseButtonCallback(app.MouseButtonCallback)

	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	for !app.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.window.SwapBuffers()
	}
}
