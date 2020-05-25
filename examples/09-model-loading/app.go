package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/model"
	"github.com/akosgarai/opengl_playground/pkg/modelimport"
	"github.com/akosgarai/opengl_playground/pkg/primitives/camera"
	"github.com/akosgarai/opengl_playground/pkg/primitives/light"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Model loading example"

	FORWARD  = glfw.KeyW
	BACKWARD = glfw.KeyS
	LEFT     = glfw.KeyA
	RIGHT    = glfw.KeyD
	UP       = glfw.KeyQ
	DOWN     = glfw.KeyE

	CameraMoveSpeed       = 0.005
	CameraDirectionSpeed  = float32(0.00500)
	DefaultModelDirectory = "examples/09-model-loading/assets"
	DefaultModelFilename  = "object.obj"
)

var (
	app      *application.Application
	Importer *modelimport.Import

	lastUpdate int64

	cameraDistance = 0.1

	glWrapper wrapper.Wrapper

	Model                 = model.New()
	PointModel            = model.New()
	MaterialModel         = model.New()
	TexturedMaterialModel = model.New()
	TexturedColorModel    = model.New()

	DirectionalLightDirection = (mgl32.Vec3{0.7, 0.7, 0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{1.0, 1.0, 1.0}
	DirectionalLightDiffuse   = mgl32.Vec3{1.0, 1.0, 1.0}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}

	PointLightAmbient  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightDiffuse  = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightSpecular = mgl32.Vec3{0.5, 0.5, 0.5}
	PointLightPosition = mgl32.Vec3{8, -0.5, -1.0}

	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.14)
	LightQuadraticTerm = float32(0.07)

	SpotLightAmbient     = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse     = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular    = mgl32.Vec3{1, 1, 1}
	SpotLightDirection   = (mgl32.Vec3{0, 1, 0}).Normalize()
	SpotLightPosition    = mgl32.Vec3{0.20, -6, -0.7}
	SpotLightCutoff      = float32(4)
	SpotLightOuterCutoff = float32(5)

	DirectionalLightSource = light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	PointLightSource = light.NewPointLight([4]mgl32.Vec3{
		PointLightPosition,
		PointLightAmbient,
		PointLightDiffuse,
		PointLightSpecular},
		[3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	SpotLightSource = light.NewSpotLight([5]mgl32.Vec3{
		SpotLightPosition,
		SpotLightDirection,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff, SpotLightOuterCutoff})

	SingleDirectory    = true
	MultiDirectoryName = ""
)

// Importer init. If we have 2 or more command line arguments,
// the first one is used as model directory, the second one
// as the model filename.
func Init() {
	args := os.Args[1:]
	if len(args) == 0 {
		// load default model
		Importer = modelimport.New(DefaultModelDirectory, DefaultModelFilename)
	} else if len(args) == 1 {
		// multi model mode. read every subdir. (exported directory handling)
		SingleDirectory = false
		MultiDirectoryName = args[0]
	} else if len(args) > 1 {
		// load the directory with the given filename.
		Importer = modelimport.New(args[0], args[1])
	}
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -5.0}, mgl32.Vec3{0, 1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	app.Update(moveTime)

	forward := 0.0
	if app.GetKeyState(FORWARD) && !app.GetKeyState(BACKWARD) {
		forward = CameraMoveSpeed * moveTime
	} else if app.GetKeyState(BACKWARD) && !app.GetKeyState(FORWARD) {
		forward = -CameraMoveSpeed * moveTime
	}
	if forward != 0 {
		app.GetCamera().Walk(float32(forward))
	}
	horizontal := 0.0
	if app.GetKeyState(LEFT) && !app.GetKeyState(RIGHT) {
		horizontal = -CameraMoveSpeed * moveTime
	} else if app.GetKeyState(RIGHT) && !app.GetKeyState(LEFT) {
		horizontal = CameraMoveSpeed * moveTime
	}
	if horizontal != 0 {
		app.GetCamera().Strafe(float32(horizontal))
	}
	vertical := 0.0
	if app.GetKeyState(UP) && !app.GetKeyState(DOWN) {
		vertical = -CameraMoveSpeed * moveTime
	} else if app.GetKeyState(DOWN) && !app.GetKeyState(UP) {
		vertical = CameraMoveSpeed * moveTime
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
		dY = CameraDirectionSpeed
	} else if KeyDowns["dDown"] && !KeyDowns["dUp"] {
		dY = -CameraDirectionSpeed
	}
	if KeyDowns["dLeft"] && !KeyDowns["dRight"] {
		dX = -CameraDirectionSpeed
	} else if KeyDowns["dRight"] && !KeyDowns["dLeft"] {
		dX = CameraDirectionSpeed
	}
	app.GetCamera().UpdateDirection(dX, dY)
}
func AddMeshToRightModel(m interfaces.Mesh) {
	switch m.(type) {
	case *mesh.TexturedMaterialMesh:
		TexturedMaterialModel.AddMesh(m)
	case *mesh.TexturedColoredMesh:
		TexturedColorModel.AddMesh(m)
	case *mesh.MaterialMesh:
		MaterialModel.AddMesh(m)
	case *mesh.PointMesh:
		PointModel.AddMesh(m)
	}
}
func main() {
	Init()
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())

	pointShader := shader.NewShader("examples/09-model-loading/shaders/point.vert", "examples/09-model-loading/shaders/point.frag", glWrapper)
	app.AddShader(pointShader)
	materialShader := shader.NewShader("examples/09-model-loading/shaders/material.vert", "examples/09-model-loading/shaders/material.frag", glWrapper)
	app.AddShader(materialShader)
	texColShader := shader.NewShader("examples/09-model-loading/shaders/texturecolor.vert", "examples/09-model-loading/shaders/texturecolor.frag", glWrapper)
	app.AddShader(texColShader)
	texMatShader := shader.NewShader("examples/09-model-loading/shaders/texturemat.vert", "examples/09-model-loading/shaders/texturemat.frag", glWrapper)
	app.AddShader(texMatShader)
	app.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	app.AddPointLightSource(PointLightSource, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	app.AddSpotLightSource(SpotLightSource, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})

	if SingleDirectory {
		Importer.Import()
		meshes := Importer.GetMeshes()
		for _, m := range meshes {
			AddMeshToRightModel(m)
		}
	} else {
		files, err := ioutil.ReadDir(MultiDirectoryName)
		if err != nil {
			fmt.Println(err)
		}
		for _, f := range files {
			Importer = modelimport.New(MultiDirectoryName+"/"+f.Name(), DefaultModelFilename)
			Importer.Import()
			meshes := Importer.GetMeshes()
			for _, m := range meshes {
				AddMeshToRightModel(m)
			}
		}
	}
	app.AddModelToShader(TexturedMaterialModel, texMatShader)
	app.AddModelToShader(TexturedColorModel, texColShader)
	app.AddModelToShader(MaterialModel, materialShader)
	app.AddModelToShader(PointModel, pointShader)

	glWrapper.Enable(wrapper.DEPTH_TEST)
	glWrapper.DepthFunc(wrapper.LESS)
	lastUpdate = time.Now().UnixNano()
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(wrapper.COLOR_BUFFER_BIT | wrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw()
		app.GetWindow().SwapBuffers()
	}
}
