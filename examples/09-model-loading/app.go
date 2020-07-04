package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/modelimport"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Model loading example"

	CameraMoveSpeed       = 0.005
	CameraDirectionSpeed  = float32(0.00500)
	CameraDistance        = 0.1
	DefaultModelDirectory = "examples/09-model-loading/assets"
	DefaultModelFilename  = "object.obj"
)

var (
	app      *application.Application
	Importer *modelimport.Import

	lastUpdate int64

	glWrapper glwrapper.Wrapper

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
		Importer = modelimport.New(DefaultModelDirectory, DefaultModelFilename, glWrapper)
	} else if len(args) == 1 {
		// multi model mode. read every subdir. (exported directory handling)
		SingleDirectory = false
		MultiDirectoryName = args[0]
	} else if len(args) > 1 {
		// load the directory with the given filename.
		Importer = modelimport.New(args[0], args[1], glWrapper)
	}
}

// Setup keymap for the camera movement
func CameraMovementMap() map[string]glfw.Key {
	cm := make(map[string]glfw.Key)
	cm["forward"] = glfw.KeyW
	cm["back"] = glfw.KeyS
	cm["up"] = glfw.KeyQ
	cm["down"] = glfw.KeyE
	cm["left"] = glfw.KeyA
	cm["right"] = glfw.KeyD
	return cm
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, 0.0, -5.0}, mgl32.Vec3{0, 1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.1, 100.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano

	app.Update(delta)
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
		pointMesh := m.(*mesh.PointMesh)
		for index, _ := range pointMesh.Vertices {
			pointMesh.Vertices[index].PointSize = float32(3 + rand.Intn(17))
			pointMesh.Vertices[index].Color = mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
		}
		PointModel.AddMesh(pointMesh)
	}
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
}

func main() {
	Init()
	runtime.LockOSThread()

	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	scrn := screen.New()
	scrn.SetCamera(CreateCamera())
	scrn.SetCameraMovementMap(CameraMovementMap())
	scrn.SetRotateOnEdgeDistance(CameraDistance)

	pointShader := shader.NewShader(baseDir()+"/shaders/point.vert", baseDir()+"/shaders/point.frag", glWrapper)
	scrn.AddShader(pointShader)
	materialShader := shader.NewShader(baseDir()+"/shaders/material.vert", baseDir()+"/shaders/material.frag", glWrapper)
	scrn.AddShader(materialShader)
	texColShader := shader.NewShader(baseDir()+"/shaders/texturecolor.vert", baseDir()+"/shaders/texturecolor.frag", glWrapper)
	scrn.AddShader(texColShader)
	texMatShader := shader.NewShader(baseDir()+"/shaders/texturemat.vert", baseDir()+"/shaders/texturemat.frag", glWrapper)
	scrn.AddShader(texMatShader)
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	scrn.AddPointLightSource(PointLightSource, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	scrn.AddSpotLightSource(SpotLightSource, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})

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
			Importer = modelimport.New(MultiDirectoryName+"/"+f.Name(), DefaultModelFilename, glWrapper)
			Importer.Import()
			meshes := Importer.GetMeshes()
			for _, m := range meshes {
				AddMeshToRightModel(m)
			}
		}
	}
	scrn.AddModelToShader(TexturedMaterialModel, texMatShader)
	scrn.AddModelToShader(TexturedColorModel, texColShader)
	scrn.AddModelToShader(MaterialModel, materialShader)
	scrn.AddModelToShader(PointModel, pointShader)
	scrn.Setup(setupApp)
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)

	lastUpdate = time.Now().UnixNano()
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		glfw.PollEvents()
		Update()
		app.Draw(glWrapper)
		app.GetWindow().SwapBuffers()
	}
}
