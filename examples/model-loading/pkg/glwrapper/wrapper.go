package glwrapper

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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
)

// Wrapper for gl.GenVertexArray function.
func GenVertexArray() uint32 {
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
func VertexAttribPointer(index uint32, size, stride int32, offset int) {
	gl.EnableVertexAttribArray(index)
	gl.VertexAttribPointer(index, size, gl.FLOAT, false, stride, gl.PtrOffset(offset))
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

// Use is a wrapper for gl.UseProgram
func Use(id uint32) {
	gl.UseProgram(id)
}

// Use is a wrapper for gl.GetUniformLocation
func GetUniformLocation(uniformName string, shaderProgramId uint32) int32 {
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
func UniformMatrix4fv(location int32, mat mgl32.Mat4) {
	gl.UniformMatrix4fv(location, 1, false, &mat[0])
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
