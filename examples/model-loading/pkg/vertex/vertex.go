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
