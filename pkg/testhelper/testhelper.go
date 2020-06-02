package testhelper

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type GLWrapperMock struct{}

func (g GLWrapperMock) GenVertexArrays() uint32               { return uint32(1) }
func (g GLWrapperMock) GenBuffers() uint32                    { return uint32(1) }
func (g GLWrapperMock) BindVertexArray(vao uint32)            {}
func (g GLWrapperMock) BindBuffer(bufferType, vbo uint32)     {}
func (g GLWrapperMock) ArrayBufferData(bufferData []float32)  {}
func (g GLWrapperMock) ElementBufferData(bufferData []uint32) {}
func (g GLWrapperMock) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
}
func (g GLWrapperMock) ActiveTexture(id uint32)          {}
func (g GLWrapperMock) BindTexture(id, textureId uint32) {}
func (g GLWrapperMock) DrawTriangleElements(count int32) {}
func (g GLWrapperMock) UseProgram(id uint32)             {}
func (g GLWrapperMock) GetUniformLocation(shaderProgramId uint32, uniformName string) int32 {
	return int32(0)
}
func (g GLWrapperMock) Uniform1i(location int32, value int32) {}
func (g GLWrapperMock) CreateProgram() uint32                 { return uint32(1) }
func (g GLWrapperMock) AttachShader(program, shader uint32)   {}
func (g GLWrapperMock) LinkProgram(program uint32)            {}
func (g GLWrapperMock) UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
}
func (g GLWrapperMock) CreateShader(shaderType uint32) uint32 { return uint32(1) }
func (g GLWrapperMock) Strs(strs string) (**uint8, func()) {
	i := uint8(1)
	pi := &i
	pii := &pi
	return pii, func() {}
}
func (g GLWrapperMock) ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {}
func (g GLWrapperMock) CompileShader(id uint32)                                                 {}
func (g GLWrapperMock) GetShaderiv(shader uint32, pname uint32, params *int32)                  {}
func (g GLWrapperMock) GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
}
func (g GLWrapperMock) Str(str string) *uint8 { i := uint8(1); return &i }
func (g GLWrapperMock) InitOpenGL()           {}
func (g GLWrapperMock) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
}
func (g GLWrapperMock) Ptr(data interface{}) unsafe.Pointer {
	type tmp struct {
		d int
		p float64
	}
	var a tmp
	return unsafe.Pointer(unsafe.Offsetof(a.p))
}
func (g GLWrapperMock) GenerateMipmap(target uint32)          {}
func (g GLWrapperMock) GenTextures(n int32, textures *uint32) {}
func (g GLWrapperMock) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
}
func (g GLWrapperMock) Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {}
func (g GLWrapperMock) Uniform1f(location int32, v0 float32)                         {}
func (g GLWrapperMock) PtrOffset(offset int) unsafe.Pointer {
	var tmp struct {
		d int
		p float64
	}
	return unsafe.Pointer(unsafe.Offsetof(tmp.p))
}
func (g GLWrapperMock) DisableVertexAttribArray(index uint32)                              {}
func (g GLWrapperMock) DrawArrays(mode uint32, first int32, count int32)                   {}
func (g GLWrapperMock) TexParameteri(target uint32, pname uint32, param int32)             {}
func (g GLWrapperMock) TexParameterfv(target uint32, pname uint32, params *float32)        {}
func (g GLWrapperMock) ClearColor(red float32, green float32, blue float32, alpha float32) {}
func (g GLWrapperMock) Clear(mask uint32)                                                  {}
func (g GLWrapperMock) Enable(cap uint32)                                                  {}
func (g GLWrapperMock) DepthFunc(xfunc uint32)                                             {}
func (g GLWrapperMock) Viewport(x int32, y int32, width int32, height int32)               {}

type WindowMock struct{}

func (wm WindowMock) GetCursorPos() (float64, float64)                    { return 0.0, 0.0 }
func (wm WindowMock) SetKeyCallback(cb glfw.KeyCallback) glfw.KeyCallback { return cb }
func (wm WindowMock) SetMouseButtonCallback(cb glfw.MouseButtonCallback) glfw.MouseButtonCallback {
	return cb
}
func (wm WindowMock) ShouldClose() bool { return false }
func (wm WindowMock) SwapBuffers()      {}

type CameraMock struct{}

func (cm CameraMock) Log() string                      { return "" }
func (cm CameraMock) GetViewMatrix() mgl32.Mat4        { return mgl32.Ident4() }
func (cm CameraMock) GetProjectionMatrix() mgl32.Mat4  { return mgl32.Ident4() }
func (cm CameraMock) Walk(float32)                     {}
func (cm CameraMock) Strafe(float32)                   {}
func (cm CameraMock) Lift(float32)                     {}
func (cm CameraMock) UpdateDirection(float32, float32) {}
func (cm CameraMock) GetPosition() mgl32.Vec3          { return mgl32.Vec3{0, 0, 0} }

// This function returns true, if the given a, b is almost equal,
// the difference between them is less than epsilon.
func Float32ApproxEqual(a, b, epsilon float32) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
