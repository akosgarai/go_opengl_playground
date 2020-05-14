package mesh

import (
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/texture"
	"github.com/akosgarai/opengl_playground/examples/model-loading/pkg/vertex"
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"

	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	Verticies vertex.Verticies
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
func (m *Mesh) SetDirection(p mgl32.Vec3) {
	m.direction = p
}
func (m *Mesh) SetSpeed(a float32) {
	m.velocity = a
}
func (m *Mesh) GetPosition() mgl32.Vec3 {
	return m.position
}
func (m *Mesh) GetDirection() mgl32.Vec3 {
	return m.direction
}

func (m *Mesh) Update(dt float64) {
	delta := float32(dt)
	motionVector := m.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * m.velocity)
	}
	m.position = m.position.Add(motionVector)
}
func (m *Mesh) ModelTransformation() mgl32.Mat4 {
	return mgl32.Translate3D(
		m.position.X(),
		m.position.Y(),
		m.position.Z()).Mul4(mgl32.HomogRotate3D(m.angle, m.axis)).Mul4(mgl32.Scale3D(
		m.scale.X(),
		m.scale.Y(),
		m.scale.Z(),
	))
}

type TexturedMesh struct {
	Mesh
	Textures texture.Textures
}

func (m *TexturedMesh) setup() {
	m.vao = wrapper.GenVertexArrays()
	m.vbo = wrapper.GenBuffers()
	m.ebo = wrapper.GenBuffers()

	wrapper.BindVertexArray(m.vao)

	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, m.vbo)
	wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_NORMAL_TEXCOORD))

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
func (m *TexturedMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-wrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	shader.SetUniform1f("material.shininess", float32(32))
	wrapper.BindVertexArray(m.vao)
	wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.Textures.UnBind()
	wrapper.BindVertexArray(0)
	wrapper.ActiveTexture(0)
}
func NewTexturedMesh(v []vertex.Vertex, i []uint32, t texture.Textures) *TexturedMesh {
	mesh := &TexturedMesh{
		Mesh{
			Verticies: v,
			Indicies:  i,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
		},
		t,
	}
	mesh.setup()
	return mesh
}

type MaterialMesh struct {
	Mesh
	Material *material.Material
}

func NewMaterialMesh(v []vertex.Vertex, i []uint32, mat *material.Material) *MaterialMesh {
	mesh := &MaterialMesh{
		Mesh{
			Verticies: v,
			Indicies:  i,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
		},
		mat,
	}
	mesh.setup()
	return mesh
}
func (m *MaterialMesh) setup() {
	m.vao = wrapper.GenVertexArrays()
	m.vbo = wrapper.GenBuffers()
	m.ebo = wrapper.GenBuffers()

	wrapper.BindVertexArray(m.vao)

	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, m.vbo)
	wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_NORMAL))

	wrapper.BindBuffer(wrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	wrapper.VertexAttribPointer(0, 3, wrapper.FLOAT, false, 4*6, wrapper.PtrOffset(0))
	// setup normal vector
	wrapper.VertexAttribPointer(1, 3, wrapper.FLOAT, false, 4*6, wrapper.PtrOffset(4*3))

	// close
	wrapper.BindVertexArray(0)
}
func (m *MaterialMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	diffuse := m.Material.GetDiffuse()
	ambient := m.Material.GetAmbient()
	specular := m.Material.GetSpecular()
	shininess := m.Material.GetShininess()
	shader.SetUniform3f("material.diffuse", diffuse.X(), diffuse.Y(), diffuse.Z())
	shader.SetUniform3f("material.ambient", ambient.X(), ambient.Y(), ambient.Z())
	shader.SetUniform3f("material.specular", specular.X(), specular.Y(), specular.Z())
	shader.SetUniform1f("material.shininess", shininess)
	wrapper.BindVertexArray(m.vao)
	wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	wrapper.BindVertexArray(0)
	wrapper.ActiveTexture(0)
}

type PointMesh struct {
	Mesh
}

func NewPointMesh() *PointMesh {
	mesh := &PointMesh{
		Mesh{
			Verticies: []vertex.Vertex{},
			Indicies:  []uint32{},

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
		},
	}
	return mesh
}
func (m *PointMesh) setup() {
	m.vao = wrapper.GenVertexArrays()
	m.vbo = wrapper.GenBuffers()

	wrapper.BindVertexArray(m.vao)

	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, m.vbo)
	wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_COLOR_SIZE))

	// setup coordinates
	wrapper.VertexAttribPointer(0, 3, wrapper.FLOAT, false, 4*7, wrapper.PtrOffset(0))
	// setup color vector
	wrapper.VertexAttribPointer(1, 3, wrapper.FLOAT, false, 4*7, wrapper.PtrOffset(4*3))
	// setup point size
	wrapper.VertexAttribPointer(2, 1, wrapper.FLOAT, false, 4*7, wrapper.PtrOffset(4*6))

	// close
	wrapper.BindVertexArray(0)
}
func (m *PointMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	wrapper.BindVertexArray(m.vao)
	wrapper.DrawArrays(wrapper.POINTS, 0, int32(len(m.Verticies)))

	wrapper.BindVertexArray(0)
	wrapper.ActiveTexture(0)
}
func (m *PointMesh) AddVertex(v vertex.Vertex) {
	m.Verticies.Add(v)
	m.setup()
}
