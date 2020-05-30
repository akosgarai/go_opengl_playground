package mesh

import (
	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"
	"github.com/akosgarai/opengl_playground/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	Verticies vertex.Verticies
	Indicies  []uint32

	vbo uint32
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
	// For calling gl functions.
	wrapper interfaces.GLWrapper
	// parent-child hierarchy
	parent    interfaces.Mesh
	parentSet bool
}

// SetScale updates the scale of the mesh.
func (m *Mesh) SetScale(s mgl32.Vec3) {
	m.scale = s
}

// SetRotationAngle updates the rotation angle of the mesh. The input value
// has to be radian.
func (m *Mesh) SetRotationAngle(a float32) {
	m.angle = a
}

// GetRotationAngle returns the rotation angle of the mesh. The output value
// is radian.
func (m *Mesh) GetRotationAngle() float32 {
	return m.angle
}

// SetRotationAxis updates the rotation axis of the mesh.
func (m *Mesh) SetRotationAxis(a mgl32.Vec3) {
	m.axis = a
}

// SetPosition updates the position of the mesh.
func (m *Mesh) SetPosition(p mgl32.Vec3) {
	m.position = p
}

// SetDirection updates the move direction of the mesh.
func (m *Mesh) SetDirection(p mgl32.Vec3) {
	m.direction = p
}

// SetSpeed updates the velocity of the mesh.
func (m *Mesh) SetSpeed(a float32) {
	m.velocity = a
}

// GetPosition returns the current position of the mesh.
func (m *Mesh) GetPosition() mgl32.Vec3 {
	return m.position
}

// GetDirection returns the current direction of the mesh.
func (m *Mesh) GetDirection() mgl32.Vec3 {
	return m.direction
}
func (m *Mesh) SetParent(msh interfaces.Mesh) {
	m.parentSet = true
	m.parent = msh
}
func (m *Mesh) GetParentTranslationTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.TranslationTransformation()
	}
	return mgl32.Ident4()
}
func (m *Mesh) GetParentRotationTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.RotationTransformation()
	}
	return mgl32.Ident4()
}
func (m *Mesh) GetParentScaleTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.ScaleTransformation()
	}
	return mgl32.Ident4()
}

// Update calulates the position change. It's input is the delta since the current draw circle.
// The movement is calculated from the direction, velocity and delta.
// motion = motionVector * (delta * velocity)
// new position = current position + motion
func (m *Mesh) Update(dt float64) {
	delta := float32(dt)
	motionVector := m.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * m.velocity)
	}
	m.position = m.position.Add(motionVector)
}

// ModelTransformation returns the transformation that we can
// use as the model transformation of this mesh.
// The matrix is calculated from the position (translate), the rotation (rotate)
// and from the scale (scale) patameters.
func (m *Mesh) ModelTransformation() mgl32.Mat4 {
	return m.TranslationTransformation().Mul4(
		m.RotationTransformation()).Mul4(
		m.ScaleTransformation())
}

// ScaleTransformation returns the scale part of the model transformation.
func (m *Mesh) ScaleTransformation() mgl32.Mat4 {
	return mgl32.Scale3D(m.scale.X(), m.scale.Y(), m.scale.Z()).Mul4(m.GetParentScaleTransformation())
}

// TranslateTransformation returns the translation part of the model transformation.
func (m *Mesh) TranslationTransformation() mgl32.Mat4 {
	translate := m.position
	return mgl32.Translate3D(translate.X(), translate.Y(), translate.Z()).Mul4(m.GetParentTranslationTransformation())
}

// RotationTransformation returns the rotation part of the model transformation.
// It is used in the export module, where we have to handle the normal vectors also.
func (m *Mesh) RotationTransformation() mgl32.Mat4 {
	return mgl32.HomogRotate3D(m.angle, m.axis).Mul4(m.GetParentRotationTransformation())
}

func (m *Mesh) Rotate(angleDeg float32, axisVector mgl32.Vec3) {
	trMat := mgl32.HomogRotate3D(mgl32.DegToRad(angleDeg), axisVector)
	if m.parentSet {
		m.position = mgl32.TransformNormal(m.position, trMat)
	} else {
		m.direction = mgl32.TransformNormal(m.direction, trMat)
		m.SetRotationAngle(m.angle + mgl32.DegToRad(angleDeg))
		m.SetRotationAxis(axisVector)
	}
}
func (m *Mesh) IsParentMesh() bool {
	return !m.parentSet
}

type TexturedMesh struct {
	Mesh
	Indicies []uint32
	Textures texture.Textures
	ebo      uint32
}

