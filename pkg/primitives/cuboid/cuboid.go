package cuboid

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Cuboid struct {
	vao    *vao.VAO
	shader rectangle.Shader

	sides [6]*rectangle.Rectangle
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
		vao:    vao.NewVAO(),
		shader: shader,
		sides:  sides,
	}
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

func (c *Cuboid) setupVao() {
	c.vao.Clear()
	for i := 0; i < 6; i++ {
		c.vao = c.sides[i].SetupExternalVao(c.vao)
	}
}
func (c *Cuboid) buildVao() {
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
	c.draw()
}
func (c *Cuboid) draw() {
	c.buildVao()
	c.shader.DrawTriangles(int32(len(c.vao.Get()) / 6))
	c.shader.Close(1)
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (c *Cuboid) DrawWithUniforms(view, projection mgl32.Mat4) {
	c.shader.Use()
	c.shader.SetUniformMat4("view", view)
	c.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	c.shader.SetUniformMat4("model", M)
	c.draw()
}

// Update
func (c *Cuboid) Update(dt float64) {
	for i := 0; i < 6; i++ {
		c.sides[i].Update(dt)
	}
}
