package vertex

import "github.com/go-gl/mathgl/mgl32"

const (
	POSITION_NORMAL          = 1
	POSITION_NORMAL_TEXCOORD = 2
	POSITION_COLOR_SIZE      = 3
)

type Vertex struct {
	// Position vector
	Position mgl32.Vec3
	// Normal vector
	Normal mgl32.Vec3
	// Texture coordinates for textured objects.
	// As the textures are 2D, we need vec2 for storing the coordinates.
	TexCoords mgl32.Vec2
	// Color vector
	Color mgl32.Vec3
	// Point size for points
	PointSize float32
}

type Verticies []Vertex

// Get returns the vertex array object as []float32. It's input describes the
// fields and the order that we want to add to the vao.
// Modes: 'POSITION_NORMAL' (1) the position and the normal vectors are returned in this order.
// 'POSITION_NORMAL_TEXCOORD' (2) the position, normal and the tex coords vectors are
// returned in this order. 'POSITION_COLOR_SIZE' (3) the position, the color and the point size
// are retuned in this order.
func (v Verticies) Get(resultMode int) []float32 {
	if resultMode == POSITION_COLOR_SIZE {
		return v.getPoint()
	}
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Normal.X())
		vao = append(vao, vertex.Normal.Y())
		vao = append(vao, vertex.Normal.Z())

		if resultMode == POSITION_NORMAL_TEXCOORD {
			vao = append(vao, vertex.TexCoords.X())
			vao = append(vao, vertex.TexCoords.Y())
		}
	}

	return vao
}
func (v Verticies) getPoint() []float32 {
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Color.X())
		vao = append(vao, vertex.Color.Y())
		vao = append(vao, vertex.Color.Z())

		vao = append(vao, vertex.PointSize)
	}

	return vao
}

// Add appends the given vertex to the verticies.
func (v *Verticies) Add(ver Vertex) {
	*v = append(*v, ver)
}
