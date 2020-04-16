package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Sphere struct {
	center mgl32.Vec3
	radius float64
	color  mgl32.Vec3

	numOfRows       int
	numOfItemsInRow int
	vao             *VAO
}

func NewSphere() *Sphere {
	return &Sphere{mgl32.Vec3{0, 0, 0}, 1, mgl32.Vec3{1, 1, 1}, 20, 20, NewVAO()}
}

// SetCenter updates the center of the sphere
func (s *Sphere) SetCenter(c mgl32.Vec3) {
	s.center = c
}

// GetCenter returns the center of the sphere
func (s *Sphere) GetCenter() mgl32.Vec3 {
	return s.center
}

// SetColor updates the color of the sphere
func (s *Sphere) SetColor(c mgl32.Vec3) {
	s.color = c
}

// GetColor returns the color of the sphere
func (s *Sphere) GetColor() mgl32.Vec3 {
	return s.color
}

// SetRadius updates the radius of the sphere
func (s *Sphere) SetRadius(r float64) {
	s.radius = r
}

// GetRadius returns the radius of the sphere
func (s *Sphere) GetRadius() float64 {
	return s.radius
}
func (s *Sphere) setupVao() {
	s.vao.Clear()
	// the coordinates will be set as a following: origo as center, 1 as radius, for drawing, the translation and scale could be done later in the model transformation.
	// Sphere top: center + v{0,radius,0}, bottom: center + v{0,-radius,0}, left: center + v{-radius,0,0}, right: center + v{radius,0,0}
	RefPoint := &mgl32.Vec3{0, 1, 0}
	step_Z := -mgl32.DegToRad(float32(360.0 / float32(s.numOfItemsInRow)))
	step_Y := -mgl32.DegToRad(float32(360.0 / float32(s.numOfRows)))
	for i := 0; i < s.numOfRows; i++ {
		i_Rotation := mgl32.HomogRotate3DZ(float32(i) * step_Z).Transpose()
		i1_Rotation := mgl32.HomogRotate3DZ(float32(i+1) * step_Z).Transpose()
		for j := 0; j < s.numOfItemsInRow; j++ {
			j1_Rotation := mgl32.HomogRotate3DY(float32(j+1) * step_Y).Transpose()
			j_Rotation := mgl32.HomogRotate3DY(float32(j) * step_Y).Transpose()
			if i == 0 {
				p1 := Point{*RefPoint, s.color}
				p2 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i1_Rotation)), s.color}
				p3 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i1_Rotation)), s.color}
				s.vao.AppendTrianglePoints(p1, p2, p3)
			} else {
				p1 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i_Rotation)), s.color}
				p2 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i_Rotation)), s.color}
				p3 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i1_Rotation)), s.color}
				p4 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i1_Rotation)), s.color}
				s.vao.AppendTrianglePoints(p1, p2, p3)
				s.vao.AppendTrianglePoints(p1, p3, p4)
			}
		}
	}
}

func (s *Sphere) Draw() {
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
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// setup color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// The sphere is represented by triangles, so we have TODO points here.
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vao.Get())/6))
}
