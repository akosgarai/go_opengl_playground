package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
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
	WindowTitle  = "Example - multiple light source"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed            = 0.005
	cameraDirectionSpeed = float32(0.100)
	CameraMoveSpeed      = 0.005
	cameraDistance       = 0.1
)

var (
	app                       *application.Application
	Bug1                      *model.Bug
	Bug2                      *mesh.TexturedMesh
	BugOneLastRotate          int64
	lastUpdate                int64
	ShaderProgramsWithViewPos []*shader.Shader
	DirectionalLightSource    *light.Light
	PointLightSource_1        *light.Light
	PointLightSource_2        *light.Light
	SpotLightSource_1         *light.Light
	SpotLightSource_2         *light.Light

	BugOneForwardMove         = float64(1000)
	DirectionalLightDirection = (mgl32.Vec3{0.7, 0.7, 0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightDiffuse   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightSpecular  = mgl32.Vec3{0.1, 0.1, 0.1}
	PointLightAmbient         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular        = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition_1      = mgl32.Vec3{8, -0.5, -1.0}
	PointLightPosition_2      = mgl32.Vec3{8, -5, -30}
	LightConstantTerm         = float32(1.0)
	LightLinearTerm           = float32(0.14)
	LightQuadraticTerm        = float32(0.07)
	SpotLightAmbient          = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse          = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular         = mgl32.Vec3{1, 1, 1}
	SpotLightDirection_1      = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightDirection_2      = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightPosition_1       = mgl32.Vec3{0.20, -6, -0.65}
	SpotLightPosition_2       = mgl32.Vec3{10.20, -6, -0.65}
	SpotLightCutoff_1         = float32(4)
	SpotLightCutoff_2         = float32(4)
	SpotLightOuterCutoff_1    = float32(5)
	SpotLightOuterCutoff_2    = float32(5)
	TexModel                  = model.New()
	MatModel                  = model.New()

	glWrapper wrapper.Wrapper
)

func CreateGrassMesh(t texture.Textures) *mesh.TexturedMesh {
	square := rectangle.NewSquare()
	v, i := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetScale(mgl32.Vec3{1000, 1, 1000})
	return m
}
func CreateCubeMesh(t texture.Textures, pos mgl32.Vec3) *mesh.TexturedMesh {
	cube := cuboid.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(pos)
	return m
}

// It generates the lamp. Now it usws the StreetLamp model for creating it.
func StreetLamp(position mgl32.Vec3) *model.StreetLamp {
	StreetLamp := model.NewStreetLamp(position)
	StreetLamp.Rotate(180, mgl32.Vec3{1, 0, 0})
	StreetLamp.Rotate(90, mgl32.Vec3{0, 1, 0})
	return StreetLamp
}

func TexturedBug(t texture.Textures) {
	sph := sphere.New(15)
	v, i := sph.TexturedMeshInput()
	Bug2 = mesh.NewTexturedMesh(v, i, t, glWrapper)
	Bug2.SetPosition(PointLightPosition_2)
	Bug2.SetDirection(mgl32.Vec3{0, 0, 1})
	Bug2.SetSpeed(moveSpeed)
	TexModel.AddMesh(Bug2)
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-11.2, -5.0, 4.2}, mgl32.Vec3{0, 1, 0}, -37.0, -2.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 1000.0)
	return camera
}
func RotateBugOne(now int64) {
	moveTime := float64(now-BugOneLastRotate) / float64(time.Millisecond)
	if moveTime > BugOneForwardMove {
		BugOneLastRotate = now
		// rotate 45 deg
		Bug1.Rotate(-45, mgl32.Vec3{0, 1, 0}.Normalize())
	}
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	RotateBugOne(nowNano)
	PointLightSource_1.SetPosition(Bug1.GetBottomPosition())
	PointLightSource_2.SetPosition(Bug2.GetPosition())
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

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource = light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	PointLightSource_1 = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition_1,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	SpotLightSource_1 = light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition_1,
		SpotLightDirection_1,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_1, SpotLightOuterCutoff_1})
	SpotLightSource_2 = light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition_2,
		SpotLightDirection_2,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_2, SpotLightOuterCutoff_2})
	PointLightSource_2 = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition_2,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})

	// Add the lightources to the application
	app.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	app.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	app.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	app.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	app.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})

	// Define the shader application for the textured meshes.
	shaderProgramTexture := shader.NewShader("examples/08-multiple-light/shaders/texture.vert", "examples/08-multiple-light/shaders/texture.frag", glWrapper)
	app.AddShader(shaderProgramTexture)

	// grass textures
	var grassTexture texture.Textures
	grassTexture.AddTexture("examples/08-multiple-light/assets/grass.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	grassTexture.AddTexture("examples/08-multiple-light/assets/grass.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)

	grassMesh := CreateGrassMesh(grassTexture)
	TexModel.AddMesh(grassMesh)

	// box textures
	var boxTexture texture.Textures
	boxTexture.AddTexture("examples/08-multiple-light/assets/box-diffuse.png", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	boxTexture.AddTexture("examples/08-multiple-light/assets/box-specular.png", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)

	// we have 3 boxes in the following coordinates.
	boxPositions := []mgl32.Vec3{
		mgl32.Vec3{-5.0, -0.51, 0.0},
		mgl32.Vec3{0.0, -0.51, 0.0},
		mgl32.Vec3{5.0, -0.51, 0.0},
	}
	for _, pos := range boxPositions {
		box := CreateCubeMesh(boxTexture, pos)
		TexModel.AddMesh(box)
	}

	// Shader application for the lamp
	shaderProgramMaterial := shader.NewShader("examples/08-multiple-light/shaders/lamp.vert", "examples/08-multiple-light/shaders/lamp.frag", glWrapper)
	app.AddShader(shaderProgramMaterial)

	lamp1 := StreetLamp(mgl32.Vec3{0.4, -12, -1.3})
	app.AddModelToShader(lamp1, shaderProgramMaterial)
	lamp2 := StreetLamp(mgl32.Vec3{10.4, -12, -1.3})
	app.AddModelToShader(lamp2, shaderProgramMaterial)

	Bug1 = model.NewBug(mgl32.Vec3{9, -0.5, -1.0}, mgl32.Vec3{0.2, 0.2, 0.2})
	Bug1.SetDirection(mgl32.Vec3{1, 0, 0})
	Bug1.SetSpeed(moveSpeed)

	app.AddModelToShader(Bug1, shaderProgramMaterial)

	// sun texture
	var sunTexture texture.Textures
	sunTexture.AddTexture("examples/08-multiple-light/assets/sun.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	sunTexture.AddTexture("examples/08-multiple-light/assets/sun.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)
	TexturedBug(sunTexture)
	app.AddModelToShader(TexModel, shaderProgramTexture)
	app.AddModelToShader(MatModel, shaderProgramMaterial)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.0, 0.0, 0.0, 1.0)

	lastUpdate = time.Now().UnixNano()
	BugOneLastRotate = lastUpdate
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
