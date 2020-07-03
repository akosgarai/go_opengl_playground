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
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Example - rooms with light sources"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
	FontFile             = "/assets/fonts/Desyrel/desyrel.ttf"
	LEFT_MOUSE_BUTTON    = glfw.MouseButtonLeft
	MENU_BUTTON          = glfw.KeyM
	TOGGLE_DOOR_BUTTON   = glfw.KeyC
)

var (
	MenuApp                *application.Application
	RoomApp                *application.Application
	ActiveApp              *application.Application
	Ground                 *model.Terrain
	Water                  *model.Liquid
	Room                   *model.Room
	Room2                  *model.Room
	Room3                  *model.Room
	Menu                   *model.BaseModel
	StartButton            *mesh.TexturedMaterialMesh
	ExitButton             *mesh.TexturedMaterialMesh
	lastUpdate             int64
	startTime              int64
	DirectionalLightSource *light.Light
	PointLightSource_1     *light.Light
	PointLightSource_2     *light.Light
	SpotLightSource_1      *light.Light
	SpotLightSource_2      *light.Light

	TexModel                  = model.New()
	DirectionalLightDirection = (mgl32.Vec3{0.5, 0.5, -0.7}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightDiffuse   = mgl32.Vec3{0.3, 0.3, 0.3}
	DirectionalLightSpecular  = mgl32.Vec3{0.3, 0.3, 0.3}
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

	DefaultMaterial   = material.Jade
	HighlightMaterial = material.Ruby

	glWrapper glwrapper.Wrapper
)

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

func CreateGround() {
	gb := model.NewTerrainBuilder()
	gb.SetWidth(4)
	gb.SetLength(4)
	gb.SetIterations(10)
	gb.SetScale(mgl32.Vec3{5, 1, 5})
	gb.SetGlWrapper(glWrapper)
	gb.SurfaceTextureGrass()
	gb.LiquidTextureWater()
	gb.SetPeakProbability(5)
	gb.SetCliffProbability(5)
	gb.SetMinHeight(-1)
	gb.SetMaxHeight(3)
	gb.SetPosition(mgl32.Vec3{0.0, 1.003, 0.0})
	gb.SetSeed(0)
	gb.SetLiquidEta(0.75)
	gb.SetLiquidAmplitude(0.125 / 2.0)
	gb.SetLiquidFrequency(1.0)
	gb.SetLiquidDetailMultiplier(10)
	gb.SetLiquidWaterLevel(0.25)
	gb.SetDebugMode(false)
	Ground, Water = gb.BuildWithLiquid()
	Water.SetTransparent(true)
}

// It creates a new camera with the necessary setup
func CreateCamera() *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0.0, -0.5, 3.0}, mgl32.Vec3{0, 1, 0}, -85.0, -0.0)
	camera.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 20.0)
	camera.SetVelocity(CameraMoveSpeed)
	camera.SetRotationStep(CameraDirectionSpeed)
	return camera
}

// Update has to be separated to model update and to game update.
// - Let's handle the delta first.
// - Then the MENU_BUTTON handler.
// - GetClosestModel
// - Switch on model first. If it is identical to menu, then compare the meshed with exit, start and do the action
// - If it is not menu, then we have the game handler.
func Update() {
	// Delta
	nowNano := time.Now().UnixNano()
	delta := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	ActiveApp.SetUniformFloat("time", float32(float64(nowNano-startTime)/float64(time.Second)))
	ActiveApp.Update(delta)
	// MENU_BUTTON is pressed.
	if ActiveApp.GetKeyState(MENU_BUTTON) {
		ActiveApp = MenuApp
		glWrapper.ClearColor(0.0, 1.0, 0.0, 1.0)
		StartButton.Material = DefaultMaterial
	}
	// Get closest stuff
	mdl, msh, distance := ActiveApp.GetClosestModelMeshDistance()
	switch mdl.(type) {
	case *model.BaseModel:
		baseModel := mdl.(*model.BaseModel)
		// Check only the menu options.
		if baseModel == Menu {
			tmMesh := msh.(*mesh.TexturedMaterialMesh)
			if distance < 0.01 {
				tmMesh.Material = HighlightMaterial
				if ActiveApp.GetMouseButtonState(LEFT_MOUSE_BUTTON) {
					if tmMesh == ExitButton {
						fmt.Println("Exit button has been pressed.\n")
						ActiveApp.GetWindow().SetShouldClose(true)
					} else if tmMesh == StartButton {
						fmt.Println("Start button has been pressed.\n")
						ActiveApp = RoomApp
						ActiveApp.GetWindow().SetKeyCallback(ActiveApp.KeyCallback)
						glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
					}
				}
			} else {
				tmMesh.Material = DefaultMaterial
			}
		}
		break
	case *model.Room:
		if ActiveApp.GetKeyState(TOGGLE_DOOR_BUTTON) {
			room := mdl.(*model.Room)
			room.PushDoorState()
		}
		break
	}
}
func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func GroundPos(pos mgl32.Vec3) mgl32.Vec3 {
	h, err := Ground.HeightAtPos(pos)
	if err == nil {
		return mgl32.Vec3{pos.X(), h - 0.001, pos.Z()}
	}
	return pos
}

func Paper(width, height float32, position mgl32.Vec3) *mesh.TexturedMaterialMesh {
	rect := rectangle.NewExact(width, height)
	v, i, bo := rect.MeshInput()
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)

	msh := mesh.NewTexturedMaterialMesh(v, i, tex, DefaultMaterial, glWrapper)
	msh.SetBoundingObject(bo)
	msh.SetPosition(position)
	return msh
}

