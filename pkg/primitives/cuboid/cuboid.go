package cuboid

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

const (
	DRAW_MODE_COLOR = 0
	DRAW_MODE_LIGHT = 1
)

type Cuboid struct {
	vao    *vao.VAO
	shader rectangle.Shader

	sides [6]*rectangle.Rectangle
	// rotation parameters
	// angle has to be in radian
	angle    float32
	axis     mgl32.Vec3
	drawMode int

	material *material.Material
}

func (c *Cuboid) Log() string {
	logString := "Cuboid:\n"
	logString += " - Rotation : Axis: Vector{" + trans.Vec3ToString(c.axis) + "}, angle: " + trans.Float32ToString(c.angle) + "}, DrawMode: " + trans.IntegerToString(c.drawMode) + "\n"
	logString += " - " + c.material.Log() + "\n"
	logString += " - Center: Vector{" + trans.Vec3ToString(c.GetCenterPoint()) + "}\n"
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
	return logString

}

func New(bottom *rectangle.Rectangle, heightLength float32, shader rectangle.Shader) *Cuboid {
	givenSidePoints := bottom.Coordinates()
	// normal vector * (-1) * heightLength is the height vector.
	// each point of the bottom + calculated prev vector -> top rectangle.
	v1 := givenSidePoints[1].Sub(givenSidePoints[0])
	v2 := givenSidePoints[3].Sub(givenSidePoints[0])
	cp := v1.Cross(v2).Normalize()
	height := cp.Mul(-1.0 * heightLength)

	var sides [6]*rectangle.Rectangle
	// bottom
	sides[0] = bottom
	// top
	topSide := [4]mgl32.Vec3{
		givenSidePoints[0].Add(height),
		givenSidePoints[3].Add(height),
		givenSidePoints[2].Add(height),
		givenSidePoints[1].Add(height),
	}
	top := rectangle.New(topSide, bottom.Colors(), shader)
	sides[1] = top
	oppositeSidePoints := top.Coordinates()
	// front
	frontSide := [4]mgl32.Vec3{
		oppositeSidePoints[0],
		oppositeSidePoints[3],
		givenSidePoints[1],
		givenSidePoints[0],
	}
	front := rectangle.New(frontSide, bottom.Colors(), shader)
	sides[2] = front
	// back
	backSide := [4]mgl32.Vec3{
		givenSidePoints[2],
		oppositeSidePoints[2],
		oppositeSidePoints[1],
		givenSidePoints[3],
	}
	back := rectangle.New(backSide, bottom.Colors(), shader)
	sides[3] = back
	// left
	leftSide := [4]mgl32.Vec3{
		givenSidePoints[2],
		givenSidePoints[1],
		oppositeSidePoints[3],
		oppositeSidePoints[2],
	}
	left := rectangle.New(leftSide, bottom.Colors(), shader)
	sides[4] = left
	// right
	rightSide := [4]mgl32.Vec3{
		givenSidePoints[0],
		givenSidePoints[3],
		oppositeSidePoints[1],
		oppositeSidePoints[0],
	}
	right := rectangle.New(rightSide, bottom.Colors(), shader)
	sides[5] = right

	return &Cuboid{
		vao:      vao.NewVAO(),
		shader:   shader,
		sides:    sides,
		angle:    0,
		axis:     mgl32.Vec3{0, 0, 0},
		drawMode: DRAW_MODE_COLOR,
		material: material.New((bottom.Colors())[0], (bottom.Colors())[0], (bottom.Colors())[0], 36.0),
	}
}

// SetMaterial updates the material of the cuboid.
func (c *Cuboid) SetMaterial(mat *material.Material) {
	c.material = mat
}

// SetColor updates every color with the given one.
func (c *Cuboid) SetColor(color mgl32.Vec3) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetColor(color)
	}
}

// SetIndexColor updates the color of the given index.
func (c *Cuboid) SetIndexColor(index int, color mgl32.Vec3) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetIndexColor(index, color)
	}
}

// SetSideColor updates the color of the given index.
func (c *Cuboid) SetSideColor(index int, color mgl32.Vec3) {
	for i := 0; i < 6; i++ {
		c.sides[index].SetColor(color)
	}
}

// SetDirection updates the direction vector.
func (c *Cuboid) SetDirection(dir mgl32.Vec3) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetDirection(dir)
	}
}

// SetIndexDirection updates the direction vector.
func (c *Cuboid) SetIndexDirection(index int, value float32) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetIndexDirection(index, value)
	}
}

// SetSpeed updates the speed.
func (c *Cuboid) SetSpeed(speed float32) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetSpeed(speed)
	}
}

