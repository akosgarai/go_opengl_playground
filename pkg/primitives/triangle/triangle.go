package triangle

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Triangle struct {
	vao    *vao.VAO
	shader *shader.Shader

	points [3]mgl32.Vec3
	colors [3]mgl32.Vec3

	direction mgl32.Vec3
	speed     float32
}

func NewTriangle(points, colors [3]mgl32.Vec3, shader *shader.Shader) *Triangle {
	return &Triangle{
		shader:    shader,
		vao:       vao.NewVAO(),
		points:    points,
		colors:    colors,
		direction: mgl32.Vec3{0, 0, 0},
		speed:     0,
	}
}

// SetColor updates every color with the given one.
func (t *Triangle) SetColor(color mgl32.Vec3) {
	for i := 0; i < 3; i++ {
		t.colors[i] = color
	}
}

// SetIndexColor updates the color of the given index.
func (t *Triangle) SetIndexColor(index int, color mgl32.Vec3) {
	t.colors[index] = color
}

// SetDirection updates the direction vector.
func (t *Triangle) SetDirection(dir mgl32.Vec3) {
	t.direction = dir
}

// SetIndexDirection updates the direction vector.
func (t *Triangle) SetIndexDirection(index int, value float32) {
	t.direction[index] = value
}

// SetSpeed updates the speed.
func (t *Triangle) SetSpeed(speed float32) {
	t.speed = speed
}

// Log returns the string representation of this object.
func (t *Triangle) Log() string {
	logString := "Triangle:\n"
	logString += " - A : Coordinate: Vector{" + trans.Vec3ToString(t.points[0]) + "}, color: Vector{" + trans.Vec3ToString(t.colors[0]) + "}\n"
	logString += " - B : Coordinate: Vector{" + trans.Vec3ToString(t.points[1]) + "}, color: Vector{" + trans.Vec3ToString(t.colors[1]) + "}\n"
	logString += " - C : Coordinate: Vector{" + trans.Vec3ToString(t.points[2]) + "}, color: Vector{" + trans.Vec3ToString(t.colors[2]) + "}\n"
	logString += " - Movement : Direction: Vector{" + trans.Vec3ToString(t.direction) + "}, speed: " + trans.Float32ToString(t.speed) + "}\n"
	return logString
}
func (t *Triangle) setupVao() {
	t.vao.Clear()
	for i := 0; i < 3; i++ {
		t.vao.AppendVectors(t.points[i], t.colors[i])
	}
}

func (t *Triangle) buildVao() uint32 {
	t.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(t.vao.Get()), gl.Ptr(t.vao.Get()), gl.STATIC_DRAW)

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

	return vertexArrayObject
}
func (t *Triangle) Draw() {
	t.shader.Use()
	// setup MVP to ident4 matrix
	MVP := mgl32.Ident4()
	t.shader.SetUniformMat4("MVP", MVP)
	t.draw()
}
func (t *Triangle) draw() {
	t.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (t *Triangle) DrawWithUniforms(view, projection mgl32.Mat4) {
	t.shader.Use()
	t.shader.SetUniformMat4("view", view)
	t.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	t.shader.SetUniformMat4("model", M)

	t.draw()
}
func (t *Triangle) Update(dt float64) {
	delta := float32(dt)
	motionVector := t.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * t.speed)
	}
	for i := 0; i < 3; i++ {
		t.points[i] = (t.points[i]).Add(motionVector)
	}
}
