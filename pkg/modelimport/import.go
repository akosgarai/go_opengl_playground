package modelimport

import (
	"errors"
	"fmt"

	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"
	"github.com/akosgarai/opengl_playground/pkg/texture"

	"github.com/akosgarai/gwob"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	glWrapper glwrapper.Wrapper
)

type Import struct {
	objectFile string
	basePath   string
	meshes     []interfaces.Mesh
	object     *gwob.Obj
	material   gwob.MaterialLib
}

func New(basePath, objectFileName string) *Import {
	return &Import{
		objectFile: objectFileName,
		basePath:   basePath,
		meshes:     []interfaces.Mesh{},
	}
}
func (i *Import) GetMeshes() []interfaces.Mesh {
	return i.meshes
}

func (i *Import) loadObjectFile(options *gwob.ObjParserOptions) error {
	var errObj error
	objectFile := i.basePath + "/" + i.objectFile
	i.object, errObj = gwob.NewObjFromFile(objectFile, options)
	return errObj
}
func (i *Import) loadMaterialFile(options *gwob.ObjParserOptions) error {
	// load material lib
	var errMtl error
	objectFile := i.basePath + "/" + i.object.Mtllib
	i.material, errMtl = gwob.ReadMaterialLibFromFile(objectFile, options)
	return errMtl
}
func (i *Import) getVerticesAndIndices(g *gwob.Group) (vertex.Verticies, []uint32) {
	// Generate meshes. loop indexes from indexbegin, to indexbegin+indexcount.
	var vertices vertex.Verticies
	var indices []uint32
	// It maps the gwob indexis to the current mesh indexes.
	indexMap := make(map[int]uint32)
	for index := g.IndexBegin; index < g.IndexBegin+g.IndexCount; index++ {
		indexValue := i.object.Indices[index]
		if _, ok := indexMap[indexValue]; !ok {
			mappedValue := uint32(len(indexMap))
			indexMap[indexValue] = mappedValue
			var vert vertex.Vertex
			positionFirstIndex := indexValue * (i.object.StrideSize / 4)
			vert.Position = mgl32.Vec3{i.object.Coord[positionFirstIndex], i.object.Coord[positionFirstIndex+1], i.object.Coord[positionFirstIndex+2]}
			if i.object.TextCoordFound {
				texIndex := positionFirstIndex + i.object.StrideOffsetTexture/4
				vert.TexCoords = mgl32.Vec2{i.object.Coord[texIndex], i.object.Coord[texIndex+1]}
			}
			if i.object.NormCoordFound {
				normIndex := positionFirstIndex + i.object.StrideOffsetNormal/4
				vert.Normal = mgl32.Vec3{i.object.Coord[normIndex], i.object.Coord[normIndex+1], i.object.Coord[normIndex+2]}
			}

			vertices = append(vertices, vert)
		}
		indices = append(indices, indexMap[indexValue])
	}
	return vertices, indices
}
func (i *Import) getMaterial(mtl *gwob.Material) *material.Material {
	return material.New(
		mgl32.Vec3{mtl.Ka[0], mtl.Ka[1], mtl.Ka[2]},
		mgl32.Vec3{mtl.Kd[0], mtl.Kd[1], mtl.Kd[2]},
		mgl32.Vec3{mtl.Ks[0], mtl.Ks[1], mtl.Ks[2]},
		float32(mtl.Ns),
	)
}
func (i *Import) getTextures(mtl *gwob.Material) texture.Textures {
	var tex texture.Textures
	if mtl.MapKa != "" {
		tex.AddTexture(mtl.MapKa, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.ambient", glWrapper)
	}
	if mtl.MapKd != "" {
		tex.AddTexture(mtl.MapKd, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", glWrapper)
	}
	if mtl.MapKs != "" {
		tex.AddTexture(mtl.MapKs, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.scalar", glWrapper)
	}
	return tex
}
func (i *Import) getColor(mtl *gwob.Material) []mgl32.Vec3 {
	return []mgl32.Vec3{
		mgl32.Vec3{mtl.Ka[0], mtl.Ka[1], mtl.Ka[2]},
	}
}
func (i *Import) makeMeshes() []error {
	var result []error
	for _, g := range i.object.Groups {
		vertices, indices := i.getVerticesAndIndices(g)
		mtl, found := i.material.Lib[g.Usemtl]
		if found {
			var mat *material.Material
			var tex texture.Textures
			if i.object.NormCoordFound {
				mat = i.getMaterial(mtl)
			}
			if i.object.TextCoordFound {
				tex = i.getTextures(mtl)
			}
			if i.object.NormCoordFound {
				if i.object.TextCoordFound {
					texturedMaterialMesh := mesh.NewTexturedMaterialMesh(vertices, indices, tex, mat, glWrapper)
					i.meshes = append(i.meshes, texturedMaterialMesh)
				} else {
					materialMesh := mesh.NewMaterialMesh(vertices, indices, mat, glWrapper)
					i.meshes = append(i.meshes, materialMesh)
				}
			} else {
				if i.object.TextCoordFound {
					color := i.getColor(mtl)
					texturedColoredMesh := mesh.NewTexturedColoredMesh(vertices, indices, tex, color, glWrapper)
					i.meshes = append(i.meshes, texturedColoredMesh)
				} else {
					result = append(result, errors.New("Could not transform to mesh."))
				}
			}
		} else {
			// For point meshes we don't need material.
			pointMesh := mesh.NewPointMesh(glWrapper)
			for _, vert := range vertices {
				pointMesh.AddVertex(vert)
			}
			i.meshes = append(i.meshes, pointMesh)
		}
	}
	return result
}
func (i *Import) Import() {
	options := &gwob.ObjParserOptions{}

	errObj := i.loadObjectFile(options)
	if errObj != nil {
		fmt.Printf("Error during object file parse. '%s'", errObj.Error())
		panic(errObj)
	}

	// load material lib
	errMtl := i.loadMaterialFile(options)
	if errMtl != nil {
		fmt.Printf("Error during material file parse. '%s'", errMtl.Error())
		panic(errMtl)
	}

	errProcess := i.makeMeshes()
	if len(errProcess) != 0 {
		fmt.Println("Error during mesh construction.")
		for _, err := range errProcess {
			fmt.Printf(" - %s\n", err.Error())
		}
		return
	}
}
