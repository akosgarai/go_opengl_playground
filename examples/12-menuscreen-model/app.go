package main

import (
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
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
	WindowTitle  = "Example - menu screen"
	FontFile     = "/assets/fonts/Desyrel/desyrel.regular.ttf"

	CameraMoveSpeed      = 0.005
	CameraDirectionSpeed = float32(0.050)
	CameraDistance       = 0.1
)

var (
	app       *application.Application
	glWrapper glwrapper.Wrapper
	PaperMesh *mesh.TexturedMesh

	lastUpdate int64
)

func Paper(width, height float32) {
	rect := rectangle.NewExact(width, height)
	v, i, _ := rect.MeshInput()
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)

	PaperMesh = mesh.NewTexturedMesh(v, i, tex, glWrapper)
}

func baseDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func main() {
	runtime.LockOSThread()
	app = application.New()
	app.SetWindow(window.InitGlfw(WindowWidth, WindowHeight, WindowTitle))
	defer glfw.Terminate()
	glWrapper.InitOpenGL()

	fontShader := shader.NewShader(baseDir()+"/shaders/font.vert", baseDir()+"/shaders/font.frag", glWrapper)
	app.AddShader(fontShader)
	paperShader := shader.NewShader(baseDir()+"/shaders/paper.vert", baseDir()+"/shaders/paper.frag", glWrapper)
	app.AddShader(paperShader)

	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	glWrapper.ClearColor(0.0, 0.25, 0.5, 1.0)
	paperModel := model.New()
	Paper(1, 1)
	PaperMesh.RotateX(-90)
	PaperMesh.SetPosition(mgl32.Vec3{-0.4, -0.3, -0.0})
	paperModel.AddMesh(PaperMesh)

	app.AddModelToShader(paperModel, paperShader)

	lastUpdate = time.Now().UnixNano()
	// register keyboard button callback
	app.GetWindow().SetKeyCallback(app.KeyCallback)
	Fonts, err := model.LoadCharsetDebug(baseDir()+FontFile, 32, 127, 40.0, 72, glWrapper)
	if err != nil {
		panic(err)
	}
	Fonts.SetSurface(PaperMesh)
	Fonts.Print("How are You?", -0.5, 0.2, 3.0/float32(WindowWidth), glWrapper)
	Fonts.SetTransparent(true)
	app.AddModelToShader(Fonts, fontShader)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
