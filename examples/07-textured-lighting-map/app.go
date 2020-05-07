package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 800
	windowTitle  = "Example - textured rotating cube"

	FORWARD  = glfw.KeyW // Go forward
	BACKWARD = glfw.KeyS // Go backward
	LEFT     = glfw.KeyA // Go left
	RIGHT    = glfw.KeyD // Go right
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed            = 0.005
	rotationSpeed        = float32(2.0)
	cameraDirectionSpeed = float32(0.00500)
	CameraMoveSpeed      = 0.005

	LightSourceRoundSpeed = 3000.0
)

var (
	app  *application.Application
	cube *cuboid.Cuboid

	lastUpdate int64

	InitialCenterPointLight = mgl32.Vec3{-3, 0, -3}

	LightSource       *light.Light
	LightSourceSphere *sphere.Sphere
	JadeCube          *cuboid.Cuboid

	cameraDistance = 0.1
	rotationAngle  = float32(0.0)
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{3.3, -10, 14.0}, mgl32.Vec3{0, 1, 0}, -101.0, 21.5)
	camera.SetupProjection(45, float32(windowWidth)/float32(windowHeight), 0.1, 100.0)
	return camera
}

// It generates the lightsource sphere.
func GenerateWhiteSphere(shaderProgram *shader.Shader) {
	LightSourceSphere = sphere.New(mgl32.Vec3{-3.0, -0.5, -3.0}, mgl32.Vec3{1, 1, 1}, 0.1, shaderProgram)
	LightSourceSphere.SetPrecision(15)
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 64.0)
	LightSourceSphere.SetMaterial(mat)
	LightSourceSphere.SetDirection((mgl32.Vec3{9, 0, -3}).Normalize())
	distance := (LightSourceSphere.GetCenterPoint().Sub(JadeCube.GetCenterPoint())).Len()

	LightSourceSphere.SetSpeed((float32(2) * float32(3.1415) * distance) / LightSourceRoundSpeed)
	LightSourceSphere.DrawMode(sphere.DRAW_MODE_LIGHT)
	app.AddItem(LightSourceSphere)
}

// It generates the Jade cube.
func GenerateCube(shaderProgram *shader.Shader) {
	bottomRect := rectangle.NewSquare(mgl32.Vec3{0.5, -0.5, -0.5}, mgl32.Vec3{-0.5, -0.5, 0.5}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0.0, 1.0, 1.0}, shaderProgram)
	JadeCube = cuboid.New(bottomRect, 1.0, shaderProgram)
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 64.0)
	JadeCube.SetMaterial(mat)
	JadeCube.DrawMode(cuboid.DRAW_MODE_TEXTURED_LIGHT)
	app.AddItem(JadeCube)
}

func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	lightSourceRotationAngleRadian := mgl32.DegToRad(float32((360 / LightSourceRoundSpeed) * moveTime))
	lightDirectionRotationMatrix := mgl32.HomogRotate3D(lightSourceRotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentLightSourceDirection := LightSourceSphere.GetDirection()
	LightSourceSphere.SetDirection(mgl32.TransformNormal(currentLightSourceDirection, lightDirectionRotationMatrix))
	LightSource.SetPosition(LightSourceSphere.GetCenterPoint())

	app.Update(moveTime)

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
		vertical = -moveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = moveSpeed * moveTime
	}
	if vertical != 0 {
		app.GetCamera().Lift(float32(vertical))
	}
	currX, currY := app.GetWindow().GetCursorPos()
	x, y := trans.MouseCoordinates(currX, currY, windowWidth, windowHeight)
	KeyDowns := make(map[string]bool)
	// dUp
	if y > 1.0-cameraDistance && y < 1.0 {
		KeyDowns["dUp"] = true
	} else {
		KeyDowns["dUp"] = false
	}
	// dDown
	if y < -1.0+cameraDistance && y > -1.0 {
		KeyDowns["dDown"] = true
	} else {
		KeyDowns["dDown"] = false
	}
	// dLeft
	if x < -1.0+cameraDistance && x > -1.0 {
		KeyDowns["dLeft"] = true
	} else {
		KeyDowns["dLeft"] = false
	}
	// dRight
	if x > 1.0-cameraDistance && x < 1.0 {
		KeyDowns["dRight"] = true
	} else {
		KeyDowns["dRight"] = false
	}

	dX := float32(0.0)
	dY := float32(0.0)
	if KeyDowns["dUp"] && !KeyDowns["dDown"] {
		dY = cameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = -cameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = -cameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = cameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}

func main() {
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(windowWidth, windowHeight, windowTitle))
	defer glfw.Terminate()
	shader.InitOpenGL()

	app.SetCamera(CreateCamera())

	LightSource = light.New(InitialCenterPointLight, mgl32.Vec3{0.2, 0.2, 0.2}, mgl32.Vec3{0.5, 0.5, 0.5}, mgl32.Vec3{1, 1, 1})

	shaderProgramTexture := shader.NewShader("examples/07-textured-lighting-map/texture.vert", "examples/07-textured-lighting-map/texture.frag")
	shaderProgramTexture.AddTexture("examples/07-textured-lighting-map/image-texture-diffuse.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.diffuse")
	shaderProgramTexture.AddTexture("examples/07-textured-lighting-map/image-texture-specular.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.specular")
	shaderProgramTexture.SetLightSource(LightSource, "light.position", "light.ambient", "light.diffuse", "light.specular")
	GenerateCube(shaderProgramTexture)

	shaderProgramWhite := shader.NewShader("examples/07-textured-lighting-map/lightsource.vert", "examples/07-textured-lighting-map/lightsource.frag")
	shaderProgramWhite.SetLightSource(LightSource, "light.position", "light.ambient", "light.diffuse", "light.specular")
	GenerateWhiteSphere(shaderProgramWhite)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		Update()
		app.DrawWithUniforms()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
