package model

import (
	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type Room struct {
	Model
}

// NewMaterialRoom returns a Room that is based on material meshes.
// The position is the center point of the floor of the room.
// The initial orientation of the floor is the xy plane.
// The floor, the roof, the back wall, left wall, right wall are 1 * 1 * 0.05 cuboids.
// The front wall holds a door that is different colored.
func NewMaterialRoom(position mgl32.Vec3) *Room {
	floorCuboid := cuboid.New(1.0, 1.0, 0.005)
	floorV, floorI := floorCuboid.MeshInput()

	floor := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	floor.SetPosition(position)

	ceiling := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	ceiling.SetPosition(mgl32.Vec3{0.0, 1.0, 0.0})
	ceiling.SetParent(floor)

	backWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	backWall.SetPosition(mgl32.Vec3{0.0, 0.5, -0.4975})
	backWall.RotateX(90)
	backWall.SetParent(floor)

	rightWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	rightWall.SetPosition(mgl32.Vec3{-0.4975, 0.5, 0.0})
	rightWall.RotateZ(90)
	rightWall.SetParent(floor)

	leftWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	leftWall.SetPosition(mgl32.Vec3{0.4975, 0.5, 0.0})
	leftWall.RotateZ(90)
	leftWall.SetParent(floor)

	// front wall parts

	frontCuboid := cuboid.New(0.6, 1.0, 0.005)
	V, I := frontCuboid.MeshInput()
	frontWallMain := mesh.NewMaterialMesh(V, I, material.Chrome, glWrapper)
	frontWallMain.SetPosition(mgl32.Vec3{0.2, 0.5, 0.4975})
	frontWallMain.RotateX(90)
	frontWallMain.SetParent(floor)
	frontTopCuboid := cuboid.New(0.4, 0.4, 0.005)
	V, I = frontTopCuboid.MeshInput()
	frontWallRest := mesh.NewMaterialMesh(V, I, material.Chrome, glWrapper)
	frontWallRest.SetPosition(mgl32.Vec3{-0.3, 0.2, 0.4975})
	frontWallRest.RotateX(90)
	frontWallRest.SetParent(floor)
	doorCuboid := cuboid.New(0.4, 0.005, 0.6)
	V, I = doorCuboid.MeshInput()
	door := mesh.NewMaterialMesh(V, I, material.Bronze, glWrapper)
	door.SetPosition(mgl32.Vec3{-0.4975, 0.7, 0.6975})
	door.RotateY(90)
	door.SetParent(floor)

	m := New()
	m.AddMesh(floor)
	m.AddMesh(ceiling)
	m.AddMesh(backWall)
	m.AddMesh(rightWall)
	m.AddMesh(leftWall)
	m.AddMesh(frontWallMain)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
	return &Room{Model: *m}
}
func NewTextureRoom(position mgl32.Vec3) *Room {
	var concreteTexture texture.Textures
	concreteTexture.AddTexture("pkg/model/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	concreteTexture.AddTexture("pkg/model/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	var doorTexture texture.Textures
	doorTexture.AddTexture("pkg/model/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	doorTexture.AddTexture("pkg/model/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	floorCuboid := cuboid.New(1.0, 1.0, 0.005)
	floorV, floorI := floorCuboid.MeshInput()

	floor := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	floor.SetPosition(position)

	ceiling := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	ceiling.SetPosition(mgl32.Vec3{0.0, 1.0, 0.0})
	ceiling.SetParent(floor)

	backWall := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	backWall.SetPosition(mgl32.Vec3{0.0, 0.5, -0.4975})
	backWall.RotateX(90)
	backWall.SetParent(floor)

	rightWall := mesh.NewTexturedMesh(floorV, floorI, concreteTexture, glWrapper)
	rightWall.SetPosition(mgl32.Vec3{-0.4975, 0.5, 0.0})
	rightWall.RotateZ(90)
	rightWall.SetParent(floor)

	leftWall := mesh.NewTexturedMesh(floorV, floorI, concreteTexture, glWrapper)
	leftWall.SetPosition(mgl32.Vec3{0.4975, 0.5, 0.0})
	leftWall.RotateZ(90)
	leftWall.SetParent(floor)

	// front wall parts

	frontCuboid := cuboid.New(0.6, 1.0, 0.005)
	V, I := frontCuboid.MeshInput()
	frontWallMain := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain.SetPosition(mgl32.Vec3{0.2, 0.5, 0.4975})
	frontWallMain.RotateX(90)
	frontWallMain.SetParent(floor)
	frontTopCuboid := cuboid.New(0.4, 0.4, 0.005)
	V, I = frontTopCuboid.MeshInput()
	frontWallRest := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallRest.SetPosition(mgl32.Vec3{-0.3, 0.2, 0.4975})
	frontWallRest.RotateX(90)
	frontWallRest.SetParent(floor)
	doorCuboid := cuboid.New(0.4, 0.005, 0.6)
	V, I = doorCuboid.MeshInput()
	door := mesh.NewTexturedMesh(V, I, doorTexture, glWrapper)
	door.SetPosition(mgl32.Vec3{-0.4975, 0.7, 0.6975})
	door.RotateY(-90)
	door.SetParent(floor)

	m := New()
	m.AddMesh(floor)
	m.AddMesh(ceiling)
	m.AddMesh(backWall)
	m.AddMesh(rightWall)
	m.AddMesh(leftWall)
	m.AddMesh(frontWallMain)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
	return &Room{Model: *m}
}
