package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/texture"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - textured spheres"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed            = 0.005
	cameraDirectionSpeed = float32(0.010)
	CameraMoveSpeed      = 0.005
	cameraDistance       = 0.1

	EarthRoundSpeed = 3000.0
	SunRoundSpeed   = 0.01
)

var (
	app              *application.Application
	Sun              *mesh.TexturedMesh
	Earth            *mesh.TexturedMesh
	lastUpdate       int64
	PointLightSource *light.Light

	rotationAngle      = float32(0.0)
	PointLightAmbient  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition = mgl32.Vec3{0.0, 0.0, 0.0}
	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.14)
	LightQuadraticTerm = float32(0.07)
	spherePrimitive    = sphere.New(20)

	glWrapper wrapper.Wrapper
)

func TexturedSphere(t texture.Textures, position mgl32.Vec3, scale float32, shaderProgram *shader.Shader) *mesh.TexturedMesh {
	v, i := spherePrimitive.TexturedMeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(position)
	m.SetScale(mgl32.Vec3{scale, scale, scale})
	return m
}

// It generates a cube map.
func CubeMap(t texture.Textures) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	return m
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -10.0}, mgl32.Vec3{0, 1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.01, 200.0)
	return camera
}
func updateSun(moveTime float64) {
	rotationAngle = rotationAngle + float32(moveTime)*SunRoundSpeed
	Sun.SetRotationAngle(mgl32.DegToRad(rotationAngle))
}
func updateEarth(moveTime float64) {
	// Calculate the  rotation matrix. Get the current one, rotate it with a calculated angle around the Y axis. (HomogRotate3D(angle float32, axis Vec3) Mat4)
	// angle calculation: (360 / LightSourceRoundSpeed) * delta) -> in radian: mat32.DegToRad()
	// Then we can transform the current direction vector to the new one. (TransformNormal(v Vec3, m Mat4) Vec3)
	// after it we can set the new direction vector of the light source.
	rotationAngleRadian := mgl32.DegToRad(float32((360 / EarthRoundSpeed) * moveTime))
	rotationMatrix := mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})
	currentDirection := Earth.GetDirection()
	Earth.SetDirection(mgl32.TransformNormal(currentDirection, rotationMatrix))
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	updateEarth(moveTime)
	updateSun(moveTime)

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
	horizontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horizontal = -moveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horizontal = moveSpeed * moveTime
	}
	if horizontal != 0 {
		app.GetCamera().Strafe(float32(horizontal))
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
	x, y := trans.MouseCoordinates(currX, currY, WindowWidth, WindowHeight)
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
		dY = -cameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = cameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = cameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = -cameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	PointLightSource = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})

	// Add the lightources to the application
	app.AddPointLightSource(PointLightSource, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader("examples/07-textured-spheres/shaders/texture.vert", "examples/07-textured-spheres/shaders/texture.frag", glWrapper)
	app.AddShader(shaderProgramTexture)

	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture("examples/07-textured-spheres/assets/sun.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture("examples/07-textured-spheres/assets/sun.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)
	Sun = TexturedSphere(sunTexture, mgl32.Vec3{0.0, 0.0, 0.0}, 1, shaderProgramTexture)
	Sun.SetRotationAxis(mgl32.Vec3{0.0, -1.0, 0.0})
	app.AddMeshToShader(Sun, shaderProgramTexture)
	// sun texture
	var earthTexture texture.Textures
	earthTexture.AddTexture("examples/07-textured-spheres/assets/earth.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	earthTexture.AddTexture("examples/07-textured-spheres/assets/earth.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)
	Earth = TexturedSphere(earthTexture, mgl32.Vec3{3.0, 0.0, 0.0}, 0.1, shaderProgramTexture)
	distance := Earth.GetPosition().Len()
	Earth.SetSpeed((float32(2) * float32(3.1415) * distance) / EarthRoundSpeed)
	Earth.SetDirection((mgl32.Vec3{0, 0, 1}).Normalize())
	app.AddMeshToShader(Earth, shaderProgramTexture)

	shaderProgramCubeMap := shader.NewShader("examples/07-textured-spheres/shaders/cubeMap.vert", "examples/07-textured-spheres/shaders/cubeMap.frag", glWrapper)
	app.AddShader(shaderProgramCubeMap)
	var cubeMapTexture texture.Textures
	cubeMapTexture.AddCubeMapTexture("examples/07-textured-spheres/assets", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "skybox", glWrapper)
	cubeMap := CubeMap(cubeMapTexture)
	cubeMap.SetScale(mgl32.Vec3{100.0, 100.0, 100.0})
	app.AddMeshToShader(cubeMap, shaderProgramCubeMap)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.0, 0.0, 0.0, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
