package mesh

import (
	wrapper "github.com/akosgarai/opengl_playground/examples/model-loading/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/texture"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/vertex"
)

type Shader interface {
	GetId() uint32
}
type Mesh struct {
	Verticies vertex.Verticies
	Textures  []texture.Texture
	Indicies  []uint32

	vbo uint32
	ebo uint32
	vao uint32
}

func New(v []vertex.Vertex, i []uint32, t []texture.Texture) *Mesh {
	mesh := &Mesh{
		Verticies: v,
		Textures:  t,
		Indicies:  i,
	}
	mesh.setup()
	return mesh
}

func (m *Mesh) Draw(shader Shader) {
	for i := 0; i < len(m.Textures); i++ {
		tex := m.Textures[i]
		wrapper.ActiveTexture(uint32(i))
		wrapper.Uniform1i(wrapper.GetUniformLocation(tex.UniformName, shader.GetId()), int32(i))
		wrapper.BindTexture(wrapper.TEXTURE_2D, tex.Id)
	}
	wrapper.BindVertexArray(m.vao)
	wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	wrapper.BindVertexArray(0)
	wrapper.ActiveTexture(0)
}

func (m *Mesh) setup() {
	m.vao = wrapper.GenVertexArray()
	m.vbo = wrapper.GenBuffers()
	m.ebo = wrapper.GenBuffers()

	wrapper.BindVertexArray(m.vao)

	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, m.vbo)
	wrapper.ArrayBufferData(m.Verticies.Get())

	wrapper.BindBuffer(wrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	wrapper.VertexAttribPointer(0, 3, 4*8, 0)
	// setup normals
	wrapper.VertexAttribPointer(1, 3, 4*8, 4*3)
	// setup texture position
	wrapper.VertexAttribPointer(2, 2, 4*8, 4*6)

	// close
	wrapper.BindVertexArray(0)
}
