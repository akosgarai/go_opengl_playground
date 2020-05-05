package cuboid

import (
	"github.com/go-gl/mathgl/mgl32"

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

	color mgl32.Vec3
}

func (c *Cuboid) Log() string {
	logString := "Cuboid:\n"
	logString += " - Rotation : Axis: Vector{" + trans.Vec3ToString(c.axis) + "}, angle: " + trans.Float32ToString(c.angle) + "}\n"
	logString += " - Color: Vector{" + trans.Vec3ToString(c.color) + "}, DrawMode: " + trans.IntegerToString(c.drawMode) + "\n"
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
	bottomPoints := bottom.Coordinates()
	// normal vector * (-1) * heightLength is the height vector.
	// each point of the bottom + calculated prev vector -> top rectangle.
	v1 := bottomPoints[1].Sub(bottomPoints[0])
	v2 := bottomPoints[3].Sub(bottomPoints[0])
	cp := v1.Cross(v2).Normalize()
	height := cp.Mul(-1.0 * heightLength)

	var sides [6]*rectangle.Rectangle
	// bottom
	sides[0] = bottom
	// top
	topSide := [4]mgl32.Vec3{
		bottomPoints[0].Add(height),
		bottomPoints[1].Add(height),
		bottomPoints[2].Add(height),
		bottomPoints[3].Add(height),
	}
	top := rectangle.New(topSide, bottom.Colors(), shader)
	sides[1] = top
	topPoints := top.Coordinates()
	// front
	frontSide := [4]mgl32.Vec3{
		topPoints[0],
		bottomPoints[0],
		bottomPoints[1],
		topPoints[1],
	}
	front := rectangle.New(frontSide, bottom.Colors(), shader)
	sides[2] = front
	// back
	backSide := [4]mgl32.Vec3{
		topPoints[3],
		bottomPoints[3],
		bottomPoints[2],
		topPoints[2],
	}
	back := rectangle.New(backSide, bottom.Colors(), shader)
	sides[3] = back
	// left
	leftSide := [4]mgl32.Vec3{
		topPoints[3],
		bottomPoints[3],
		bottomPoints[0],
		topPoints[0],
	}
	left := rectangle.New(leftSide, bottom.Colors(), shader)
	sides[4] = left
	// right
	rightSide := [4]mgl32.Vec3{
		topPoints[1],
		bottomPoints[1],
		bottomPoints[2],
		topPoints[2],
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
		color:    (bottom.Colors())[0],
	}
}

// SetColor updates every color with the given one.
func (c *Cuboid) SetColor(color mgl32.Vec3) {
	c.color = color
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
		c.shader.SetUniform3f("objectColor", c.color.X(), c.color.Y(), c.color.Z())
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
