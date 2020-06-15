package main

import (
	"fmt"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	app                       *application.Application
	TexModel                  = model.New()
	lastUpdate                int64
	DirectionalLightSource    *light.Light
	DirectionalLightDirection = (mgl32.Vec3{0.5, 0.5, -0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{0.5, 0.5, 0.5}

	glWrapper glwrapper.Wrapper
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - Terrain"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
)

type Terrain struct {
	width     int
	length    int
	iteration int
	maxHeight float32
	minHeight float32
	seed      int64
	Model     model.BaseModel
}

// It should create a terrain surface. A flat surface is on the x-z plane.
// In the x axis, it is from -x/2 to x/2. In the z axis it is from -length/2, lenth/2.
// In the y axis, it is from minH to maxX.
func NewTerrain(width, length, iterations int, minH, maxH float32, seed int64, t texture.Textures) *mesh.TexturedMesh {

	//terrainMaxDiff := maxH - minH
	//iteraxitionStep := terrainMaxDiff / float32(iterations)

	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}

	var vertices vertex.Vertices
	var indices []uint32

	for w := 0; w <= width; w++ {
		for l := 0; l <= length; l++ {
			texIndex := (w % 2) + (l%2)*2
			vertices = append(vertices, vertex.Vertex{
				Position:  mgl32.Vec3{-float32(width)/2.0 + float32(w), 0, -float32(length)/2.0 + float32(l)},
				Normal:    mgl32.Vec3{0, 1, 0},
				TexCoords: textureCoords[texIndex],
			})
		}
	}
	for w := 0; w <= width-1; w++ {
		for l := 0; l <= length-1; l++ {
			i0 := uint32(w*(length+1) + l)
			i1 := uint32(1) + i0
			i2 := uint32(length+1) + i0
			i3 := uint32(1) + i2
			indices = append(indices, i0)
			indices = append(indices, i1)
			indices = append(indices, i2)

			indices = append(indices, i2)
			indices = append(indices, i1)
			indices = append(indices, i3)
		}
	}
	return mesh.NewTexturedMesh(vertices, indices, t, glWrapper)
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
	camera := camera.NewCamera(mgl32.Vec3{0.0, -0.5, 3.0}, mgl32.Vec3{0, 1, 0}, 0.0, 0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 20.0)
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

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func CreateGrassMesh(t texture.Textures) *mesh.TexturedMesh {
	square := rectangle.NewSquare()
	v, i, bo := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, t, glWrapper)
	m.SetScale(mgl32.Vec3{20, 1, 20})
	m.SetBoundingObject(bo)
	return m
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	app.SetCamera(CreateCamera())
	app.SetCameraMovementMap(CameraMovementMap())
	app.SetRotateOnEdgeDistance(CameraDistance)

	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShader(glWrapper)
	app.AddShader(shaderProgramTexture)

	var grassTexture texture.Textures
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	grassMesh := NewTerrain(20, 20, 1, 0, 0, 0, grassTexture)
	grassMesh.SetPosition(mgl32.Vec3{0.0, 1.003, 0.0})
	TexModel.AddMesh(grassMesh)
	app.AddModelToShader(TexModel, shaderProgramTexture)

	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource = light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	// Add the lightources to the application
	app.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
