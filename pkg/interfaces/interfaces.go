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
type DirectionalLight interface {
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
}
type PointLight interface {
	GetPosition() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
}
type SpotLight interface {
	GetPosition() mgl32.Vec3
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
	GetCutoff() float32
	GetOuterCutoff() float32
}
