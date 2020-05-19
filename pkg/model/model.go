package model

import (
	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"

	"github.com/akosgarai/gwob"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	glWrapper glwrapper.Wrapper
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

// NewModelFromFile function gets a path as input. It loads the stuff from the paths to a model
// that the function returns.
func NewModelFromFile(path string) *Model {
	m := &Model{}
	m.loadModel(path)
	return m
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
func (m *Model) loadModel(path string) {
	options := &gwob.ObjParserOptions{}

	o, errObj := gwob.NewObjFromFile(path, options)

	if errObj != nil {
		panic(errObj)
	}
	// load material lib
	mtlLib, errMtl := gwob.ReadMaterialLibFromFile(o.Mtllib, options)
	if errMtl != nil {
		panic(errMtl)
	}

	for _, g := range o.Groups {
		// Generate meshes. loop indexes from indexbegin, to indexbegin+indexcount.
		var vertices vertex.Verticies
		var indices []uint32
		// It maps the gwob indexis to the current mesh indexes.
		indexMap := make(map[int]uint32)
		for i := g.IndexBegin; i < g.IndexBegin+g.IndexCount; i++ {
			indexValue := o.Indices[i]
			if _, ok := indexMap[indexValue]; !ok {
				mappedValue := uint32(len(indexMap))
				indexMap[indexValue] = mappedValue
				var vert vertex.Vertex
				positionFirstIndex := indexValue * (o.StrideSize / 4)
				vert.Position = mgl32.Vec3{o.Coord[positionFirstIndex], o.Coord[positionFirstIndex+1], o.Coord[positionFirstIndex+2]}
				if o.TextCoordFound {
					texIndex := positionFirstIndex + o.StrideOffsetTexture/4
					vert.TexCoords = mgl32.Vec2{o.Coord[texIndex], o.Coord[texIndex+1]}
				}
				if o.NormCoordFound {
					normIndex := positionFirstIndex + o.StrideOffsetNormal/4
					vert.Normal = mgl32.Vec3{o.Coord[normIndex], o.Coord[normIndex+1], o.Coord[normIndex+2]}
				}

				vertices = append(vertices, vert)
			}
			indices = append(indices, indexMap[indexValue])
		}
		mtl, found := mtlLib.Lib[g.Usemtl]
		if found {
			mat := material.New(
				mgl32.Vec3{mtl.Ka[0], mtl.Ka[1], mtl.Ka[2]},
				mgl32.Vec3{mtl.Kd[0], mtl.Kd[1], mtl.Kd[2]},
				mgl32.Vec3{mtl.Ks[0], mtl.Ks[1], mtl.Ks[2]},
				float32(64),
			)
			materialMesh := mesh.NewMaterialMesh(vertices, indices, mat, glWrapper)
			m.meshes = append(m.meshes, materialMesh)
		}
	}
}
