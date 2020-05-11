package vertex

import "github.com/go-gl/mathgl/mgl32"

type Vertex struct {
	// Position vector
	Position mgl32.Vec3
	// Normal vector
	Normal mgl32.Vec3
	// Texture coordinates for textured objects.
	// As the textures are 2D, we need vec2 for storing the coordinates.
	TexCoords mgl32.Vec2
}

type Verticies []Vertex

func (v Verticies) VAO() []float32 {
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Normal.X())
		vao = append(vao, vertex.Normal.Y())
		vao = append(vao, vertex.Normal.Z())

		vao = append(vao, vertex.TexCoords.X())
		vao = append(vao, vertex.TexCoords.Y())
	}

	return vao
}
