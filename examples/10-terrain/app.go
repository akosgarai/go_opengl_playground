package main

import (
	"fmt"
	"math/rand"
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
	WaterModel                = model.New()
	lastUpdate                int64
	startTime                 int64
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

func generateHeightMap(width, length, iterations, peakProbability int, minH, maxH float32, seed int64) [][]float32 {
	// init map with 0.0-s
	var heightMap [][]float32
	for w := 0; w <= width; w++ {
		heightMap = append(heightMap, []float32{})
		for l := 0; l <= length; l++ {
			heightMap[w] = append(heightMap[w], 0.0)
		}
	}
	terrainMaxDiff := maxH - minH
	iterationStep := terrainMaxDiff / float32(iterations)

	rand.Seed(seed)
	fmt.Printf("Seed: %d\n", seed)
	for i := 0; i < iterations; i++ {
		value := minH + float32(i)*iterationStep
		fmt.Printf("Value: %f\n", value)
		for l := 0; l <= length; l++ {
			for w := 0; w <= width; w++ {
				if heightMap[l][w] != 0 {
					continue
				}

				rndNum := rand.Intn(100)
				//fmt.Printf("Random: %d\n", rndNum)
				if adjacentElevation(w, l, value-iterationStep, peakProbability, width, length, heightMap) || rndNum < peakProbability {
					heightMap[l][w] = value
				}
			}
		}
	}
	//fmt.Printf("HeightMap: %v\n", heightMap)
	return heightMap
}

func adjacentElevation(w, h int, elevation float32, cliffProbability, width, height int, elements [][]float32) bool {
	for y := max(0, h-1); y <= min(height-1, h+1); y++ {
		for x := max(0, w-1); x <= min(width-1, w+1); x++ {
			if elements[y][x] == elevation {
				// if this element is *not* randomly a cliff, return true
				return rand.Intn(100) > cliffProbability
			}
		}
	}

	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if b > a {
		return a
	}
	return b
}

// It should create a terrain surface. A flat surface is on the x-z plane.
// In the x axis, it is from -x/2 to x/2. In the z axis it is from -length/2, lenth/2.
// In the y axis, it is from minH to maxX.
func NewTerrain(width, length, iterations int, minH, maxH float32, seed int64, t texture.Textures) *mesh.TexturedMesh {

	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
	var vertices vertex.Vertices
	var indices []uint32
	// generate heights.
	peakProbability := 5
	heightMap := generateHeightMap(width, length, iterations, peakProbability, minH, maxH, seed)

	for l := 0; l <= length; l++ {
		for w := 0; w <= width; w++ {
			texIndex := (w % 2) + (l%2)*2
			// The normal vector calculation is necessary to be implemented.
			// How can I calculate the normal vector? it should be based on the height map.
			// Get the heights in the neighbor points. Get the direction vectors and then
			// it could be calculated.
			var iL, iW int
			if l == length {
				iL = l
			} else {
				iL = l + 1
			}
			if w == width {
				iW = w
			} else {
				iW = w + 1
			}
			currentPos := mgl32.Vec3{-float32(width)/2.0 + float32(w), heightMap[l][w], -float32(length)/2.0 + float32(l)}
			nextPosX := mgl32.Vec3{-float32(width)/2.0 + float32(iW), heightMap[l][iW], -float32(length)/2.0 + float32(l)}
			nextPosY := mgl32.Vec3{-float32(width)/2.0 + float32(w), heightMap[iL][w], -float32(length)/2.0 + float32(iL)}
			normal := nextPosX.Sub(currentPos).Cross(nextPosY.Sub(currentPos)).Normalize()
			vertices = append(vertices, vertex.Vertex{
				Position:  currentPos,
				Normal:    normal,
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
	//fmt.Println(vertices)
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
	camera := camera.NewCamera(mgl32.Vec3{-13.0, 6.5, 2.5}, mgl32.Vec3{0, -1, 0}, -3.5, -16.5)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 20.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}
func Update() {
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	app.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
	app.Update(delta)
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
func CreateWaterMesh() *mesh.TexturedMesh {
	var waterTexture texture.Textures
	waterTexture.AddTexture(baseDir()+"/assets/water.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	waterTexture.AddTexture(baseDir()+"/assets/water.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	square := rectangle.NewSquare()
	v, i, _ := square.MeshInput()
	m := mesh.NewTexturedMesh(v, i, waterTexture, glWrapper)
	m.SetScale(mgl32.Vec3{20, 1, 20})
	m.SetPosition(mgl32.Vec3{0.0, 1.0, 0.0})
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
	shaderProgramTexture := shader.NewTextureShaderBlending(glWrapper)
	app.AddShader(shaderProgramTexture)
	shaderProgramWater := shader.NewShader(baseDir()+"/shaders/water.vert", baseDir()+"/shaders/water.frag", glWrapper)
	app.AddShader(shaderProgramWater)

	var grassTexture texture.Textures
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	grassTexture.AddTexture(baseDir()+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	grassMesh := NewTerrain(4, 4, 10, -1, 3, 0, grassTexture)
	grassMesh.SetPosition(mgl32.Vec3{0.0, 1.003, 0.0})
	grassMesh.SetScale(mgl32.Vec3{5, 1, 5})
	TexModel.AddMesh(grassMesh)
	app.AddModelToShader(TexModel, shaderProgramTexture)

	WaterModel.AddMesh(CreateWaterMesh())
	WaterModel.SetTransparent(true)
	app.AddModelToShader(WaterModel, shaderProgramWater)
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
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)

	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate
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
