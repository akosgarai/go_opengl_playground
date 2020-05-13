package interfaces

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Shader interface {
	Use()
	SetUniformMat4(string, mgl32.Mat4)
	GetId() uint32
	SetUniform3f(string, float32, float32, float32)
	SetUniform1f(string, float32)
	SetUniform1i(string, int32)
}
