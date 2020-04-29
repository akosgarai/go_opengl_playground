package square

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Square struct {
	vao    *vao.VAO
	shader *shader.Shader

	points [4]mgl32.Vec3
	colors [4]mgl32.Vec3

	direction mgl32.Vec3
	speed     float32
}

func NewSquare(points, colors [4]mgl32.Vec3, shader *shader.Shader) *Square {
	return &Square{
		shader:    shader,
		vao:       vao.NewVAO(),
		points:    points,
		colors:    colors,
		direction: mgl32.Vec3{0, 0, 0},
		speed:     0,
	}
}

// Log returns the string representation of this object.
func (s *Square) Log() string {
	logString := "Triangle:\n"
	logString += " - A : Coordinate: Vector{" + trans.Vec3ToString(s.points[0]) + "}, color: Vector{" + trans.Vec3ToString(s.colors[0]) + "}\n"
	logString += " - B : Coordinate: Vector{" + trans.Vec3ToString(s.points[1]) + "}, color: Vector{" + trans.Vec3ToString(s.colors[1]) + "}\n"
	logString += " - C : Coordinate: Vector{" + trans.Vec3ToString(s.points[2]) + "}, color: Vector{" + trans.Vec3ToString(s.colors[2]) + "}\n"
	logString += " - D : Coordinate: Vector{" + trans.Vec3ToString(s.points[3]) + "}, color: Vector{" + trans.Vec3ToString(s.colors[3]) + "}\n"
	return logString
}

// SetColor updates every color with the given one.
func (s *Square) SetColor(color mgl32.Vec3) {
	for i := 0; i < 4; i++ {
		s.colors[i] = color
	}
}

// SetIndexColor updates the color of the given index.
func (s *Square) SetIndexColor(index int, color mgl32.Vec3) {
	s.colors[index] = color
}

// SetDirection updates the direction vector.
func (s *Square) SetDirection(dir mgl32.Vec3) {
	s.direction = dir
}

// SetIndexDirection updates the direction vector.
func (s *Square) SetIndexDirection(index int, value float32) {
	s.direction[index] = value
}

// SetSpeed updates the speed.
func (s *Square) SetSpeed(speed float32) {
	s.speed = speed
}

func (s *Square) setupVao() {
	s.vao.Clear()
	indicies := [6]int{0, 1, 2, 0, 2, 3}
	for i := 0; i < 6; i++ {
		s.vao.AppendVectors(s.points[indicies[i]], s.colors[indicies[i]])
	}
}
func (s *Square) buildVao() uint32 {
	s.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.vao.Get()), gl.Ptr(s.vao.Get()), gl.STATIC_DRAW)

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

func (s *Square) Draw() {
	s.shader.Use()
	// setup MVP to ident4 matrix
	MVP := mgl32.Ident4()
	s.shader.SetUniformMat4("MVP", MVP)
	s.draw()
}
func (s *Square) draw() {
	s.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (s *Square) DrawWithUniforms(view, projection mgl32.Mat4) {
	s.shader.Use()
	s.shader.SetUniformMat4("view", view)
	s.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	s.shader.SetUniformMat4("model", M)

	s.draw()
}
func (s *Square) Update(dt float64) {
	delta := float32(dt)
	motionVector := s.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * s.speed)
	}
	for i := 0; i < 4; i++ {
		s.points[i] = s.points[i].Add(motionVector)
	}
}
