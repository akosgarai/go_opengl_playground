package shader

import (
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"

	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	id uint32
}

// NewShader returns a Shader. It's inputs are the filenames of the shaders.
// It reads the files and compiles them. The shaders are attached to the shader program.
func NewShader(vertexShaderPath, fragmentShaderPath string) *Shader {
	vertexShaderSource, err := LoadShaderFromFile(vertexShaderPath)
	if err != nil {
		panic(err)
	}
	vertexShader, err := CompileShader(vertexShaderSource, wrapper.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := LoadShaderFromFile(fragmentShaderPath)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(fragmentShaderSource, wrapper.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := wrapper.CreateProgram()
	wrapper.AttachShader(program, vertexShader)
	wrapper.AttachShader(program, fragmentShader)
	wrapper.LinkProgram(program)

	return &Shader{
		id: program,
	}
}

// Use is a wrapper for gl.UseProgram
func (s *Shader) Use() {
	wrapper.UseProgram(s.id)
}

// GetId returns the program identifier of the shader.
func (s *Shader) GetId() uint32 {
	return s.id
}

// SetUniformMat4 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix4fv function
func (s *Shader) SetUniformMat4(uniformName string, mat mgl32.Mat4) {
	location := wrapper.GetUniformLocation(s.id, uniformName)
	wrapper.UniformMatrix4fv(location, 1, false, &mat[0])
}

// SetUniform3f gets an uniform name string and 3 float values as input and
// calls the gl.Uniform3f function
func (s *Shader) SetUniform3f(uniformName string, v1, v2, v3 float32) {
	location := wrapper.GetUniformLocation(s.id, uniformName)
	wrapper.Uniform3f(location, v1, v2, v3)
}

// SetUniform1f gets an uniform name string and a float value as input and
// calls the gl.Uniform1f function
func (s *Shader) SetUniform1f(uniformName string, v1 float32) {
	location := wrapper.GetUniformLocation(s.id, uniformName)
	wrapper.Uniform1f(location, v1)
}

// SetUniform1i gets an uniform name string and an integer value as input and
// calls the gl.Uniform1i function
func (s *Shader) SetUniform1i(uniformName string, v1 int32) {
	location := wrapper.GetUniformLocation(s.id, uniformName)
	wrapper.Uniform1i(location, v1)
}
