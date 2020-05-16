package primitives

import (
	"math"

	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Sphere struct {
	Points    []mgl32.Vec3
	Indicies  []uint32
	TexCoords []mgl32.Vec2
}

// based on this: http://www.songho.ca/opengl/gl_sphere.html
func NewSphere(precision int) *Sphere {
	var points []mgl32.Vec3
	var indicies []uint32
	var texCoords []mgl32.Vec2
	sectorStep := (2 * math.Pi) / float64(precision)
	stackStep := (math.Pi) / float64(precision)

	for i := 0; i <= precision; i++ {
		stackAngle := (math.Pi / 2) - (float64(i) * stackStep)
		xy := math.Cos(stackAngle)
		z := math.Sin(stackAngle)

		k1 := i * (precision + 1)
		k2 := k1 + precision + 1

		for j := 0; j <= precision; j++ {
			sectorAngle := float64(j) * sectorStep
			x := xy * math.Cos(sectorAngle)
			y := xy * math.Sin(sectorAngle)

			points = append(points, mgl32.Vec3{float32(x), float32(y), float32(z)})
			// textures [0-1], s,t
			s := float32(j) / float32(precision)
			t := float32(i) / float32(precision)
			texCoords = append(texCoords, mgl32.Vec2{s, t})
			// indicies
			if !(i == precision || j == precision) {
				if i != 0 {
					indicies = append(indicies, uint32(k1))
					indicies = append(indicies, uint32(k2))
					indicies = append(indicies, uint32(k1+1))
				}
				if i != precision-1 {
					indicies = append(indicies, uint32(k1+1))
					indicies = append(indicies, uint32(k2))
					indicies = append(indicies, uint32(k2+1))
				}

			}

			k1 = k1 + 1
			k2 = k2 + 1
		}
	}
	return &Sphere{
		Points:    points,
		Indicies:  indicies,
		TexCoords: texCoords,
	}
}

// MaterialMeshInput method returns the verticies, indicies inputs for the NewMaterialMesh function.
func (s *Sphere) MaterialMeshInput() (vertex.Verticies, []uint32) {
	var verticies vertex.Verticies
	for i := 0; i < len(s.Points); i++ {
		verticies = append(verticies, vertex.Vertex{
			Position: s.Points[i],
			Normal:   s.Points[i],
		})
	}
	return verticies, s.Indicies
}
