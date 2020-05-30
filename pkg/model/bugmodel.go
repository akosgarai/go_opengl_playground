package model

import (
	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	glWrapper glwrapper.Wrapper
)

type Bug struct {
	Model
}

// NewBug returns a bug instance. A Bug is a sphered mesh system. Its 'Body'
// is the parent mesh. Its position is absolute in the world coordinate system.
// The other 3 mesh, the 'Bottom', the 'Eye1' and 'Eye2' are child meshes, their
// position is relative to the parent.
func NewBug(position, scale mgl32.Vec3) *Bug {
	sphereBase := sphere.New(20)
	i, v := sphereBase.MaterialMeshInput()
	// Body supposed to be other green. Like green rubber
	Body := mesh.NewMaterialMesh(i, v, material.Greenrubber, glWrapper)
	Body.SetScale(scale)
	Body.SetPosition(position)

	// Bottom supposed to be greenish color / material like emerald
	Bottom := mesh.NewMaterialMesh(i, v, material.Emerald, glWrapper)
	Bottom.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	Bottom.SetPosition(mgl32.Vec3{scale.X() * -1, 0, 0})
	Bottom.SetParent(Body)
	// Eyes are red. (red plastic)
	Eye1 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye1.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	initPosBase := (mgl32.Vec3{1, 1, 1}).Normalize()
	initPosScaled := mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye1.SetPosition(initPosScaled)
	Eye1.SetParent(Body)

	Eye2 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye2.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	initPosBase = (mgl32.Vec3{1, 1, -1}).Normalize()
	initPosScaled = mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye2.SetPosition(initPosScaled)
	Eye2.SetParent(Body)

	m := New()
	m.AddMesh(Bottom)
	m.AddMesh(Body)
	m.AddMesh(Eye1)
	m.AddMesh(Eye2)

	return &Bug{Model: *m}
}

// GetBottomPosition returns the current position of the bottom mesh.
func (b *Bug) GetBottomPosition() mgl32.Vec3 {
	return b.meshes[0].GetPosition()
}

// GetBodyPosition returns the current position of the body mesh.
func (b *Bug) GetBodyPosition() mgl32.Vec3 {
	return b.meshes[1].GetPosition()
}

// GetEye1Position returns the current position of the eye1 mesh.
func (b *Bug) GetEye1Position() mgl32.Vec3 {
	return b.meshes[2].GetPosition()
}

// GetEye2Position returns the current position of the eye2 mesh.
func (b *Bug) GetEye2Position() mgl32.Vec3 {
	return b.meshes[3].GetPosition()
}
