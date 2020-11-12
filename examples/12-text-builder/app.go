package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	WindowTitle = "Example - Text to meshes"
)

var (
	app *application.Application
	// window related variables
	WindowWidth      = 800
	WindowHeight     = 800
	WindowDecorated  = true
	WindowFullScreen = false
	// glwrapper
	glWrapper glwrapper.Wrapper
)

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func init() {
	runtime.LockOSThread()

	app = application.New(glWrapper)
	// Setup the window
	builder := setupWindowBuilder()
	app.SetWindow(builder.Build())
	// Init opengl.
	glWrapper.InitOpenGL()
	scrn := CreateAppScreen()
	app.AddScreen(scrn)
	app.ActivateScreen(scrn)
}
func CreateAppScreen() *screen.Screen {
	scrn := screen.New()
	scrn.SetWrapper(glWrapper)
	scrn.SetupCamera(CreateCamera(), CameraMovementOptions())
	scrn.Setup(setupApp)
	rect := rectangle.NewExact(0.8, 0.5)
	// Colored model
	col := []mgl32.Vec3{mgl32.Vec3{1, 0, 0}}
	V, I, _ := rect.ColoredMeshInput(col)
	ColoredModel := model.New()
	cm := mesh.NewColorMesh(V, I, col, scrn.GetWrapper())
	ColoredModel.AddMesh(cm)
	colorShader := shader.NewShader(baseDir()+"/shaders/color.vert", baseDir()+"/shaders/color.frag", scrn.GetWrapper())
	scrn.AddShader(colorShader)
	scrn.AddModelToShader(ColoredModel, colorShader)
	// Directional lightsource is necessary for the materials.
	// directional light is coming from the up direction but not from too up.
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{0.5, 0.5, 0.5},
		mgl32.Vec3{0.5, 0.5, 0.5},
		mgl32.Vec3{0.5, 0.5, 0.5},
	})
	scrn.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	// Material model
	mat := material.Gold
	V, I, _ = rect.MeshInput()
	MaterialModel := model.New()
	mm := mesh.NewMaterialMesh(V, I, mat, scrn.GetWrapper())
	MaterialModel.AddMesh(mm)
	matShader := shader.NewMaterialShader(scrn.GetWrapper())
	scrn.AddShader(matShader)
	scrn.AddModelToShader(MaterialModel, matShader)
	// Textured model
	var normalTex texture.Textures
	normalTex.AddTexture(baseDir()+"/../../assets/sample-texture.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	normalTex.AddTexture(baseDir()+"/../../assets/sample-texture.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", glWrapper)
	PlaneTexModel := model.New()
	tm := mesh.NewTexturedMesh(V, I, normalTex, scrn.GetWrapper())
	PlaneTexModel.AddMesh(tm)
	texShader := shader.NewTextureShader(scrn.GetWrapper())
	scrn.AddShader(texShader)
	scrn.AddModelToShader(PlaneTexModel, texShader)
	// Colored Texture model
	ColorTexModel := model.New()
	V, I, _ = rect.TexturedColoredMeshInput(col)
	ctm := mesh.NewTexturedColoredMesh(V, I, normalTex, col, scrn.GetWrapper())
	ColorTexModel.AddMesh(ctm)
	colorTexShader := shader.NewShader(baseDir()+"/shaders/color-texture.vert", baseDir()+"/shaders/color-texture.frag", scrn.GetWrapper())
	scrn.AddShader(colorTexShader)
	scrn.AddModelToShader(ColorTexModel, colorTexShader)
	// Material Texture model
	MaterialTexModel := model.New()
	V, I, _ = rect.MeshInput()
	mtm := mesh.NewTexturedMaterialMesh(V, I, normalTex, mat, scrn.GetWrapper())
	MaterialTexModel.AddMesh(mtm)
	matTexShader := shader.NewTextureMatShaderBlending(scrn.GetWrapper())
	scrn.AddShader(matTexShader)
	scrn.AddModelToShader(MaterialTexModel, matTexShader)
	// Material transparent Texture model
	MaterialTransTexModel := model.New()
	var transTex texture.Textures
	transTex.TransparentTexture(1, 1, 1, "tex.diffuse", scrn.GetWrapper())
	transTex.TransparentTexture(1, 1, 1, "tex.specular", scrn.GetWrapper())
	mttm := mesh.NewTexturedMaterialMesh(V, I, transTex, mat, scrn.GetWrapper())
	MaterialTransTexModel.AddMesh(mttm)
	matTransTexShader := shader.NewTextureMatShaderBlending(scrn.GetWrapper())
	scrn.AddShader(matTransTexShader)
	scrn.AddModelToShader(MaterialTransTexModel, matTransTexShader)
	surfaceMeshes := []interfaces.Mesh{cm, mm, tm, ctm, mtm, mttm}
	// rotation of the meshes
	for i := 0; i < len(surfaceMeshes); i++ {
		surfaceMeshes[i].RotateZ(-90)
		surfaceMeshes[i].RotateX(90)
		y := float32(1-i/2) * 0.6
		z := 0.5 - float32(i%2)
		surfaceMeshes[i].SetPosition(mgl32.Vec3{1.7274, y, z})
	}
	// charset + text printing to each kind of surfaces.
	MenuFonts, err := model.LoadCharset(baseDir()+"/../../assets/fonts/Desyrel/desyrel.regular.ttf", 32, 127, 40.0, 300, scrn.GetWrapper())
	if err != nil {
		panic(err)
	}
	MenuFonts.SetTransparent(true)
	fontShader := shader.NewFontShader(scrn.GetWrapper())
	scrn.AddShader(fontShader)
	scrn.AddModelToShader(MenuFonts, fontShader)
	fontColor := []mgl32.Vec3{mgl32.Vec3{1, 1, 1}}
	surfaceTexts := []string{"Color", "Material", "Texture", "Col-Tex", "Mat-Tex", "Mat-Tex-Blending"}
	scale := float32(1.0) / float32(3000.0)
	for i := 0; i < len(surfaceMeshes); i++ {
		w := MenuFonts.TextWidth(surfaceTexts[i], scale)
		MenuFonts.PrintTo(surfaceTexts[i], -w/2, 0, -0.003, scale, scrn.GetWrapper(), surfaceMeshes[i], fontColor)
	}
	return scrn
}
func CreateCamera() interfaces.Camera {
	cam := camera.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0.0, 0.0)
	cam.SetupProjection(45, float32(WindowWidth)/float32(WindowHeight), 0.001, 10)
	fmt.Printf("%s\n", cam.Log())
	return cam
}

