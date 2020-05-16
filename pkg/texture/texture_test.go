package texture

import (
	"testing"
	"unsafe"

	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
)

type glWrapperMock struct{}

func (g glWrapperMock) GenVertexArrays() uint32               { return uint32(1) }
func (g glWrapperMock) GenBuffers() uint32                    { return uint32(1) }
func (g glWrapperMock) BindVertexArray(vao uint32)            {}
func (g glWrapperMock) BindBuffer(bufferType, vbo uint32)     {}
func (g glWrapperMock) ArrayBufferData(bufferData []float32)  {}
func (g glWrapperMock) ElementBufferData(bufferData []uint32) {}
func (g glWrapperMock) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
}
func (g glWrapperMock) ActiveTexture(id uint32)          {}
func (g glWrapperMock) BindTexture(id, textureId uint32) {}
func (g glWrapperMock) DrawTriangleElements(count int32) {}
func (g glWrapperMock) UseProgram(id uint32)             {}
func (g glWrapperMock) GetUniformLocation(shaderProgramId uint32, uniformName string) int32 {
	return int32(0)
}
func (g glWrapperMock) Uniform1i(location int32, value int32) {}
func (g glWrapperMock) CreateProgram() uint32                 { return uint32(1) }
func (g glWrapperMock) AttachShader(program, shader uint32)   {}
func (g glWrapperMock) LinkProgram(program uint32)            {}
func (g glWrapperMock) UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
}
func (g glWrapperMock) CreateShader(shaderType uint32) uint32 { return uint32(1) }
func (g glWrapperMock) Strs(strs string) (**uint8, func()) {
	i := uint8(1)
	pi := &i
	pii := &pi
	return pii, func() {}
}
func (g glWrapperMock) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {}
func (g glWrapperMock) CompileShader(id uint32)                                                 {}
func (g glWrapperMock) GetShaderiv(shader uint32, pname uint32, params *int32)                  {}
func (g glWrapperMock) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
}
func (g glWrapperMock) Str(str string) *uint8 { i := uint8(1); return &i }
func (g glWrapperMock) InitOpenGL()           {}
func (g glWrapperMock) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (g glWrapperMock) Ptr(data interface{}) unsafe.Pointer {
	type tmp struct {
		d int
		p float64
	}
	var a tmp
	return unsafe.Pointer(unsafe.Offsetof(a.p))
}
func (g glWrapperMock) GenerateMipmap(target uint32)          {}
func (g glWrapperMock) GenTextures(n int32, textures *uint32) {}
func (g glWrapperMock) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
}
func (g glWrapperMock) Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {}
func (g glWrapperMock) Uniform1f(location int32, v0 float32)                         {}
func (g glWrapperMock) PtrOffset(offset int) unsafe.Pointer {
	var tmp struct {
		d int
		p float64
	}
	return unsafe.Pointer(unsafe.Offsetof(tmp.p))
}
func (g glWrapperMock) DisableVertexAttribArray(index uint32)                              {}
func (g glWrapperMock) DrawArrays(mode uint32, first int32, count int32)                   {}
func (g glWrapperMock) TexParameteri(target uint32, pname uint32, param int32)             {}
func (g glWrapperMock) TexParameterfv(target uint32, pname uint32, params *float32)        {}
func (g glWrapperMock) ClearColor(red float32, green float32, blue float32, alpha float32) {}
func (g glWrapperMock) Clear(mask uint32)                                                  {}
func (g glWrapperMock) Enable(cap uint32)                                                  {}
func (g glWrapperMock) DepthFunc(xfunc uint32)                                             {}
func (g glWrapperMock) Viewport(x int32, y int32, width int32, height int32)               {}

var testGlWrapper glWrapperMock

func TestLoadImageFromFile(t *testing.T) {
	_, err := loadImageFromFile("this-image-does-not-exist.jpg")
	if err == nil {
		t.Error("Image load should be failed.")
	}
	_, err = loadImageFromFile("testing.jpg")
	if err != nil {
		t.Error("Issue during load.")
	}
}
func TestAddTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("testing.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", testGlWrapper)
	if len(textures) != 1 {
		t.Error("AddTexture should be successful")
	}
}
func TestAddTextureInvalidName(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("AddTexture should have panicked due to the missing file!")
			}
		}()
		var textures Textures
		textures.AddTexture("this-image-does-not-exist.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", testGlWrapper)
	}()
}
func TestGenTextures(t *testing.T) {
	genTextures(testGlWrapper)
}
func TestBindTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("testing.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures[0].Bind()
}
func TestUnBindTexture(t *testing.T) {
	var textures Textures
	textures.AddTexture("testing.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures[0].UnBind()
}
func TestUnBindTextures(t *testing.T) {
	var textures Textures
	textures.AddTexture("testing.jpg", wrapper.CLAMP_TO_EDGE, wrapper.CLAMP_TO_EDGE, wrapper.LINEAR, wrapper.LINEAR, "material.diffuse", testGlWrapper)
	textures.UnBind()
}
