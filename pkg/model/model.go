package model

import (
	"fmt"

	"github.com/akosgarai/opengl_playground/pkg/export"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"

	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	meshes    []interfaces.Mesh
	directory string
}

func New() *Model {
	return &Model{
		meshes: []interfaces.Mesh{},
	}
}

// AddMesh function adds a mesh to the meshes.
func (m *Model) AddMesh(msh interfaces.Mesh) {
	m.meshes = append(m.meshes, msh)
}

// Update function loops over each of the meshes and calls their Update function.
func (m *Model) Update(dt float64) {
	for i, _ := range m.meshes {
		m.meshes[i].Update(dt)
	}
}

// Draw function loops over each of the meshes and calls their Draw function.
func (m *Model) Draw(s interfaces.Shader) {
	for i, _ := range m.meshes {
		m.meshes[i].Draw(s)
	}
}

// Export function exports the meshes to a file
func (m *Model) Export(path string) {
	exporter := export.New(m.meshes)
	err := exporter.Export(path)
	if err != nil {
		fmt.Printf("Export failed. '%s'\n", err.Error())
	}
}

// SetSpeed function loops over each of the meshes and calls their SetSpeed function.
func (m *Model) SetSpeed(s float32) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].SetSpeed(s)
		}
	}
}

// SetDirection function loops over each of the meshes and calls their SetDirection function.
func (m *Model) SetDirection(p mgl32.Vec3) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].SetDirection(p)
		}
	}
}

// RotateWithAngle function rotates the model with the given angle (has to be degree).
// It calls the TransformOrigin function of each mesh.
func (m *Model) RotateWithAngle(angleDeg float32, axisVector mgl32.Vec3) {
	for i, _ := range m.meshes {
		m.meshes[i].TransformOrigin(
			mgl32.HomogRotate3D(mgl32.DegToRad(angleDeg), axisVector))
	}
}

// RotateWithAngleV2 function rotates the model with the given angle (has to be degree).
// It calls the TransformOrigin function of each mesh.
func (m *Model) RotateWithAngleV2(angleDeg float32, axisVector mgl32.Vec3) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].SetRotationAngle(mgl32.DegToRad(angleDeg))
			m.meshes[i].SetRotationAxis(axisVector)
		}
	}
}
