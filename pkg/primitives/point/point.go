package point

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Shader interface {
	Use()
	SetUniformMat4(string, mgl32.Mat4)
}

type Point struct {
	coordinate mgl32.Vec3
	color      mgl32.Vec3
	size       float32

	direction mgl32.Vec3
	speed     float32
}

type Points struct {
	vao    *vao.VAO
	shader Shader
	points []*Point
}

// SetColor updates the Color of the point.
func (p *Point) SetColor(color mgl32.Vec3) {
	p.color = color
}

// SetSpeed updates the speed of the point.
func (p *Point) SetSpeed(speed float32) {
	p.speed = speed
}

// SetDirection updates the direction vector.
func (p *Point) SetDirection(dir mgl32.Vec3) {
	p.direction = dir
}

// SetIndexDirection updates the direction vector.
func (p *Point) SetIndexDirection(index int, value float32) {
	p.direction[index] = value
}

// Update is responsible for updating the state of the point
func (p *Point) Update(dt float64) {
	delta := float32(dt)
	motionVector := p.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * p.speed)
	}
	p.coordinate = (p.coordinate).Add(motionVector)
}

func New(shader Shader) *Points {
	return &Points{
		vao:    vao.NewVAO(),
		shader: shader,
		points: []*Point{},
	}
}

// Add inserts a new point to the points. The inputs: coordinate, color, size.
// It returns the point for further processing (eg: setup direction & speed.)
func (p *Points) Add(coords, color mgl32.Vec3, size float32) *Point {
	point := &Point{
		coordinate: coords,
		color:      color,
		size:       size,

		direction: mgl32.Vec3{0, 0, 0},
		speed:     0.0,
	}
	p.points = append(p.points, point)
	return point
}

// Update calls update for each points
func (p *Points) Update(dt float64) {
	for i, _ := range p.points {
		p.points[i].Update(dt)
	}
}

// Log is the string representation of the object
func (p *Points) Log() string {
	logString := "Points:\n"
	for _, item := range p.points {
		logString += " - Coordinate: Vector{" + trans.Vec3ToString(item.coordinate) + "}, color: Vector{" + trans.Vec3ToString(item.color) + "}, size: " + trans.Float32ToString(item.size) + "\n"
	}
	return logString
}
func (p *Points) setupVao() {
	p.vao.Clear()
	for index, _ := range p.points {
		p.vao.AppendPoint(p.points[index].coordinate, p.points[index].color, p.points[index].size)
	}
}
func (p *Points) buildVao() {
	p.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(p.vao.Get()), gl.Ptr(p.vao.Get()), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*7, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*7, gl.PtrOffset(4*3))
	// setup size
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 1, gl.FLOAT, false, 4*7, gl.PtrOffset(4*6))
}
func (p *Points) DrawWithUniforms(view, projection mgl32.Mat4) {
	p.shader.Use()
	p.shader.SetUniformMat4("view", view)
	p.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	p.shader.SetUniformMat4("model", M)
	p.draw()
}
func (p *Points) Draw() {
	p.shader.Use()
	p.draw()
}
func (p *Points) draw() {
	p.buildVao()

	gl.DrawArrays(gl.POINTS, 0, int32(len(p.vao.Get())/7))
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.BindVertexArray(0)
}

func (p *Points) Count() int {
	return len(p.points)
}
