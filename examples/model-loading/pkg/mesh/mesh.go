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

	// the center position of the mesh. the model transformation is calculated based on this.
	position mgl32.Vec3
	// movement paramteres
	direction mgl32.Vec3
	velocity  float32
	// rotation parameters
	// angle has to be in radian
	angle float32
	axis  mgl32.Vec3
	// for scaling - if a want to make other rectangles than unit ones.
	// This vector contains the scale factor for each axis.
	scale mgl32.Vec3
}

func New(v []vertex.Vertex, i []uint32, t texture.Textures) *Mesh {
	mesh := &Mesh{
		Verticies: v,
		Textures:  t,
		Indicies:  i,

		position:  mgl32.Vec3{0, 0, 0},
		direction: mgl32.Vec3{0, 0, 0},
		velocity:  0,
		angle:     0,
		axis:      mgl32.Vec3{0, 0, 0},
		scale:     mgl32.Vec3{1, 1, 1},
	}
	mesh.setup()
	return mesh
}
func (m *Mesh) SetScale(s mgl32.Vec3) {
	m.scale = s
}
func (m *Mesh) SetRotationAngle(a float32) {
	m.angle = a
}
func (m *Mesh) SetRotationAxis(a mgl32.Vec3) {
	m.axis = a
}
func (m *Mesh) SetPosition(p mgl32.Vec3) {
	m.position = p
}

func (m *Mesh) Draw(shader Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-wrapper.TEXTURE0))
	}
	M := m.modelTransformation()
	shader.SetUniformMat4("model", M)
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
func (m *Mesh) Update(dt float64) {
	delta := float32(dt)
	motionVector := m.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * m.velocity)
	}
	m.position = m.position.Add(motionVector)
}
func (m *Mesh) modelTransformation() mgl32.Mat4 {
	return mgl32.Translate3D(
		m.position.X(),
		m.position.Y(),
		m.position.Z()).Mul4(mgl32.HomogRotate3D(m.angle, m.axis)).Mul4(mgl32.Scale3D(
		m.scale.X(),
		m.scale.Y(),
		m.scale.Z(),
	))
}
