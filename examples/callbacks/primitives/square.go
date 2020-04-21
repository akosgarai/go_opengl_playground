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

// Log returns the string representation of this object.
func (s *Square) Log() string {
	logString := "Square:\n"
	logString += " - A : Vector{" + Vec3ToString(s.A.Coordinate) + "}\n"
	logString += " - B : Vector{" + Vec3ToString(s.B.Coordinate) + "}\n"
	logString += " - C : Vector{" + Vec3ToString(s.C.Coordinate) + "}\n"
	logString += " - D : Vector{" + Vec3ToString(s.D.Coordinate) + "}\n"
	logString += " - precision : " + string(s.precision) + "\n"
	return logString
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
	// Create the vao object
	s.setupVao()
	// setup shader program during draw.
	gl.UseProgram(s.shaderProgram)

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
	// draw the stuff.
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vao.Get())/6))
}

func (s *Square) DrawWithUniforms(view, projection mgl32.Mat4) {
	gl.UseProgram(s.shaderProgram)

	viewLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
	projectionLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	modelLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("model\x00"))

	M := mgl32.Translate3D(s.D.Coordinate.X(), s.D.Coordinate.Y(), s.D.Coordinate.Z())
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])

	s.setupVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
