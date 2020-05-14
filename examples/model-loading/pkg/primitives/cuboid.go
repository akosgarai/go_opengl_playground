package primitives

import (
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Cuboid struct {
	Points   [24]mgl32.Vec3
	Normals  [6]mgl32.Vec3
	Indicies []uint32
}

// NewCube returns a unit cube (side = 1).
// The center point is the origo.
// point[1-4] the bottom side, point[5-8] top.
func NewCube() *Cuboid {
	normals := [6]mgl32.Vec3{
		mgl32.Vec3{0, -1, 0}, // bottom
		mgl32.Vec3{0, 1, 0},  // top
		mgl32.Vec3{0, 0, -1}, // front
		mgl32.Vec3{0, 0, 1},  // back
		mgl32.Vec3{-1, 0, 0}, // left
		mgl32.Vec3{1, 0, 0},  // right
	}
	// bottom
	a := mgl32.Vec3{-0.5, -0.5, -0.5}
	b := mgl32.Vec3{0.5, -0.5, -0.5}
	c := mgl32.Vec3{0.5, -0.5, 0.5}
	d := mgl32.Vec3{-0.5, -0.5, 0.5}
	// top
	e := mgl32.Vec3{-0.5, 0.5, -0.5}
	f := mgl32.Vec3{0.5, 0.5, -0.5}
	g := mgl32.Vec3{0.5, 0.5, 0.5}
	h := mgl32.Vec3{-0.5, 0.5, 0.5}
	points := [24]mgl32.Vec3{
		// bottom
		a, b, c, d,
		// top
		h, g, f, e,
		// front
		e, f, b, a,
		// back
		d, c, g, h,
		// left
		e, a, d, h,
		// right
		b, f, g, c,
	}
	indicies := []uint32{
		0, 1, 2, 0, 2, 3, // bottom
		4, 5, 6, 4, 6, 7, // top
		8, 9, 10, 8, 10, 11, // front
		12, 13, 14, 12, 14, 15, // back
		16, 17, 18, 16, 18, 19, // left
		20, 21, 22, 20, 22, 23, // right
	}
	return &Cuboid{
		Points:   points,
		Normals:  normals,
		Indicies: indicies,
	}
}

// If i want to use (yes) the element buffer, i have to return 24 verticies instead of 8.
// Because the normal vectors are not shared, i have to define the 6 side with 4 points.
// Then i can define the normals properly.
func (c *Cuboid) MeshInput() (vertex.Verticies, []uint32) {
	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
	var verticies vertex.Verticies
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			pointIndex := i*4 + j
			verticies = append(verticies, vertex.Vertex{
				Position:  c.Points[pointIndex],
				Normal:    c.Normals[i],
				TexCoords: textureCoords[j],
			})
		}
	}
	return verticies, c.Indicies
}
