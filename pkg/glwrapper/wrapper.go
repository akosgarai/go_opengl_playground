package glwrapper

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	ARRAY_BUFFER         = gl.ARRAY_BUFFER
	ELEMENT_ARRAY_BUFFER = gl.ELEMENT_ARRAY_BUFFER
	TEXTURE_2D           = gl.TEXTURE_2D
	VERTEX_SHADER        = gl.VERTEX_SHADER
	FRAGMENT_SHADER      = gl.FRAGMENT_SHADER
	COMPILE_STATUS       = gl.COMPILE_STATUS
	INFO_LOG_LENGTH      = gl.INFO_LOG_LENGTH
	FALSE                = gl.FALSE
	TEXTURE0             = gl.TEXTURE0
	TEXTURE1             = gl.TEXTURE1
	TEXTURE2             = gl.TEXTURE2
	TEXTURE_WRAP_R       = gl.TEXTURE_WRAP_R
	TEXTURE_WRAP_S       = gl.TEXTURE_WRAP_S
	TEXTURE_MIN_FILTER   = gl.TEXTURE_MIN_FILTER
	TEXTURE_MAG_FILTER   = gl.TEXTURE_MAG_FILTER
	RGBA                 = gl.RGBA
	UNSIGNED_BYTE        = gl.UNSIGNED_BYTE
	FLOAT                = gl.FLOAT
	POINTS               = gl.POINTS
	TRIANGLES            = gl.TRIANGLES
	TEXTURE_BORDER_COLOR = gl.TEXTURE_BORDER_COLOR
	CLAMP_TO_EDGE        = gl.CLAMP_TO_EDGE
	LINEAR               = gl.LINEAR
	COLOR_BUFFER_BIT     = gl.COLOR_BUFFER_BIT
	DEPTH_BUFFER_BIT     = gl.DEPTH_BUFFER_BIT
	DEPTH_TEST           = gl.DEPTH_TEST
	LESS                 = gl.LESS
)

// Wrapper for gl.GenVertexArrays function.
func GenVertexArrays() uint32 {
	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	return vertexArrayObject
}

// Wrapper for gl.GenBuffers function.
func GenBuffers() uint32 {
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	return vertexBufferObject
}

// Wrapper for gl.BindVertexArray function.
func BindVertexArray(vao uint32) {
	gl.BindVertexArray(vao)
}

// Wrapper for gl.BindBuffer function.
func BindBuffer(bufferType, vbo uint32) {
	gl.BindBuffer(bufferType, vbo)
}

// Wrapper for gl.BufferData function but for ARRAY_BUFFER.
func ArrayBufferData(bufferData []float32) {
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(bufferData), gl.Ptr(bufferData), gl.STATIC_DRAW)
}

// Wrapper for gl.BufferData function, but for ELEMENT_ARRAY_BUFFER.
func ElementBufferData(bufferData []uint32) {
	// a 32-bit uint has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(bufferData), gl.Ptr(bufferData), gl.STATIC_DRAW)
}

// VertexAttribPointer enables and sets the pointer.
func VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	gl.EnableVertexAttribArray(index)
	gl.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
}

// Wrapper for gl.ActiveTexture function.
func ActiveTexture(id uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + id)
}

// Wrapper for gl.BindTexture function.
func BindTexture(id, textureId uint32) {
	gl.BindTexture(id, textureId)
}

// Wrapper for gl.DrawElements function in triangle mode.
func DrawTriangleElements(count int32) {
	gl.DrawElements(gl.TRIANGLES, count, gl.UNSIGNED_INT, gl.PtrOffset(0))
}

// Wrapper for gl.UseProgram function.
func UseProgram(id uint32) {
	gl.UseProgram(id)
}

// Use is a wrapper for gl.GetUniformLocation
func GetUniformLocation(shaderProgramId uint32, uniformName string) int32 {
	return gl.GetUniformLocation(shaderProgramId, gl.Str(uniformName+"\x00"))
}

// Uniform1i gets an uniform name string and 3 float values as input and
// calls the gl.Uniform1i function
func Uniform1i(location int32, value int32) {
	gl.Uniform1i(location, value)
}

// Wrapper for gl.Use function.
func CreateProgram() uint32 {
	program := gl.CreateProgram()
	return program
}

// Wrapper for gl.AttachShader function.
func AttachShader(program, shader uint32) {
	gl.AttachShader(program, shader)
}

// Wrapper for gl.LinkProgram function.
func LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

// Wrapper for gl.UniformMatrix4fv function.
func UniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
	gl.UniformMatrix4fv(location, count, transpose, value)
}

// Wrapper for gl.CreateShader function.
func CreateShader(shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	return shader
}

// Wrapper for gl.Strs function.
func Strs(strs string) (**uint8, func()) {
	csources, free := gl.Strs(strs)
	return csources, free
}

// Wrapper for gl.ShaderSource function.
func ShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	gl.ShaderSource(shader, count, xstring, length)
}

// Wrapper for gl.CompileShader function.
func CompileShader(id uint32) {
	gl.CompileShader(id)
}

// Wrapper for gl.GetShaderiv function.
func GetShaderiv(shader uint32, pname uint32, params *int32) {
	gl.GetShaderiv(shader, pname, params)
}

// Wrapper for gl.GetShaderInfoLog function.
func GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	gl.GetShaderInfoLog(shader, bufSize, length, infoLog)
}

// Wrapper for gl.Str function.
func Str(str string) *uint8 {
	return gl.Str(str)
}

// InitOpenGL is for initializing the gl lib. It also prints out the gl version.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

// Wrapper for gl.TexImage2D function.
func TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, pixels)
}

// Wrapper for gl.Ptr function.
func Ptr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}

// Wrapper for gl.GenerateMipmap function.
func GenerateMipmap(target uint32) {
	gl.GenerateMipmap(target)
}

// Wrapper for gl.GenTextures function.
func GenTextures(n int32, textures *uint32) {
	gl.GenTextures(n, textures)
}

// Wrapper for gl.UniformMatrix3fv function.
func UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
	gl.UniformMatrix3fv(location, count, transpose, value)
}

// Wrapper for gl.Uniform3f function
func Uniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(location, v0, v1, v2)
}

// Wrapper for gl.Uniform1f function.
func Uniform1f(location int32, v0 float32) {
	gl.Uniform1f(location, v0)
}

// Wrapper for gl.PtrOffset function.
func PtrOffset(offset int) unsafe.Pointer {
	return gl.PtrOffset(offset)
}

// Wrapper for gl.DisableVertexAttribArray function.
func DisableVertexAttribArray(index uint32) {
	gl.DisableVertexAttribArray(index)
}

// Wrapper for gl.DrawArrays function.
func DrawArrays(mode uint32, first int32, count int32) {
	gl.DrawArrays(mode, first, count)
}

// Wrapper for gl.TexParameteri function.
func TexParameteri(target uint32, pname uint32, param int32) {
	gl.TexParameteri(target, pname, param)
}

// Wrapper fro gl.TexParameterfv function.
func TexParameterfv(target uint32, pname uint32, params *float32) {
	gl.TexParameterfv(target, pname, params)
}

// Wrapper for gl.ClearColor function.
func ClearColor(red float32, green float32, blue float32, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

// Wrapper fro gl.Clear function.
func Clear(mask uint32) {
	gl.Clear(mask)
}

// Wrapper for gl.Enable function.
func Enable(cap uint32) {
	gl.Enable(cap)
}

// Wrapper for gl.DepthFunc function.
func DepthFunc(xfunc uint32) {
	gl.DepthFunc(xfunc)
}
