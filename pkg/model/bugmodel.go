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

// position: body is a unit sphere, center is the given position. bottom half unit radius,
// its center position is {-1,0,0} from the body center.
func NewBug(position, scale mgl32.Vec3) *Bug {
	sphereBase := sphere.New(20)
	i, v := sphereBase.MaterialMeshInput()
	// Bottom supposed to be greenish color / material like emerald
	Bottom := mesh.NewMaterialMesh(i, v, material.Emerald, glWrapper)
	Bottom.SetScale(mgl32.Vec3{scale.X() * 0.5, scale.Y() * 0.5, scale.Z() * 0.5})
	Bottom.InitPos(mgl32.Vec3{scale.X() * -1, 0, 0})
	Bottom.SetPosition(position.Add(mgl32.Vec3{scale.X() * -1, 0, 0}))
	// Body supposed to be other green. Like green rubber
	Body := mesh.NewMaterialMesh(i, v, material.Greenrubber, glWrapper)
	Body.SetScale(scale)
	Body.InitPos(mgl32.Vec3{0, 0, 0})
	Body.SetPosition(position)
	// Eyes are red. (red plastic)
	Eye1 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye1.SetScale(mgl32.Vec3{scale.X() * 0.1, scale.Y() * 0.1, scale.Z() * 0.1})
	initPosBase := (mgl32.Vec3{1, 1, 1}).Normalize()
	initPosScaled := mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye1.InitPos(initPosScaled)
	Eye1.SetPosition(position.Add(initPosScaled))
	Eye2 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye2.SetScale(mgl32.Vec3{scale.X() * 0.1, scale.Y() * 0.1, scale.Z() * 0.1})
	initPosBase = (mgl32.Vec3{1, 1, -1}).Normalize()
	initPosScaled = mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye2.InitPos(initPosScaled)
	Eye2.SetPosition(position.Add(initPosScaled))

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
