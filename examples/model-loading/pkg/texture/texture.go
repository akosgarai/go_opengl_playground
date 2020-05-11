package texture

type Texture struct {
	// The id of the texture. eg gl.TEXTURE0
	Id uint32
	// The type of the texture. diffuse or specular.
	Type string
}
