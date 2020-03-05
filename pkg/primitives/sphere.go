package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sphere struct {
	center Vector
	radius float64
	color  Vector

	numOfRows       int
	numOfItemsInRow int
}

func NewSphere() *Sphere {
	return &Sphere{Vector{0, 0, 0}, 1, Vector{1, 1, 1}, 80, 80}
}

// SetCenter updates the center of the sphere
func (s *Sphere) SetCenter(c Vector) {
	s.center = c
}

// GetCenter returns the center of the sphere
func (s *Sphere) GetCenter() Vector {
	return s.center
}

// SetColor updates the color of the sphere
func (s *Sphere) SetColor(c Vector) {
	s.color = c
}

// GetColor returns the color of the sphere
func (s *Sphere) GetColor() Vector {
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
func (s *Sphere) appendPointToVao(currentVao []float32, p Point) []float32 {
	currentVao = append(currentVao, float32(p.Coordinate.X))
	currentVao = append(currentVao, float32(p.Coordinate.Y))
	currentVao = append(currentVao, float32(p.Coordinate.Z))
	currentVao = append(currentVao, float32(p.Color.X))
	currentVao = append(currentVao, float32(p.Color.Y))
	currentVao = append(currentVao, float32(p.Color.Z))
	return currentVao
}
func (s *Sphere) sideByPointToVao(currentVao []float32, pa, pb, pc, pd Point) []float32 {
	currentVao = s.appendPointToVao(currentVao, pa)
	currentVao = s.appendPointToVao(currentVao, pb)
	currentVao = s.appendPointToVao(currentVao, pc)
	currentVao = s.appendPointToVao(currentVao, pa)
	currentVao = s.appendPointToVao(currentVao, pc)
	currentVao = s.appendPointToVao(currentVao, pd)
	return currentVao
}
func (s *Sphere) setupVao() []float32 {
	var vao []float32
	// the coordinates will be set as a following: origo as center, 1 as radius, for drawing, the translation and scale could be done later in the model transformation.
	// Sphere top: center + v{0,radius,0}, bottom: center + v{0,-radius,0}, left: center + v{-radius,0,0}, right: center + v{radius,0,0}
	// Idea : start drawing triangles from both direction (top, bottom). step the coordinates and calculate the triangles, add them to vao.
	// - step for y coord, : radius * 2 / numOfRows
	RefPoint := &Vector{0, 1, 0}
	unitXRotation := (float64(360.0 / s.numOfItemsInRow))
	unitYRotation := (float64(360.0 / s.numOfRows))
	for i := 0; i < s.numOfRows; i++ {
		for j := 0; j < s.numOfItemsInRow; j++ {
			// define 4 points. ref.rotate?(i*?).rotate(j*?), ref.rotate((i+1)*?).rotate(j*?), ref.rotate((i+1)*?).rotate((j+1)*?), ref.rotate(i*?).rateate((j+1)*?)
			p1 := Point{*(RotationXMatrix4x4(float64(i) * unitXRotation).Mul4(RotationYMatrix4x4(float64(j) * unitYRotation)).MultiVector(*RefPoint)), s.color}
			p2 := Point{*(RotationXMatrix4x4(float64(i+1) * unitXRotation).Mul4(RotationYMatrix4x4(float64(j) * unitYRotation)).MultiVector(*RefPoint)), s.color}
			p3 := Point{*(RotationXMatrix4x4(float64(i+1) * unitXRotation).Mul4(RotationYMatrix4x4(float64(j+1) * unitYRotation)).MultiVector(*RefPoint)), s.color}
			p4 := Point{*(RotationXMatrix4x4(float64(i) * unitXRotation).Mul4(RotationYMatrix4x4(float64(j+1) * unitYRotation)).MultiVector(*RefPoint)), s.color}
			vao = s.sideByPointToVao(vao, p1, p4, p3, p2)
		}
	}
	return vao
}

func (s *Sphere) Draw() {
	points := s.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

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
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points)/6))
}
