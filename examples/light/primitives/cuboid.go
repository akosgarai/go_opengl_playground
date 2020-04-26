package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Cuboid struct {
	precision     int
	vao           *VAO
	shaderProgram uint32

	material *Material
	sides    [6]*Rectangle
}

func NewCuboid(bottom *Rectangle, heightLength float32, mat *Material, prec int, shaderProgram uint32) *Cuboid {
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
	top := NewRectangle(topSide, mat, prec, shaderProgram)
	top.SetInvertNormal(!bottom.IsNormalInverted())
	sides[1] = top
	// front
	frontSide := [4]mgl32.Vec3{
		top.Points[0],
		bottom.Points[0],
		bottom.Points[1],
		top.Points[1],
	}
	front := NewRectangle(frontSide, mat, prec, shaderProgram)
	front.SetInvertNormal(true)
	sides[2] = front
	// back
	backSide := [4]mgl32.Vec3{
		top.Points[3],
		bottom.Points[3],
		bottom.Points[2],
		top.Points[2],
	}
	back := NewRectangle(backSide, mat, prec, shaderProgram)
	sides[3] = back
	// left
	leftSide := [4]mgl32.Vec3{
		top.Points[3],
		bottom.Points[3],
		bottom.Points[0],
		top.Points[0],
	}
	left := NewRectangle(leftSide, mat, prec, shaderProgram)
	sides[4] = left
	// right
	rightSide := [4]mgl32.Vec3{
		top.Points[1],
		bottom.Points[1],
		bottom.Points[2],
		top.Points[2],
	}
	right := NewRectangle(rightSide, mat, prec, shaderProgram)
	right.SetInvertNormal(true)
	sides[5] = right

	return &Cuboid{
		precision:     prec,
		vao:           NewVAO(),
		shaderProgram: shaderProgram,
		material:      mat,
		sides:         sides,
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
	c.vao.Clear()
	for i := 0; i < 6; i++ {
		c.vao = c.sides[i].SetupExternalVao(c.vao)
	}
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
	gl.UseProgram(c.shaderProgram)

	viewLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
	projectionLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])
	modelLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("model\x00"))
	M := mgl32.Ident4()
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])

	// diffuse color
	diffuseLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("diffuseColor\x00"))
	diffCol := c.material.GetDiffuse()
	gl.Uniform3f(diffuseLocation, diffCol.X(), diffCol.Y(), diffCol.Z())
	// specular color
	specularLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("specularColor\x00"))
	specCol := c.material.GetSpecular()
	gl.Uniform3f(specularLocation, specCol.X(), specCol.Y(), specCol.Z())
	// shininess
	shininessLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("shininess\x00"))
	gl.Uniform1f(shininessLocation, c.material.GetShininess())
	// light position
	lightPosLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("lightPosition\x00"))
	gl.Uniform3f(lightPosLocation, lightPos.X(), lightPos.Y(), lightPos.Z())
	// normal matrix
	normalMatLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("normal\x00"))
	N := mgl32.Mat4Normal(M.Mul4(view))
	gl.UniformMatrix3fv(normalMatLocation, 1, false, &N[0])

	c.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(c.vao.Get())/6))
}
