package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/primitives"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
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
	WindowTitle  = "Example - mesh experiment"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	moveSpeed            = 0.005
	cameraDirectionSpeed = float32(0.05)
	CameraMoveSpeed      = 0.005
	cameraDistance       = 0.1

	rotationSpeed = float32(2.0)
)

var (
	app *application.Application

	RotatingCube *mesh.TexturedMesh
	LiftingCube  *mesh.TexturedMesh

	lastUpdate int64

	DirectionalLightDirection = (mgl32.Vec3{0.7, 0.7, 0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightDiffuse   = mgl32.Vec3{0.1, 0.1, 0.1}
	DirectionalLightSpecular  = mgl32.Vec3{0.1, 0.1, 0.1}
	PointLightAmbient         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse         = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular        = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition_1      = mgl32.Vec3{8, -0.5, -1.0}
	LightConstantTerm         = float32(1.0)
	LightLinearTerm           = float32(0.14)
	LightQuadraticTerm        = float32(0.07)
	SpotLightAmbient          = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse          = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular         = mgl32.Vec3{1, 1, 1}
	SpotLightDirection_1      = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightPosition_1       = mgl32.Vec3{0.20, -6, -0.7}
	SpotLightCutoff_1         = float32(4)
	SpotLightOuterCutoff_1    = float32(5)

	rotationAngle = float32(0.0)

	glWrapper wrapper.Wrapper
)

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{-11.5, -4.8, 9.5}, mgl32.Vec3{0, 1, 0}, -37.0, -2.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}

func GenerateGrassMesh(t texture.Textures) *mesh.TexturedMesh {
	square := primitives.NewSquare()
	v, i := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetScale(mgl32.Vec3{100, 1, 100})
	return m
}
func GenerateCubeMesh(t texture.Textures, pos mgl32.Vec3) *mesh.TexturedMesh {
	cube := primitives.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(pos)
	return m
}
func GenerateRotatingCubeMesh(t texture.Textures, pos mgl32.Vec3) *mesh.TexturedMesh {
	cube := primitives.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(pos)
	m.SetRotationAxis(mgl32.Vec3{0, 1, 0})
	return m
}
func GenerateLiftingCubeMesh(t texture.Textures, pos mgl32.Vec3) *mesh.TexturedMesh {
	cube := primitives.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetPosition(pos)
	m.SetDirection(mgl32.Vec3{0, 1, 0})
	m.SetSpeed(moveSpeed)
	return m
}
func GenerateMaterialCubeMesh(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	cube := primitives.NewCube()
	v, i := cube.MeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	return m
}
func GenerateMaterialSphereMesh(mat *material.Material, pos mgl32.Vec3) *mesh.MaterialMesh {
	sphere := primitives.NewSphere(20)
	v, i := sphere.MaterialMeshInput()
	m := mesh.NewMaterialMesh(v, i, mat, glWrapper)
	m.SetPosition(pos)
	return m
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	rotationAngle = rotationAngle + float32(moveTime)*rotationSpeed
	RotatingCube.SetRotationAngle(mgl32.DegToRad(mgl32.DegToRad(rotationAngle)))

	if LiftingCube.GetPosition().Y() < -10 || LiftingCube.GetPosition().Y() > -0.5 {
		LiftingCube.SetDirection(LiftingCube.GetDirection().Mul(-1))
	}

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

	textureShader := shader.NewShader("examples/model-loading/shaders/texture.vert", "examples/model-loading/shaders/texture.frag", glWrapper)
	app.AddShader(textureShader)
	materialShader := shader.NewShader("examples/model-loading/shaders/material.vert", "examples/model-loading/shaders/material.frag", glWrapper)
	app.AddShader(materialShader)

	var TexturesGrass texture.Textures
	TexturesGrass.AddTexture("examples/model-loading/assets/grass.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	TexturesGrass.AddTexture("examples/model-loading/assets/grass.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)
	var TexturesCube texture.Textures
	TexturesCube.AddTexture("examples/model-loading/assets/texture-diffuse.png", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", glWrapper)
	TexturesCube.AddTexture("examples/model-loading/assets/texture-specular.png", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.specular", glWrapper)

	grassMesh := GenerateGrassMesh(TexturesGrass)
	app.AddMeshToShader(grassMesh, textureShader)
	cubeMesh := GenerateCubeMesh(TexturesCube, mgl32.Vec3{0, -0.5, 0})
	app.AddMeshToShader(cubeMesh, textureShader)
	RotatingCube = GenerateRotatingCubeMesh(TexturesCube, mgl32.Vec3{3, -0.5, 0})
	app.AddMeshToShader(RotatingCube, textureShader)
	LiftingCube = GenerateLiftingCubeMesh(TexturesCube, mgl32.Vec3{2, -10, 3})
	app.AddMeshToShader(LiftingCube, textureShader)
	materialCube := GenerateMaterialCubeMesh(material.Silver, mgl32.Vec3{3, -0.5, -3})
	app.AddMeshToShader(materialCube, materialShader)
	materialSphere := GenerateMaterialSphereMesh(material.Silver, mgl32.Vec3{-3, -0.5, 3})
	app.AddMeshToShader(materialSphere, materialShader)

	// setup lighsources.
	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	PointLightSource_1 := light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition_1,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	SpotLightSource_1 := light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition_1,
		SpotLightDirection_1,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_1, SpotLightOuterCutoff_1})
	app.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	app.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	app.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	glWrapper.ClearColor(0.0, 0.0, 0.0, 1.0)

	lastUpdate = time.Now().UnixNano()
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
