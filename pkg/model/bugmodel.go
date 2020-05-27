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

// position: body is a unit sphere, center is the given position. bottom half unit radius,
// its center position is {-1,0,0} from the body center.
func NewBug(position mgl32.Vec3) *Model {
	sphereBase := sphere.New(20)
	i, v := sphereBase.MaterialMeshInput()
	// Bottom supposed to be greenish color / material like emerald
	Bottom := mesh.NewMaterialMesh(i, v, material.Emerald, glWrapper)
	Bottom.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	Bottom.InitPos(mgl32.Vec3{-1, 0, 0})
	Bottom.SetPosition(position.Add(mgl32.Vec3{-1, 0, 0}))
	// Body supposed to be other green. Like green rubber
	Body := mesh.NewMaterialMesh(i, v, material.Greenrubber, glWrapper)
	Body.InitPos(mgl32.Vec3{0, 0, 0})
	Body.SetPosition(position)
	// Eyes are red. (red plastic)
	Eye1 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye1.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	Eye1.InitPos((mgl32.Vec3{1, 1, 1}).Normalize())
	Eye1.SetPosition(position.Add((mgl32.Vec3{1, 1, 1}).Normalize()))
	Eye2 := mesh.NewMaterialMesh(i, v, material.Ruby, glWrapper)
	Eye2.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	Eye2.InitPos((mgl32.Vec3{1, 1, -1}).Normalize())
	Eye2.SetPosition(position.Add((mgl32.Vec3{1, 1, -1}).Normalize()))

	m := New()
	m.AddMesh(Bottom)
	m.AddMesh(Body)
	m.AddMesh(Eye1)
	m.AddMesh(Eye2)

	return m
}
