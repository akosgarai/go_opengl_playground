package main

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/application"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/window"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
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

type Glyph struct {
	glyphWidth    int
	glyphHeight   int
	bearingHeight int
	bearingWidth  int
	advance       int
	tex           texture.Textures
}

type Charset struct {
	*model.BaseModel
	fonts               map[rune]*Glyph
	maxWidth, maxHeight int
}

func Paper(width, height float32) {
	rect := rectangle.NewExact(width, height)
	v, i, _ := rect.MeshInput()
	var tex texture.Textures
	tex.AddTexture(baseDir()+"/assets/paper.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "paper", glWrapper)

	PaperMesh = mesh.NewTexturedMesh(v, i, tex, glWrapper)
}
func (c *Charset) Print(text string, x, y, scale float32) {
	indices := []rune(text)
	fmt.Printf("The following text will be printed: '%s' as '%v'\n", text, indices)
	if len(indices) == 0 {
		return
	}
	// the low rune value from the LoadCharset function.
	lc := rune(32)
	var mshStore []interfaces.Mesh
	for i := range indices {
		runeIndex := indices[i]
		//skip runes that are not in font chacter range
		if int(runeIndex)-int(lc) > len(c.fonts) || runeIndex < lc {
			fmt.Printf("%c %d\n", runeIndex, runeIndex)
			continue
		}
		ch := c.fonts[runeIndex]
		//calculate position and size for current rune
		xpos := x + float32(ch.bearingWidth)*scale
		ypos := y + float32(ch.glyphHeight-ch.bearingHeight)*scale
		w := float32(ch.glyphWidth) * scale
		h := float32(ch.glyphHeight) * scale
		rect := rectangle.NewExact(w, h)
		cols := []mgl32.Vec3{
			mgl32.Vec3{0.0, 1.0, 0.0},
		}
		v, i, _ := rect.TexturedColoredMeshInput(cols)
		rotTr := PaperMesh.RotationTransformation()
		position := mgl32.Vec3{x + float32(ch.bearingWidth+ch.glyphWidth/2)*scale, 0.01, y - float32(ch.bearingHeight-ch.glyphHeight/2)*scale}
		msh := mesh.NewTexturedColoredMesh(v, i, ch.tex, cols, glWrapper)
		msh.SetPosition(mgl32.TransformCoordinate(position, rotTr))
		msh.SetParent(PaperMesh)
		mshStore = append(mshStore, msh)
		fmt.Printf("%c %d\n\tpos: %#v\n\tch: %#v\n\tw: %f, h: %f, xpos: %f, ypos: %f, adv: %f\n", runeIndex, runeIndex, position, ch, w, h, xpos, ypos, float32(ch.advance)*scale)
		x += float32(ch.advance) * scale
	}
	for i := len(mshStore) - 1; i >= 0; i-- {
		c.Model.AddMesh(mshStore[i])
	}
}

func LoadCharset(fontFile string, low, high rune, scale float64) *Charset {
	filePath := baseDir() + fontFile
	fd, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}
	// Read the truetype font.
	ttf, err := truetype.Parse(data)
	if err != nil {
		panic(err)
	}
	fonts := make(map[rune]*Glyph)
	var maxWidth, maxHeight int
	for ch := low; ch <= high; ch++ {
		//create new face to measure glyph diamensions
		ttfFace := truetype.NewFace(ttf, &truetype.Options{
			Size:    scale,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		gBnd, gAdv, ok := ttfFace.GlyphBounds(ch)
		if ok != true {
			panic("NOTOK")
		}

		gh := int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)
		gw := int32((gBnd.Max.X - gBnd.Min.X) >> 6)
		//if gylph has no diamensions set to a max value
		if gw == 0 || gh == 0 {
			gBnd = ttf.Bounds(fixed.Int26_6(scale))
			gw = int32((gBnd.Max.X - gBnd.Min.X) >> 6)
			gh = int32((gBnd.Max.Y - gBnd.Min.Y) >> 6)

			//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if gw == 0 || gh == 0 {
				gw = 1
				gh = 1
			}
		}
		//The glyph's ascent and descent equal -bounds.Min.Y and +bounds.Max.Y.
		gAscent := int(-gBnd.Min.Y) >> 6
		gdescent := int(gBnd.Max.Y) >> 6
		//fmt.Printf("gAscent: %d, gdescent: %d\n", gAscent, gdescent)

		//create image to draw glyph
		fg, bg := image.White, image.Black
		rect := image.Rect(0, 0, int(gw), int(gh))
		rgba := image.NewRGBA(rect)
		draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
		//create a freetype context for drawing
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(ttf)
		c.SetFontSize(scale)
		c.SetClip(rgba.Bounds())
		c.SetDst(rgba)
		c.SetSrc(fg)
		c.SetHinting(font.HintingFull)
		//set the glyph dot
		px := 0 - (int(gBnd.Min.X) >> 6)
		py := (gAscent)
		pt := freetype.Pt(px, py)

		// Draw the text from mask to image
		_, err = c.DrawString(string(ch), pt)
		if err != nil {
			panic(err)
		}

		// Generate texture
		g := &Glyph{
			glyphWidth:    int(gw),
			glyphHeight:   int(gh),
			bearingWidth:  (int(gBnd.Min.X) >> 6),
			bearingHeight: gdescent,
			advance:       (int(gAdv) >> 6),
			tex:           NewTextures(len(fonts), "tex", rgba),
		}
		if int(gw) > maxWidth {
			maxWidth = int(gw)
		}
		if int(gh) > maxHeight {
			maxHeight = int(gh)
		}
		fonts[ch] = g
	}
	fmt.Printf("MaxWidth: %d, MaxHeight: %d\n", maxWidth, maxHeight)
	return &Charset{model.New(), fonts, maxWidth, maxHeight}
}
func NewTextures(textureId int, uniformName string, rgba *image.RGBA) texture.Textures {
	var id uint32
	glWrapper.GenTextures(1, &id)
	tex := &texture.Texture{
		TextureName: id,
		TargetId:    glwrapper.TEXTURE_2D,
		Id:          glwrapper.TEXTURE0 + uint32(textureId),
		UniformName: uniformName,
		Wrapper:     glWrapper,
		FilePath:    "tex",
	}
	tex.Bind()
	defer tex.UnBind()
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_WRAP_S, glwrapper.CLAMP_TO_EDGE)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_WRAP_R, glwrapper.CLAMP_TO_EDGE)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_MAG_FILTER, glwrapper.LINEAR)
	tex.Wrapper.TexParameteri(glwrapper.TEXTURE_2D, glwrapper.TEXTURE_MIN_FILTER, glwrapper.LINEAR)
	tex.Wrapper.TexImage2D(tex.TargetId, 0, glwrapper.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0, glwrapper.RGBA, uint32(glwrapper.UNSIGNED_BYTE), tex.Wrapper.Ptr(rgba.Pix))
	var textures texture.Textures
	textures = append(textures, tex)
	return textures

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
	Fonts := LoadCharset(FontFile, 32, 127, 40.0)
	Fonts.Print("How are You?", -0.5, 0.2, 3.0/float32(WindowWidth))
	Fonts.SetTransparent(true)
	app.AddModelToShader(Fonts, fontShader)

	for !app.GetWindow().ShouldClose() {
		glWrapper.Clear(glwrapper.COLOR_BUFFER_BIT | glwrapper.DEPTH_BUFFER_BIT)
		app.Draw()
		glfw.PollEvents()
		app.GetWindow().SwapBuffers()
	}
}