func (m *TexturedMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_NORMAL_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it sets up the model uniform, and the shininess.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *TexturedMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	shader.SetUniform1f("material.shininess", float32(32))
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

// NewTexturedMesh gets the verticies, indicies, textures, glwrapper as inputs and makes the necessary setup for a
// standing (not moving) textured mesh before returning it. The vbo, vao, ebo is also set.
func NewTexturedMesh(v []vertex.Vertex, i []uint32, t texture.Textures, wrapper interfaces.GLWrapper) *TexturedMesh {
	mesh := &TexturedMesh{
		Mesh: Mesh{
			Verticies: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
		Indicies: i,
		Textures: t,
	}
	mesh.setup()
	return mesh
}

type MaterialMesh struct {
	Mesh
	Indicies []uint32
	Material *material.Material
	ebo      uint32
}

// NewMaterialMesh gets the verticies, indicies, material, glwrapper as inputs and makes the necessary setup for a
// standing (not moving) material mesh before returning it. The vbo, vao, ebo is also set.
func NewMaterialMesh(v []vertex.Vertex, i []uint32, mat *material.Material, wrapper interfaces.GLWrapper) *MaterialMesh {
	mesh := &MaterialMesh{
		Mesh: Mesh{
			Verticies: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
		Indicies: i,
		Material: mat,
	}
	mesh.setup()
	return mesh
}
func (m *MaterialMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_NORMAL))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(0))
	// setup normal vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(4*3))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the material with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). It also sets up the model uniform.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
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
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type PointMesh struct {
	Mesh
}

// NewPointMesh has only a glwrapper input, because it returns an empty mesh (without Verticies).
// Due to this, the vao, vbo setup is unnecessary now.
func NewPointMesh(wrapper interfaces.GLWrapper) *PointMesh {
	mesh := &PointMesh{
		Mesh{
			Verticies: []vertex.Vertex{},

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
	}
	return mesh
}
func (m *PointMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_COLOR_SIZE))

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(0))
	// setup color vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(4*3))
	// setup point size
	m.wrapper.VertexAttribPointer(2, 1, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the  model uniform with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function).
// Then it binds the vertex array and draws the mesh with arrays (points). Finally it cleans up.
func (m *PointMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawArrays(glwrapper.POINTS, 0, int32(len(m.Verticies)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

// AddVertex inserts a new vertex to the verticies. Then it calls setup
// Because the verticies are changed, so that we have to generate the vao again.
func (m *PointMesh) AddVertex(v vertex.Vertex) {
	m.Verticies.Add(v)
	m.setup()
}

type ColorMesh struct {
	Mesh
	Indicies []uint32
	Color    []mgl32.Vec3
	ebo      uint32
}

// NewColorMesh gets the verticies, indicies, colors, glwrapper as inputs and makes the necessary setup for a
// standing (not moving) colored mesh before returning it. The vbo, vao, ebo is also set.
func NewColorMesh(v []vertex.Vertex, i []uint32, color []mgl32.Vec3, wrapper interfaces.GLWrapper) *ColorMesh {
	mesh := &ColorMesh{
		Mesh: Mesh{
			Verticies: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
		Indicies: i,
		Color:    color,
	}
	mesh.setup()
	return mesh
}
func (m *ColorMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_COLOR))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(0))
	// setup color vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(4*3))

	// closeColorMesh
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the  model uniform with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function).
// Then it binds the vertex array and draws the mesh with arrays (points). Finally it cleans up.
func (m *ColorMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type TexturedColoredMesh struct {
	Mesh
	Indicies []uint32
	Textures texture.Textures
	Color    []mgl32.Vec3
	ebo      uint32
}

// NewTexturedColoredMesh gets the verticies, indicies, textures, colors, glwrapper as inputs and makes the necessary setup for a
// standing (not moving) textured colored mesh before returning it. The vbo, vao, ebo is also set.
func NewTexturedColoredMesh(v []vertex.Vertex, i []uint32, t texture.Textures, color []mgl32.Vec3, wrapper interfaces.GLWrapper) *TexturedColoredMesh {
	mesh := &TexturedColoredMesh{
		Mesh: Mesh{
			Verticies: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
		Indicies: i,
		Textures: t,
		Color:    color,
	}
	mesh.setup()
	return mesh
}
func (m *TexturedColoredMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_COLOR_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. Its input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it sets up the model uniform.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *TexturedColoredMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type TexturedMaterialMesh struct {
	Mesh
	Indicies []uint32
	Textures texture.Textures
	Material *material.Material
	ebo      uint32
}

// NewTexturedMaterialMesh gets the verticies, indicies, textures, material, glwrapper as inputs and makes the necessary setup for a
// standing (not moving) textured material mesh before returning it. The vbo, vao, ebo is also set.
func NewTexturedMaterialMesh(v []vertex.Vertex, i []uint32, t texture.Textures, material *material.Material, wrapper interfaces.GLWrapper) *TexturedMaterialMesh {
	mesh := &TexturedMaterialMesh{
		Mesh: Mesh{
			Verticies: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			angle:     0,
			axis:      mgl32.Vec3{0, 0, 0},
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,
		},
		Indicies: i,
		Textures: t,
		Material: material,
	}
	mesh.setup()
	return mesh
}
func (m *TexturedMaterialMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Verticies.Get(vertex.POSITION_NORMAL_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indicies)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. Its input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it binds the material and sets up the model uniform.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *TexturedMaterialMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
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
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indicies)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}
