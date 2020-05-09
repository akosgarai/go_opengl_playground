package main

import (
	"runtime"
	"time"

	"github.com/akosgarai/opengl_playground/pkg/application"
	"github.com/akosgarai/opengl_playground/pkg/composite/bug"
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
	Bug1                      *bug.Bug
	Bug2                      *bug.Bug
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
	SpotLightPosition_1       = mgl32.Vec3{0.20, -6, -0.7}
	SpotLightPosition_2       = mgl32.Vec3{10.20, -6, -0.7}
	SpotLightCutoff_1         = float32(4)
	SpotLightCutoff_2         = float32(4)
	SpotLightOuterCutoff_1    = float32(5)
	SpotLightOuterCutoff_2    = float32(5)
)

// It generates a square.
func Grass(shaderProgram *shader.Shader) {
	square := rectangle.NewSquare(mgl32.Vec3{-1000, 0, -1000}, mgl32.Vec3{1000, 0, 1000}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0}, shaderProgram)
	square.SetPrecision(1)
	square.DrawMode(cuboid.DRAW_MODE_TEXTURED_LIGHT)
	app.AddItem(square)
}

// It generates the box.
func Box(shaderProgram *shader.Shader) {
	numberOfBoxes := 3
	baseX := float32(-6)
	sideLength := float32(2)
	baseY := float32(-1)
	diffBetweenBoxes := float32(3)
	for i := 0; i < numberOfBoxes; i++ {
		x1_pos := baseX + float32(i)*(sideLength+diffBetweenBoxes)
		bottomRect := rectangle.NewSquare(mgl32.Vec3{x1_pos + sideLength, 0, baseY}, mgl32.Vec3{x1_pos, 0, baseY + sideLength}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0.0, 1.0, 1.0}, shaderProgram)
		box := cuboid.New(bottomRect, sideLength, shaderProgram)
		mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 64.0)
		box.SetMaterial(mat)
		box.DrawMode(cuboid.DRAW_MODE_TEXTURED_LIGHT)
		app.AddItem(box)
	}
}

// It generates the lamp
func Lamp(shaderProgram *shader.Shader, baseX, baseZ float32, lightPosition mgl32.Vec3) {
	width := float32(0.4)
	height := float32(6)
	length := float32(2.5)

	bottomRect := rectangle.NewSquare(mgl32.Vec3{baseX, 0, baseZ}, mgl32.Vec3{baseX + width, 0, baseZ + width}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 0, 0}, shaderProgram)
	pole := cuboid.New(bottomRect, height, shaderProgram)
	pole.SetMaterial(material.Chrome)
	pole.SetPrecision(10)
	pole.DrawMode(cuboid.DRAW_MODE_LIGHT)
	app.AddItem(pole)
	bottomRect = rectangle.NewSquare(mgl32.Vec3{baseX, -height, baseZ}, mgl32.Vec3{baseX + width, -height - width, baseZ}, mgl32.Vec3{0, 0, -1}, mgl32.Vec3{0, 0, 0}, shaderProgram)
	top := cuboid.New(bottomRect, length, shaderProgram)
	top.SetMaterial(material.Chrome)
	top.DrawMode(cuboid.DRAW_MODE_LIGHT)
	top.SetPrecision(10)
	app.AddItem(top)
	bulb := sphere.New(lightPosition, mgl32.Vec3{1, 1, 1}, float32(0.1), shaderProgram)
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 256.0)
	bulb.SetMaterial(mat)
	bulb.SetPrecision(15)
	bulb.DrawMode(cuboid.DRAW_MODE_LIGHT)
	app.AddItem(bulb)
}

