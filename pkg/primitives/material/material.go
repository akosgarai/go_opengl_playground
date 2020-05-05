package material

import (
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
)

type Material struct {
	ambient   mgl32.Vec3
	diffuse   mgl32.Vec3
	specular  mgl32.Vec3
	shininess float32
}

func New(ambient, diffuse, specular mgl32.Vec3, shininess float32) *Material {
	return &Material{
		ambient:   ambient,
		diffuse:   diffuse,
		specular:  specular,
		shininess: shininess,
	}
}

func (m *Material) Log() string {
	logString := "Material\n"
	logString += " - Ambient: Vector{" + trans.Vec3ToString(m.ambient) + "}\n"
	logString += " - Diffuse: Vector{" + trans.Vec3ToString(m.diffuse) + "}\n"
	logString += " - Specualar: Vector{" + trans.Vec3ToString(m.specular) + "}\n"
	logString += " - Shininess: " + trans.Float32ToString(m.shininess) + "\n"
	return logString
}

// GetDiffuse returns the diffuse color of the material
func (m *Material) GetAmbient() mgl32.Vec3 {
	return m.ambient
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

var (
	TestMaterialGreen = &Material{
		diffuse:   mgl32.Vec3{0, 1, 0},
		specular:  mgl32.Vec3{0, 1, 0},
		shininess: 0,
	}
	TestMaterialRed = &Material{
		diffuse:   mgl32.Vec3{1, 0, 0},
		specular:  mgl32.Vec3{1, 0, 0},
		shininess: 0,
	}
)
