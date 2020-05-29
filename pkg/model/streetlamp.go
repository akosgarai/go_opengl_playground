package model

import (
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"

	"github.com/go-gl/mathgl/mgl32"
)

type StreetLamp struct {
	Model
}

// NewStreetLamp returns a street lamp like model. It has a pole, a top and a bulb.
// The bulb could be used as a light source (spot light). The bulb poisition is in the
// 0, height, 0 coordinate. The top position is width/2, height + width / 2, 0. The pole
// is in the length/2, height/2, 0 position.
func NewStreetLamp(position mgl32.Vec3) *StreetLamp {
	width := float32(0.4)
	height := float32(6)
	length := float32(1.25)
	bulbScale := mgl32.Vec3{0.1, 0.1, 0.1}
	bulbMaterial := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 256.0)
	sph := sphere.New(15)
	sphereV, sphereI := sph.MaterialMeshInput()
	poleCuboid := cuboid.New(width, width, height)
	poleV, poleI := poleCuboid.MeshInput()
	topCuboid := cuboid.New(length, width, width)
	topV, topI := topCuboid.MeshInput()
	// bulb
	bulb := mesh.NewMaterialMesh(sphereV, sphereI, bulbMaterial, glWrapper)
	bulb.InitPos(mgl32.Vec3{0, 0, 0})
	bulb.SetPosition(mgl32.Vec3{position.X(), position.Y() + height, position.Z()})
	bulb.SetScale(bulbScale)

	// top
	top := mesh.NewMaterialMesh(topV, topI, material.Chrome, glWrapper)
	top.InitPos(mgl32.Vec3{width / 2, width / 2, 0})
	top.SetScale(bulbScale.Mul(100))
	top.SetParent(bulb)

	// pole
	pole := mesh.NewMaterialMesh(poleV, poleI, material.Chrome, glWrapper)
	pole.InitPos(mgl32.Vec3{(length - width) / 2, (-height - width) / 2, 0})
	//pole.SetScale(bulbScale.Mul(100))
	pole.SetParent(top)

	m := New()
	m.AddMesh(bulb)
	m.AddMesh(top)
	m.AddMesh(pole)

	return &StreetLamp{Model: *m}
}

// GetBulbPosition returns the current position of the bulb mesh.
func (s *StreetLamp) GetBulbPosition() mgl32.Vec3 {
	return s.meshes[0].GetPosition()
}

// GetPolePosition returns the current position of the pole mesh.
func (s *StreetLamp) GetPolePosition() mgl32.Vec3 {
	return s.meshes[1].GetPosition()
}

// GetTopPosition returns the current position of the top mesh.
func (s *StreetLamp) GetTopPosition() mgl32.Vec3 {
	return s.meshes[2].GetPosition()
}
