package shader

import (
	"github.com/go-gl/mathgl/mgl32"
)

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
}

type DirectionalLightSource struct {
	LightSource          DirectionalLight
	DirectionUniformName string
	AmbientUniformName   string
	DiffuseUniformName   string
	SpecularUniformName  string
}
type PointLightSource struct {
	LightSource              PointLight
	PositionUniformName      string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
}
type SpotLightSource struct {
	LightSource              SpotLight
	PositionUniformName      string
	DirectionUniformName     string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
	CutoffUniformName        string
}
