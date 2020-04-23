package primitives

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Material struct {
	diffuse   mgl32.Vec3
	specular  mgl32.Vec3
	shininess float32
}

// GetDiffuse returns the diffuse color of the material
func (m *Material) GetDiffuse() mgl32.Vec3 {
	return m.diffuse
}

// GetSpecular returns the specular color of the material
func (m *Material) GetSpecular() mgl32.Vec3 {
	return m.specular
}

// GetShininess returns the shininess of the material
func (m *Material) GetShininess() float32 {
	return m.shininess
}
