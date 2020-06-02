package cuboid

import (
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	textureCoords = [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
)

type Cuboid struct {
	Points  [24]mgl32.Vec3
	Normals [6]mgl32.Vec3
	Indices []uint32
}

func getEmptyCuboid() *Cuboid {
	return &Cuboid{
		Points:  [24]mgl32.Vec3{},
		Normals: [6]mgl32.Vec3{},
		Indices: []uint32{},
	}
}

func (c *Cuboid) calculatePoints(sideWidth, sideLength, sideHeight float32) {
	// bottom
	aa := mgl32.Vec3{-sideWidth / 2, -sideHeight / 2, -sideLength / 2}
	bb := mgl32.Vec3{sideWidth / 2, -sideHeight / 2, -sideLength / 2}
	cc := mgl32.Vec3{sideWidth / 2, -sideHeight / 2, sideLength / 2}
	dd := mgl32.Vec3{-sideWidth / 2, -sideHeight / 2, sideLength / 2}
	// top
	ee := mgl32.Vec3{-sideWidth / 2, sideHeight / 2, -sideLength / 2}
	ff := mgl32.Vec3{sideWidth / 2, sideHeight / 2, -sideLength / 2}
	gg := mgl32.Vec3{sideWidth / 2, sideHeight / 2, sideLength / 2}
	hh := mgl32.Vec3{-sideWidth / 2, sideHeight / 2, sideLength / 2}
	c.Points = [24]mgl32.Vec3{
		// bottom
		aa, bb, cc, dd,
		// top
		hh, gg, ff, ee,
		// front
		ee, ff, bb, aa,
		// back
		dd, cc, gg, hh,
		// left
		ee, aa, dd, hh,
		// right
		bb, ff, gg, cc,
	}
}
func (c *Cuboid) calculateNormals() {
	c.Normals = [6]mgl32.Vec3{
		mgl32.Vec3{0, -1, 0}, // bottom
		mgl32.Vec3{0, 1, 0},  // top
		mgl32.Vec3{0, 0, -1}, // front
		mgl32.Vec3{0, 0, 1},  // back
		mgl32.Vec3{-1, 0, 0}, // left
		mgl32.Vec3{1, 0, 0},  // right
	}
}
func (c *Cuboid) calculateIndices() {
	c.Indices = []uint32{
		0, 1, 2, 0, 2, 3, // bottom
		4, 5, 6, 4, 6, 7, // top
		8, 9, 10, 8, 10, 11, // front
		12, 13, 14, 12, 14, 15, // back
		16, 17, 18, 16, 18, 19, // left
		20, 21, 22, 20, 22, 23, // right
	}
}

// New function returns a cuboid. The inputs are the width,height,length attributes.
// If the edges of the cuboid are parallel with the x,y,z axises, then  the width
// means the length in the 'x' axis, the length is the length in the 'z' axis,
// the height is the length in the 'y' axis.
func New(sideWidth, sideLength, sideHeight float32) *Cuboid {
	cuboid := getEmptyCuboid()
	cuboid.calculatePoints(sideWidth, sideLength, sideHeight)
	cuboid.calculateNormals()
	cuboid.calculateIndices()
	return cuboid
}

// NewCube returns a unit cube (side = 1).
// The center point is the origo.
// point[1-4] the bottom side, point[5-8] top.
func NewCube() *Cuboid {
	cuboid := getEmptyCuboid()
	cuboid.calculatePoints(1, 1, 1)
	cuboid.calculateNormals()
	cuboid.calculateIndices()
	return cuboid
}

// MeshInput method returns the vertices, indices inputs for the NewTexturedMesh function.
func (c *Cuboid) MeshInput() (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			pointIndex := i*4 + j
			vertices = append(vertices, vertex.Vertex{
				Position:  c.Points[pointIndex],
				Normal:    c.Normals[i],
				TexCoords: textureCoords[j],
			})
		}
	}
	return vertices, c.Indices
}

// ColoredMeshInput method returns the vertices, indices inputs for the New Mesh function.
func (c *Cuboid) ColoredMeshInput(col []mgl32.Vec3) (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			pointIndex := i*4 + j
			vertices = append(vertices, vertex.Vertex{
				Position: c.Points[pointIndex],
				Color:    col[i%len(col)],
			})
		}
	}
	return vertices, c.Indices
}

// TexturedColoredMeshInput method returns the vertices, indices inputs for the NewTexturedColoredMesh function.
func (c *Cuboid) TexturedColoredMeshInput(col []mgl32.Vec3) (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			pointIndex := i*4 + j
			vertices = append(vertices, vertex.Vertex{
				Position:  c.Points[pointIndex],
				Color:     col[i%len(col)],
				TexCoords: textureCoords[j],
			})
		}
	}
	return vertices, c.Indices
}
