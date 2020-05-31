package model

import (
	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cylinder"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	"github.com/akosgarai/opengl_playground/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type StreetLamp struct {
	Model
}

// NewStreetLamp returns a street lamp like model. The StreetLamp is a mesh system.
// The position input describes the position of the lightsource 'Bulb'. The 'Top' is the
// child of the 'Bulb', and the 'Pole' is the children of 'Top', so that their coordinates
// are relative to their parents.
func NewMaterialStreetLamp(position mgl32.Vec3) *StreetLamp {
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
	bulb.SetPosition(mgl32.Vec3{position.X(), position.Y() + height, position.Z()})
	bulb.SetScale(bulbScale)

	// top
	top := mesh.NewMaterialMesh(topV, topI, material.Chrome, glWrapper)
	top.SetPosition(mgl32.Vec3{width / 2, width / 2, 0})
	top.SetScale(bulbScale.Mul(100))
	top.SetParent(bulb)

	// pole
	pole := mesh.NewMaterialMesh(poleV, poleI, material.Chrome, glWrapper)
	pole.SetPosition(mgl32.Vec3{(length - width) / 2, (-height - width) / 2, 0})
	pole.SetParent(top)

	m := New()
	m.AddMesh(bulb)
	m.AddMesh(top)
	m.AddMesh(pole)

	return &StreetLamp{Model: *m}
}
func NewTexturedStreetLamp(position mgl32.Vec3) *StreetLamp {
	width := float32(0.4)
	height := float32(6)
	length := float32(1.25)
	bulbScale := mgl32.Vec3{0.1, 0.1, 0.1}
	var metalTexture texture.Textures
	metalTexture.AddTexture("pkg/model/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	metalTexture.AddTexture("pkg/model/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	var bulbTexture texture.Textures
	bulbTexture.AddTexture("pkg/model/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	bulbTexture.AddTexture("pkg/model/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	bulbMaterial := material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, 256.0)
	sph := sphere.New(15)
	sphereV, sphereI := sph.TexturedMeshInput()
	poleCylinder := cylinder.New(width/2, 20, height)
	poleV, poleI := poleCylinder.TexturedMeshInput()
	topCuboid := cuboid.New(length, width, width)
	topV, topI := topCuboid.MeshInput()
	// bulb
	bulb := mesh.NewTexturedMaterialMesh(sphereV, sphereI, bulbTexture, bulbMaterial, glWrapper)
	bulb.SetPosition(mgl32.Vec3{position.X(), position.Y() + height, position.Z()})
	bulb.SetScale(bulbScale)

	// top
	top := mesh.NewTexturedMesh(topV, topI, metalTexture, glWrapper)
	top.SetPosition(mgl32.Vec3{width / 2, width / 2, 0})
	top.SetScale(bulbScale.Mul(100))
	top.SetParent(bulb)

	// pole
	pole := mesh.NewTexturedMesh(poleV, poleI, metalTexture, glWrapper)
	pole.SetPosition(mgl32.Vec3{(length - width) / 2, (-height - width) / 2, 0})
	pole.Rotate(90, mgl32.Vec3{0, 0, 1})
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
// Transformations are applied, due to the relative position.
func (s *StreetLamp) GetPolePosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(s.meshes[1].GetPosition(), s.meshes[1].ModelTransformation())
}

// GetTopPosition returns the current position of the top mesh.
// Transformations are applied, due to the relative position.
func (s *StreetLamp) GetTopPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(s.meshes[2].GetPosition(), s.meshes[2].ModelTransformation())
}
