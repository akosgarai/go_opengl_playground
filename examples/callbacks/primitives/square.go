package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Square struct {
	A         Point
	B         Point
	C         Point
	D         Point
	precision int
}

func NewSquare(a, b, c, d Point) *Square {
	return &Square{a, b, c, d, 10}
}

func (s *Square) SetColor(color mgl32.Vec3) {
	s.A.SetColor(color)
	s.B.SetColor(color)
	s.C.SetColor(color)
	s.D.SetColor(color)
}
func (s *Square) SetPrecision(p int) {
	s.precision = p
}

func (s *Square) appendPointToVao(currentVao []float32, p Point) []float32 {
	currentVao = append(currentVao, p.Coordinate.X())
	currentVao = append(currentVao, p.Coordinate.Y())
	currentVao = append(currentVao, p.Coordinate.Z())
	currentVao = append(currentVao, p.Color.X())
	currentVao = append(currentVao, p.Color.Y())
	currentVao = append(currentVao, p.Color.Z())
	return currentVao
}
func (s *Square) setupVao() []float32 {
	var points []float32
	verticalStep := s.B.Coordinate.Sub(s.A.Coordinate).Mul(1.0 / float32(s.precision))
	horisontalStep := s.D.Coordinate.Sub(s.A.Coordinate).Mul(1.0 / float32(s.precision))

	for horisontalLoopIndex := 0; horisontalLoopIndex < s.precision; horisontalLoopIndex++ {
		for verticalLoopIndex := 0; verticalLoopIndex < s.precision; verticalLoopIndex++ {
			a := Point{
				s.A.Coordinate.Add(
					verticalStep.Mul(float32(verticalLoopIndex))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex))),
				s.A.Color,
			}
			b := Point{
				s.A.Coordinate.Add(
					verticalStep.Mul(float32(verticalLoopIndex))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex + 1))),
				s.B.Color,
			}
			c := Point{
				s.A.Coordinate.Add(
					verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex + 1))),
				s.C.Color,
			}
			d := Point{
				s.A.Coordinate.Add(
					verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex))),
				s.D.Color,
			}
			points = s.appendPointToVao(points, a)
			points = s.appendPointToVao(points, b)
			points = s.appendPointToVao(points, c)
			points = s.appendPointToVao(points, a)
			points = s.appendPointToVao(points, c)
			points = s.appendPointToVao(points, d)
		}
	}

	return points
}

func (s *Square) Draw() {
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
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(points)/6))
}
