package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Cuboid struct {
	precision int
	vao       *vao.VAO
	shader    *shader.Shader

	material *Material
	sides    [6]*Rectangle
	static   bool
	vaoIsSet bool
}

func NewCuboid(bottom *Rectangle, heightLength float32, mat *Material, prec int, shader *shader.Shader) *Cuboid {
	height := bottom.GetNormal().Mul(-1.0 * heightLength)

	var sides [6]*Rectangle
	// bottom
	sides[0] = bottom
	// top
	topSide := [4]mgl32.Vec3{
		bottom.Points[0].Add(height),
		bottom.Points[1].Add(height),
		bottom.Points[2].Add(height),
		bottom.Points[3].Add(height),
	}
	top := NewRectangle(topSide, mat, prec, shader)
	top.SetInvertNormal(!bottom.IsNormalInverted())
	sides[1] = top
	// front
	frontSide := [4]mgl32.Vec3{
		top.Points[0],
		bottom.Points[0],
		bottom.Points[1],
		top.Points[1],
	}
	front := NewRectangle(frontSide, mat, prec, shader)
	front.SetInvertNormal(true)
	sides[2] = front
	// back
	backSide := [4]mgl32.Vec3{
		top.Points[3],
		bottom.Points[3],
		bottom.Points[2],
		top.Points[2],
	}
	back := NewRectangle(backSide, mat, prec, shader)
	sides[3] = back
	// left
	leftSide := [4]mgl32.Vec3{
		top.Points[3],
		bottom.Points[3],
		bottom.Points[0],
		top.Points[0],
	}
	left := NewRectangle(leftSide, mat, prec, shader)
	sides[4] = left
	// right
	rightSide := [4]mgl32.Vec3{
		top.Points[1],
		bottom.Points[1],
		bottom.Points[2],
		top.Points[2],
	}
	right := NewRectangle(rightSide, mat, prec, shader)
	right.SetInvertNormal(true)
	sides[5] = right

	return &Cuboid{
		precision: prec,
		vao:       vao.NewVAO(),
		shader:    shader,
		material:  mat,
		sides:     sides,
		static:    true,
		vaoIsSet:  false,
	}
}
func (c *Cuboid) Log() string {
	logString := "Cuboid:\n"
	logString += " - Top:\n"
	logString += c.sides[1].Log()
	logString += " - Bottom:\n"
	logString += c.sides[0].Log()
	logString += " - Front:\n"
	logString += c.sides[2].Log()
	logString += " - Back:\n"
	logString += c.sides[3].Log()
	logString += " - Left:\n"
	logString += c.sides[4].Log()
	logString += " - Right:\n"
	logString += c.sides[5].Log()
	logString += " - precision : " + IntegerToString(c.precision) + "\n"
	logString += " - " + c.material.Log() + "\n"
	return logString

}
func (c *Cuboid) Update(dt float64) {
}
func (c *Cuboid) setupVao() {
	if c.static {
		if c.vaoIsSet {
			return
		}
	}
	c.vao.Clear()
	for i := 0; i < 6; i++ {
		c.vao = c.sides[i].SetupExternalVao(c.vao)
	}
	c.vaoIsSet = true
}
func (c *Cuboid) buildVao() {
	// Create the vao object
	c.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(c.vao.Get()), gl.Ptr(c.vao.Get()), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
}

// DrawWithLight is for drawing the rectangle to the screen. but with lightsource.
func (c *Cuboid) DrawWithLight(view, projection mgl32.Mat4, lightPos mgl32.Vec3) {
	c.shader.Use()

	c.shader.SetUniformMat4("view", view)
	c.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	c.shader.SetUniformMat4("model", M)

	// diffuse color
	diffCol := c.material.GetDiffuse()
	c.shader.SetUniform3f("diffuseColor", diffCol.X(), diffCol.Y(), diffCol.Z())
	// specular color
	specCol := c.material.GetSpecular()
	c.shader.SetUniform3f("specularColor", specCol.X(), specCol.Y(), specCol.Z())
	// shininess
	c.shader.SetUniform1f("shininess", c.material.GetShininess())
	// light position
	c.shader.SetUniform3f("lightPosition", lightPos.X(), lightPos.Y(), lightPos.Z())
	// normal matrix
	N := mgl32.Mat4Normal(M.Mul4(view))
	c.shader.SetUniformMat3("normal", N)

	c.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.vao.Get())/6))
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
}
