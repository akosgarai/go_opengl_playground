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
func NewStreetLamp(position, scale mgl32.Vec3) *StreetLamp {
	width := float32(0.0666)
	height := float32(1)
	length := float32(0.208)
	bulbScale := mgl32.Vec3{0.016, 0.016, 0.016}
	bulbMaterial := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1, 1, 1}, 256.0)
	cube := cuboid.NewCube()
	cubeV, cubeI := cube.MeshInput()
	sph := sphere.New(15)
	sphereV, sphereI := sph.MaterialMeshInput()
	// bulb
	bulb := mesh.NewMaterialMesh(sphereV, sphereI, bulbMaterial, glWrapper)
	bulb.InitPos(mgl32.Vec3{0, scale.Y() * height, 0})
	bulb.SetPosition(mgl32.Vec3{position.X(), position.Y() + scale.Y()*height, position.Z()})
	bulb.SetScale(mgl32.Vec3{bulbScale.X() * scale.X(), bulbScale.Y() * scale.Y(), bulbScale.Z() * scale.Z()})

	// pole
	pole := mesh.NewMaterialMesh(cubeV, cubeI, material.Chrome, glWrapper)
	pole.InitPos(mgl32.Vec3{scale.X() * length / 2, scale.Y() * height / 2, 0})
	pole.SetPosition(mgl32.Vec3{(position.X() + scale.X()*(length/2)), (position.Y() + scale.Y()*(height/2)), position.Z()})
	pole.SetScale(mgl32.Vec3{scale.X() * width, scale.Y() * height, scale.Z() * width})

	// top
	top := mesh.NewMaterialMesh(cubeV, cubeI, material.Chrome, glWrapper)
	top.InitPos(mgl32.Vec3{scale.X() * width / 2, scale.Y() * (height + (width / 2)), 0})
	top.SetPosition(mgl32.Vec3{(position.X() + scale.X()*(width/2)), (position.Y() + scale.Y()*(height+(width/2))), position.Z()})
	top.SetScale(mgl32.Vec3{scale.X() * length, scale.Y() * width, scale.Z() * width})

	m := New()
	m.AddMesh(bulb)
	m.AddMesh(pole)
	m.AddMesh(top)

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
