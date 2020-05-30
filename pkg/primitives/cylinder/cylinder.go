package cylinder

import (
	"math"

	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Cylinder struct {
	Points    []mgl32.Vec3
	Normals   []mgl32.Vec3
	Indices   []uint32
	TexCoords []mgl32.Vec2
}

// New function returns a cylinder.
// Based on the following example:
// http://www.songho.ca/opengl/gl_cylinder.html
// 'rad' - the radius of the circle. 'prec' - the precision of the circle.
// 'length' - the length of the body of the cylinder.
func New(rad float32, prec int, length float32) *Cylinder {
	var points []mgl32.Vec3
	var normals []mgl32.Vec3
	var indices []uint32
	var texCoords []mgl32.Vec2

	circleVertices := circleWithRadius(rad, prec)

	// sides
	for i := 0; i < 2; i++ {
		height := -length/2 + float32(i)*length
		texCoord := float32(1.0 - i)
		k := 0
		for j := 0; j < prec; j++ {
			uX := circleVertices[k]
			uY := circleVertices[k+1]
			uZ := circleVertices[k+2]
			// position vector
			points = append(points, mgl32.Vec3{uX, uY, height})
			// normal vectors
			normals = append(normals, mgl32.Vec3{uX, uY, uZ})
			// texture coordinate
			texCoords = append(texCoords, mgl32.Vec2{float32(j) / float32(prec), texCoord})
			k = k + 3
		}
	}
	// it will be used for the generating the indices. - the staring index for the base / top circels.
	baseCenterIndex := len(points) / 3
	topCenterIndex := baseCenterIndex + prec + 1
	for i := 0; i < 2; i++ {
		height := -length/2 + float32(i)*length
		normal := float32(-1 + i*2)

		// center point
		points = append(points, mgl32.Vec3{0.0, 0.0, height})
		normals = append(normals, mgl32.Vec3{0.0, 0.0, normal})
		texCoords = append(texCoords, mgl32.Vec2{0.5, 0.5})

		k := 0
		for j := 0; j < prec; j++ {
			// position vector
			points = append(points, mgl32.Vec3{circleVertices[k], circleVertices[k+1], height})
			// normal vectors
			normals = append(normals, mgl32.Vec3{0.0, 0.0, normal})
			// texture coordinate
			texCoords = append(texCoords, mgl32.Vec2{circleVertices[k]*0.5 + 0.5, circleVertices[k+1]*0.5 + 0.5})
			k = k + 3
		}
	}
	// indices
	k1 := 0        // first vertex index of the bottom
	k2 := prec + 1 // first vertex index of the top
	// indices for the side surface
	for i := 0; i <= prec; i++ {
		indices = append(indices, uint32(k1))
		indices = append(indices, uint32(k1+1))
		indices = append(indices, uint32(k2))

		indices = append(indices, uint32(k2))
		indices = append(indices, uint32(k1+1))
		indices = append(indices, uint32(k2+1))

		k1 = k1 + 1
		k2 = k2 + 1
	}
	// indices for the base surface
	k := baseCenterIndex + 1
	for i := 0; i < prec; i++ {
		if i < prec-1 {
			indices = append(indices, uint32(baseCenterIndex))
			indices = append(indices, uint32(k+1))
			indices = append(indices, uint32(k))
		} else {
			// the last triangle
			indices = append(indices, uint32(baseCenterIndex))
			indices = append(indices, uint32(baseCenterIndex+1))
			indices = append(indices, uint32(k))
		}
		k = k + 1
	}
	// indices for the top surface
	k = topCenterIndex + 1
	for i := 0; i < prec; i++ {
		if i < prec-1 {
			indices = append(indices, uint32(topCenterIndex))
			indices = append(indices, uint32(k))
			indices = append(indices, uint32(k+1))
		} else {
			indices = append(indices, uint32(topCenterIndex))
			indices = append(indices, uint32(k))
			indices = append(indices, uint32(topCenterIndex+1))
		}
		k = k + 1
	}

	return &Cylinder{
		Points:    points,
		Normals:   normals,
		Indices:   indices,
		TexCoords: texCoords,
	}
}

// circleWithRadius returns the position vectors of
// a circle on XY plane,
func circleWithRadius(radian float32, precision int) []float32 {
	var positionVectors []float32
	sectorStep := float64(2*math.Pi) / float64(precision)
	for i := 0; i < precision; i++ {
		sectorAngle := float64(i) * sectorStep
		positionVectors = append(positionVectors, float32(math.Cos(sectorAngle))*radian)
		positionVectors = append(positionVectors, float32(math.Sin(sectorAngle))*radian)
		positionVectors = append(positionVectors, 0)
	}
	return positionVectors
}

// MaterialMeshInput method returns the verticies, indicies inputs for the NewMaterialMesh function.
func (c *Cylinder) MaterialMeshInput() (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: c.Points[i],
			Normal:   c.Normals[i],
		})
	}
	return vertices, c.Indices
}

// ColorMeshInput method returns the verticies, indicies inputs for the NewColorMesh function.
func (c *Cylinder) ColoredMeshInput(col []mgl32.Vec3) (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: c.Points[i],
			Color:    col[i%len(col)],
		})
	}
	return vertices, c.Indices
}

// TexturedMeshInput method returns the verticies, indicies inputs for the NewTexturedMesh function.
func (c *Cylinder) TexturedMeshInput() (vertex.Verticies, []uint32) {
	var vertices vertex.Verticies
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position:  c.Points[i],
			Normal:    c.Normals[i],
			TexCoords: c.TexCoords[i],
		})
	}
	return vertices, c.Indices
}