func main() {
	runtime.LockOSThread()
	MenuApp = application.New()
	RoomApp = application.New()
	Window := window.InitGlfw(WindowWidth, WindowHeight, WindowTitle)
	MenuApp.SetWindow(Window)
	RoomApp.SetWindow(Window)
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	RoomApp.SetCamera(CreateCamera())
	RoomApp.SetCameraMovementMap(CameraMovementMap())
	RoomApp.SetRotateOnEdgeDistance(CameraDistance)
	RoomApp.SetUniformFloat("fog.minDistance", 1.0)
	RoomApp.SetUniformFloat("fog.maxDistance", 8.0)
	RoomApp.SetUniformVector("fog.color", mgl32.Vec3{0.4, 0.4, 0.4})

	// Shader application for the material objects
	shaderProgramMaterial := shader.NewMaterialShaderWithFog(glWrapper)
	RoomApp.AddShader(shaderProgramMaterial)
	// Shader application for the textured meshes.
	shaderProgramTexture := shader.NewTextureShaderBlendingWithFog(glWrapper)
	RoomApp.AddShader(shaderProgramTexture)
	// Shader application for the blending & fog
	shaderProgramTextureMat := shader.NewTextureMatShaderBlendingWithFog(glWrapper)
	RoomApp.AddShader(shaderProgramTextureMat)
	// Shader application for the liquid surface
	shaderProgramLiquid := shader.NewTextureShaderLiquidWithFog(glWrapper)
	RoomApp.AddShader(shaderProgramLiquid)

	CreateGround()
	RoomApp.AddModelToShader(Ground, shaderProgramTexture)
	RoomApp.AddModelToShader(Water, shaderProgramLiquid)

	Room = model.NewMaterialRoom(GroundPos(mgl32.Vec3{-5.0, 1.0, 0.0}), glWrapper)
	RoomApp.AddModelToShader(Room, shaderProgramMaterial)
	Room2 = model.NewMaterialRoom(GroundPos(mgl32.Vec3{-6.0, 1.0, 6.0}), glWrapper)
	RoomApp.AddModelToShader(Room2, shaderProgramMaterial)
	Room3 = model.NewTextureRoom(GroundPos(mgl32.Vec3{-3.5, 1.0, 5.7}), glWrapper)
	RoomApp.AddModelToShader(Room3, shaderProgramTextureMat)

	streetLamp := model.NewMaterialStreetLamp(GroundPos(mgl32.Vec3{-6.6, -2.0, 1.3}), 1.3, glWrapper)
	streetLamp.RotateX(90)
	streetLamp.RotateY(-90)
	RoomApp.AddModelToShader(streetLamp, shaderProgramMaterial)
	streetLampDark := model.NewTexturedStreetLamp(GroundPos(mgl32.Vec3{-4.8, -2.0, 6.9}), 1.3, glWrapper)
	streetLampDark.RotateX(90)
	streetLampDark.RotateY(-90)
	RoomApp.AddModelToShader(streetLampDark, shaderProgramTextureMat)

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
		streetLamp.GetBulbPosition(),
		SpotLightDirection_1,
		SpotLightAmbient,
		SpotLightDiffuse,
		SpotLightSpecular},
		[5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, SpotLightCutoff_1, SpotLightOuterCutoff_1})
	SpotLightSource_2 = light.NewSpotLight([5]mgl32.Vec3{
		streetLampDark.GetBulbPosition(),
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
	RoomApp.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	//app.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	//app.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	RoomApp.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	RoomApp.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)

	lastUpdate = time.Now().UnixNano()
	startTime = lastUpdate
	// register keyboard button callback
	RoomApp.GetWindow().SetKeyCallback(RoomApp.KeyCallback)

	// Setup menu application
	fontShader := shader.NewShader(baseDir()+"/shaders/font.vert", baseDir()+"/shaders/font.frag", glWrapper)
	MenuApp.AddShader(fontShader)
	paperShader := shader.NewShader(baseDir()+"/shaders/paper.vert", baseDir()+"/shaders/paper.frag", glWrapper)
	MenuApp.AddShader(paperShader)

	Menu = model.New()
	StartButton = Paper(1, 0.2, mgl32.Vec3{-0.0, 0.3, -0.0})
	StartButton.RotateX(-90)
	Menu.AddMesh(StartButton)
	ExitButton = Paper(1, 0.2, mgl32.Vec3{-0.0, -0.3, -0.0})
	ExitButton.RotateX(-90)
	Menu.AddMesh(ExitButton)
	MenuApp.AddModelToShader(Menu, paperShader)
	cols1 := []mgl32.Vec3{
		mgl32.Vec3{0.0, 1.0, 0.0},
	}
	cols2 := []mgl32.Vec3{
		mgl32.Vec3{0.0, 0.0, 1.0},
	}
	MenuFonts, err := model.LoadCharset(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	MenuFonts.PrintTo(" - Start - ", -0.4, -0.03, 0.01, 3.0/float32(WindowWidth), glWrapper, StartButton, cols1)
	MenuFonts.PrintTo(" - Exit - ", -0.4, -0.03, 0.01, 3.0/float32(WindowWidth), glWrapper, ExitButton, cols2)
	MenuFonts.SetTransparent(true)
	MenuApp.AddModelToShader(MenuFonts, fontShader)
	// register keyboard button callback
	MenuApp.GetWindow().SetMouseButtonCallback(MenuApp.MouseButtonCallback)
	MenuApp.GetWindow().SetKeyCallback(window.DummyKeyCallback)

	ActiveApp = MenuApp

	for !ActiveApp.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		Update()
		ActiveApp.Draw()
		glfw.PollEvents()
		ActiveApp.GetWindow().SwapBuffers()
	}
}
