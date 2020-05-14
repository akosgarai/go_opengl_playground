package texture

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
)

// LoadImageFromFile takes a filepath string argument.
// It loads the file, decodes it as PNG or jpg, and returns the image and error
func loadImageFromFile(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	return img, err

}

// GenTextures returns the generated uint32 name of the texture.
func genTextures() uint32 {
	var id uint32
	wrapper.GenTextures(1, &id)
	return id
}

type Texture struct {
	// The id of the texture. eg gl.TEXTURE0
	Id uint32
	// The generated name that was given by the GenTextures command
	TextureName uint32
	// The target that we use for BindTexture. (eg: TEXTURE_2D)
	TargetId uint32

	// The Uniform name of the texture
	UniformName string
}

func (t *Texture) Bind() {
	wrapper.ActiveTexture(t.Id)
	wrapper.BindTexture(t.TargetId, t.TextureName)
}
func (t *Texture) UnBind() {
	wrapper.BindTexture(t.TargetId, wrapper.TEXTURE0)
}

type Textures []*Texture

func (t *Textures) AddTexture(filePath string, wrapR, wrapS, minificationFilter, magnificationFilter int32, uniformName string) {
	img, err := loadImageFromFile(filePath)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("not 32 bit color")
	}

	tex := &Texture{
		TextureName: genTextures(),
		TargetId:    wrapper.TEXTURE_2D,
		Id:          wrapper.TEXTURE0 + uint32(len(*t)),
		UniformName: uniformName,
	}

	tex.Bind()
	defer tex.UnBind()

	wrapper.TexParameteri(wrapper.TEXTURE_2D, wrapper.TEXTURE_WRAP_R, wrapR)
	wrapper.TexParameteri(wrapper.TEXTURE_2D, wrapper.TEXTURE_WRAP_S, wrapS)
	wrapper.TexParameteri(wrapper.TEXTURE_2D, wrapper.TEXTURE_MIN_FILTER, minificationFilter)
	wrapper.TexParameteri(wrapper.TEXTURE_2D, wrapper.TEXTURE_MAG_FILTER, magnificationFilter)

	wrapper.TexImage2D(tex.TargetId, 0, wrapper.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, wrapper.RGBA, uint32(wrapper.UNSIGNED_BYTE), wrapper.Ptr(rgba.Pix))

	wrapper.GenerateMipmap(tex.TextureName)

	*t = append(*t, tex)
}

func (t Textures) UnBind() {
	for i, _ := range t {
		t[i].UnBind()
	}
}
