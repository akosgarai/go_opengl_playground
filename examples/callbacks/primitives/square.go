package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Square struct {
	A             Point
	B             Point
	C             Point
	D             Point
	precision     int
	vao           *VAO
	shaderProgram uint32
}

func NewSquare(a, b, c, d Point) *Square {
	return &Square{
		A:         a,
		B:         b,
		C:         c,
		D:         d,
		precision: 10,
		vao:       NewVAO(),
	}
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

// SetShaderProgram updates the shaderProgram of the square.
func (s *Square) SetShaderProgram(p uint32) {
	s.shaderProgram = p
}

func (s *Square) setupVao() {
	s.vao.Clear()
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
			s.vao.AppendPoint(a)
			s.vao.AppendPoint(b)
			s.vao.AppendPoint(c)
			s.vao.AppendPoint(a)
			s.vao.AppendPoint(c)
			s.vao.AppendPoint(d)
		}
	}
}

func (s *Square) Draw() {
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
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vao.Get())/6))
}
