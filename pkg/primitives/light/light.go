package light

import (
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
)

type Light struct {
	position mgl32.Vec3

	ambient  mgl32.Vec3
	diffuse  mgl32.Vec3
	specular mgl32.Vec3
}

func New(position, ambientComponent, diffuseComponent, specularComponent mgl32.Vec3) *Light {
	return &Light{
		position: position,
		ambient:  ambientComponent,
		diffuse:  diffuseComponent,
		specular: specularComponent,
	}
}

// Log returns the current state of the object
func (l *Light) Log() string {
	logString := "Light\n"
	logString += " - Position: Vector{" + trans.Vec3ToString(l.position) + "}\n"
	logString += " - Ambient: Vector{" + trans.Vec3ToString(l.ambient) + "}\n"
	logString += " - Diffuse: Vector{" + trans.Vec3ToString(l.diffuse) + "}\n"
	logString += " - Specualar: Vector{" + trans.Vec3ToString(l.specular) + "}\n"
	return logString
}

// GetDiffuse returns the diffuse color of the material
func (l *Light) GetAmbient() mgl32.Vec3 {
	return l.ambient
}

// GetDiffuse returns the diffuse color of the material
func (l *Light) GetDiffuse() mgl32.Vec3 {
	return l.diffuse
}

// GetSpecular returns the specular color of the material
func (l *Light) GetSpecular() mgl32.Vec3 {
	return l.specular
}

// GetPosition returns the shininess of the material
func (l *Light) GetPosition() mgl32.Vec3 {
	return l.position
}

// SetPosition updates the position of the light
func (l *Light) SetPosition(pos mgl32.Vec3) {
	l.position = pos
}