// SetPrecision updates the precision of the rectangles of the cuboid
func (c *Cuboid) SetPrecision(p int) {
	for i := 0; i < 6; i++ {
		c.sides[i].SetPrecision(p)
	}
}

// SetAngle updates the angle.
// Input has to be radian.
func (c *Cuboid) SetAngle(angle float32) {
	c.angle = angle
}

// SetAxis updates the axis.
func (c *Cuboid) SetAxis(axis mgl32.Vec3) {
	c.axis = axis
}

// GetDirection returns the direction of the cuboid, aka the direction of the first side.
func (c *Cuboid) GetDirection() mgl32.Vec3 {
	return c.sides[0].GetDirection()
}

// GetCenterPoint return the center point of the cuboid.
// In other words it returns the cross point of the diagonals
func (c *Cuboid) GetCenterPoint() mgl32.Vec3 {
	// calculate the diagonal:side[0][0] - side[1][2], then
	// multiply it with 1/2 (make it half), and add it to side[0][X].
	diagonalHalf := (c.sides[1].Coordinates()[2]).Sub(c.sides[0].Coordinates()[0]).Mul(0.5)
	return (c.sides[0].Coordinates()[0]).Add(diagonalHalf)
}

func (c *Cuboid) setupVao() {
	c.vao.Clear()
	for i := 0; i < 6; i++ {
		c.vao = c.sides[i].SetupExternalVao(c.vao)
	}
}
func (c *Cuboid) buildVaoWithTexture() {
	// Create the vao object
	c.setupVao()

	c.shader.BindBufferData(c.vao.Get())

	c.shader.BindVertexArray()
	// setup points
	c.shader.VertexAttribPointer(0, 3, 4*8, 0)
	// setup color
	c.shader.VertexAttribPointer(1, 3, 4*8, 4*3)
	c.shader.VertexAttribPointer(2, 2, 4*8, 4*6)
}
func (c *Cuboid) buildVaoWithoutTexture() {
	// Create the vao object
	c.setupVao()

	c.shader.BindBufferData(c.vao.Get())

	c.shader.BindVertexArray()
	// setup points
	c.shader.VertexAttribPointer(0, 3, 4*6, 0)
	// setup color
	c.shader.VertexAttribPointer(1, 3, 4*6, 4*3)
}

// Draw is for drawing the cuboid to the screen.
func (c *Cuboid) Draw() {
	c.shader.Use()

	c.setupColorUniform()
	if !c.shader.HasTexture() {
		c.drawWithoutTextures()
	} else {
		c.drawWithTextures()
	}
}
func (c *Cuboid) drawWithTextures() {
	c.buildVaoWithTexture()
	c.shader.DrawTriangles(int32(len(c.vao.Get()) / 8))
	c.shader.Close(2)
}
func (c *Cuboid) drawWithoutTextures() {
	c.buildVaoWithoutTexture()
	c.shader.DrawTriangles(int32(len(c.vao.Get()) / 6))
	c.shader.Close(1)
}

func (c *Cuboid) modelTransformation() mgl32.Mat4 {
	return mgl32.HomogRotate3D(c.angle, c.axis)
}
func (c *Cuboid) setupColorUniform() {
	if c.drawMode == DRAW_MODE_LIGHT {
		diffuse := c.material.GetDiffuse()
		ambient := c.material.GetAmbient()
		specular := c.material.GetSpecular()
		shininess := c.material.GetShininess()
		c.shader.SetUniform3f("material.diffuse", diffuse.X(), diffuse.Y(), diffuse.Z())
		c.shader.SetUniform3f("material.ambient", ambient.X(), ambient.Y(), ambient.Z())
		c.shader.SetUniform3f("material.specular", specular.X(), specular.Y(), specular.Z())
		c.shader.SetUniform1f("material.shininess", shininess)
	}
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (c *Cuboid) DrawWithUniforms(view, projection mgl32.Mat4) {
	c.shader.Use()
	c.shader.SetUniformMat4("view", view)
	c.shader.SetUniformMat4("projection", projection)
	M := c.modelTransformation()
	c.shader.SetUniformMat4("model", M)

	c.setupColorUniform()

	if !c.shader.HasTexture() {
		c.drawWithoutTextures()
	} else {
		c.drawWithTextures()
	}
}

// Update
func (c *Cuboid) Update(dt float64) {
	for i := 0; i < 6; i++ {
		c.sides[i].Update(dt)
	}
}

// DrawMode updates the draw mode after validation. If it fails, it keeps the original value.
func (c *Cuboid) DrawMode(mode int) {
	if mode != DRAW_MODE_COLOR && mode != DRAW_MODE_LIGHT {
		return
	}
	c.drawMode = mode
	for i := 0; i < 6; i++ {
		c.sides[i].DrawMode(mode)
	}
}