// Setup options for the camera
func CameraMovementOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["mode"] = "default"
	cm["rotateOnEdgeDistance"] = float32(0.0)
	return cm
}
func setupApp(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.ClearColor(0.0, 1.0, 0.0, 1.0)
	glWrapper.Viewport(0, 0, int32(WindowWidth), int32(WindowHeight))
}
func setupWindowBuilder() *window.WindowBuilder {
	Builder := window.NewWindowBuilder()
	fullScreen := os.Getenv("FULL")
	if fullScreen == "1" {
		WindowFullScreen = true
		WindowDecorated = false
		WindowWidth, WindowHeight = Builder.GetCurrentMonitorResolution()
	} else {
		width := os.Getenv("WIDTH")
		if width != "" {
			val, err := strconv.Atoi(width)
			if err == nil {
				WindowWidth = val
			}
		}
		height := os.Getenv("HEIGHT")
		if height != "" {
			val, err := strconv.Atoi(height)
			if err == nil {
				WindowHeight = val
			}
		}
		decorated := os.Getenv("DECORATED")
		if decorated == "0" {
			WindowDecorated = false
		}
	}
	Builder.SetFullScreen(WindowFullScreen)
	Builder.SetDecorated(WindowDecorated)
	Builder.SetTitle(WindowTitle)
	Builder.SetWindowSize(WindowWidth, WindowHeight)
	return Builder
}

func main() {
	// Terminate window at the end.
	defer glfw.Terminate()
	// main event loop
	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw(glWrapper)
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
