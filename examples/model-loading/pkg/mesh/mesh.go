package mesh

import (
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/texture"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/vertex"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"

	"github.com/go-gl/mathgl/mgl32"
)

type Shader interface {
	GetId() uint32
	SetUniformMat4(string, mgl32.Mat4)
	SetUniform1f(string, float32)
	SetUniform1i(string, int32)
}
type Mesh struct {
	Verticies vertex.Verticies
	Textures  texture.Textures
	Indicies  []uint32

	vbo uint32
	ebo uint32
	vao uint32
}

func New(v []vertex.Vertex, i []uint32, t texture.Textures) *Mesh {
	mesh := &Mesh{
		Verticies: v,
		Textures:  t,
		Indicies:  i,
	}
	mesh.setup()
	return mesh
}

func (m *Mesh) Draw(shader Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-wrapper.TEXTURE0))
	}
	shader.SetUniformMat4("model", mgl32.Ident4())
	shader.SetUniform1f("material.shininess", float32(32))
	wrapper.BindVertexArray(m.vao)
	wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.Textures.UnBind()
	wrapper.BindVertexArray(0)
	wrapper.ActiveTexture(0)
}

func (m *Mesh) setup() {
	m.vao = wrapper.GenVertexArrays()
	m.vbo = wrapper.GenBuffers()
	m.ebo = wrapper.GenBuffers()

	wrapper.BindVertexArray(m.vao)

	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, m.vbo)
	wrapper.ArrayBufferData(m.Verticies.Get())

	wrapper.BindBuffer(wrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	wrapper.VertexAttribPointer(0, 3, wrapper.FLOAT, false, 4*8, wrapper.PtrOffset(0))
	// setup normals
	wrapper.VertexAttribPointer(1, 3, wrapper.FLOAT, false, 4*8, wrapper.PtrOffset(4*3))
	// setup texture position
	wrapper.VertexAttribPointer(2, 2, wrapper.FLOAT, false, 4*8, wrapper.PtrOffset(4*6))

	// close
	wrapper.BindVertexArray(0)
}