// It generates the bug. This stuff will be the point light source. In the furure, i want to make it fly around.
func Bug(shaderProgram *shader.Shader) {
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 128.0)
	Bug1 = bug.Firefly(PointLightPosition_1, 0.1, [3]*material.Material{mat, material.Blackplastic, material.Ruby}, shaderProgram)
	Bug1.SetDirection(mgl32.Vec3{1, 0, 0})
	Bug1.SetSpeed(moveSpeed)
	app.AddItem(Bug1)
}
func BigBug(shaderProgram *shader.Shader) {
	mat := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 128.0)
	Bug2 = bug.Firefly(PointLightPosition_2, 1.0, [3]*material.Material{mat, material.Blackplastic, material.Ruby}, shaderProgram)
	Bug2.SetDirection(mgl32.Vec3{0, 0, 1})
	Bug2.SetSpeed(moveSpeed)
	app.AddItem(Bug2)
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
		rotationAngleRadian := mgl32.DegToRad(-45)
		directionRotationMatrix := mgl32.HomogRotate3D(rotationAngleRadian, mgl32.Vec3{0, -1, 0})

		currenDirection := Bug1.GetDirection()
		Bug1.SetDirection(mgl32.TransformNormal(currenDirection, directionRotationMatrix))
	}
}
func Update() {
	nowNano := time.Now().UnixNano()
	moveTime := float64(nowNano-lastUpdate) / float64(time.Millisecond)
	lastUpdate = nowNano
	for index, _ := range ShaderProgramsWithViewPos {
		ShaderProgramsWithViewPos[index].SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	}
	RotateBugOne(nowNano)
	PointLightSource_1.SetPosition(Bug1.GetCenterPoint())
	PointLightSource_2.SetPosition(Bug2.GetCenterPoint())
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
	shader.InitOpenGL()

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
	//Define the shader application for the grass
	shaderProgramGrass := shader.NewShader("examples/08-multiple-light/shaders/texture.vert", "examples/08-multiple-light/shaders/texture.frag")
	shaderProgramGrass.AddTexture("examples/08-multiple-light/assets/grass.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.diffuse")
	shaderProgramGrass.AddTexture("examples/08-multiple-light/assets/grass.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.specular")
	shaderProgramGrass.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	shaderProgramGrass.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	shaderProgramGrass.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	shaderProgramGrass.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	shaderProgramGrass.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})
	shaderProgramGrass.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	ShaderProgramsWithViewPos = append(ShaderProgramsWithViewPos, shaderProgramGrass)
	Grass(shaderProgramGrass)
	// Shader application for the box
	shaderProgramBox := shader.NewShader("examples/08-multiple-light/shaders/texture.vert", "examples/08-multiple-light/shaders/texture.frag")
	shaderProgramBox.AddTexture("examples/08-multiple-light/assets/box-diffuse.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.diffuse")
	shaderProgramBox.AddTexture("examples/08-multiple-light/assets/box-specular.png", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "material.specular")
	shaderProgramBox.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	shaderProgramBox.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	shaderProgramBox.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	shaderProgramBox.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	shaderProgramBox.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})
	shaderProgramBox.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	ShaderProgramsWithViewPos = append(ShaderProgramsWithViewPos, shaderProgramBox)
	Box(shaderProgramBox)
	// Shader application for the lamp
	shaderProgramLamp := shader.NewShader("examples/08-multiple-light/shaders/lamp.vert", "examples/08-multiple-light/shaders/lamp.frag")
	shaderProgramLamp.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	shaderProgramLamp.AddPointLightSource(PointLightSource_1, [7]string{"pointLight[0].position", "pointLight[0].ambient", "pointLight[0].diffuse", "pointLight[0].specular", "pointLight[0].constant", "pointLight[0].linear", "pointLight[0].quadratic"})
	shaderProgramLamp.AddPointLightSource(PointLightSource_2, [7]string{"pointLight[1].position", "pointLight[1].ambient", "pointLight[1].diffuse", "pointLight[1].specular", "pointLight[1].constant", "pointLight[1].linear", "pointLight[1].quadratic"})
	shaderProgramLamp.AddSpotLightSource(SpotLightSource_1, [10]string{"spotLight[0].position", "spotLight[0].direction", "spotLight[0].ambient", "spotLight[0].diffuse", "spotLight[0].specular", "spotLight[0].constant", "spotLight[0].linear", "spotLight[0].quadratic", "spotLight[0].cutOff", "spotLight[0].outerCutOff"})
	shaderProgramLamp.AddSpotLightSource(SpotLightSource_2, [10]string{"spotLight[1].position", "spotLight[1].direction", "spotLight[1].ambient", "spotLight[1].diffuse", "spotLight[1].specular", "spotLight[1].constant", "spotLight[1].linear", "spotLight[1].quadratic", "spotLight[1].cutOff", "spotLight[1].outerCutOff"})
	shaderProgramLamp.SetViewPosition(app.GetCamera().GetPosition(), "viewPosition")
	ShaderProgramsWithViewPos = append(ShaderProgramsWithViewPos, shaderProgramLamp)
	Lamp(shaderProgramLamp, 0, -3, SpotLightPosition_1)
	Lamp(shaderProgramLamp, 10, -3, SpotLightPosition_2)
	Bug(shaderProgramLamp)
	BigBug(shaderProgramLamp)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	lastUpdate = time.Now().UnixNano()
	BugOneLastRotate = lastUpdate
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
