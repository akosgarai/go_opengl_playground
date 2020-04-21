package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Rectangle struct {
	precision     int
	vao           *VAO
	shaderProgram uint32

	color  mgl32.Vec3
	points [4]mgl32.Vec3
}

// Log returns the string representation of this object.
func (r *Rectangle) Log() string {
	logString := "Rectangle:\n"
	logString += " - A : Vector{" + Vec3ToString(r.points[0]) + "}\n"
	logString += " - B : Vector{" + Vec3ToString(r.points[1]) + "}\n"
	logString += " - C : Vector{" + Vec3ToString(r.points[2]) + "}\n"
	logString += " - D : Vector{" + Vec3ToString(r.points[3]) + "}\n"
	logString += " - color : Vector{" + Vec3ToString(r.color) + "}\n"
	logString += " - precision : " + string(r.precision) + "\n"
	return logString
}

// SetColor updates the color of the rectangle
func (r *Rectangle) SetColor(c mgl32.Vec3) {
	r.color = c
}

// GetColor returns the color of the rectangle
func (r *Rectangle) GetColor() mgl32.Vec3 {
	return r.color
}

// SetPrecision updates the precision of the rectangle
func (r *Rectangle) SetPrecision(p int) {
	r.precision = p
}

// SetShaderProgram updates the shaderProgram of the rectangle.
func (r *Rectangle) SetShaderProgram(p uint32) {
	r.shaderProgram = p
}

func (r *Rectangle) setupVao() {
	r.vao.Clear()
	verticalStep := r.points[1].Sub(r.points[0]).Mul(1.0 / float32(r.precision))
	horisontalStep := r.points[3].Sub(r.points[0]).Mul(1.0 / float32(r.precision))

	for horisontalLoopIndex := 0; horisontalLoopIndex < r.precision; horisontalLoopIndex++ {
		for verticalLoopIndex := 0; verticalLoopIndex < r.precision; verticalLoopIndex++ {
			a := Point{
				r.points[0].Add(
					verticalStep.Mul(float32(verticalLoopIndex))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex))),
				r.color,
			}
			b := Point{
				r.points[0].Add(
					verticalStep.Mul(float32(verticalLoopIndex))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex + 1))),
				r.color,
			}
			c := Point{
				r.points[0].Add(
					verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex + 1))),
				r.color,
			}
			d := Point{
				r.points[0].Add(
					verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
					horisontalStep.Mul(float32(horisontalLoopIndex))),
				r.color,
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
func (r *Rectangle) buildVao() {
	// Create the vao object
	r.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(r.vao.Get()), gl.Ptr(r.vao.Get()), gl.STATIC_DRAW)

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
}

// Draw is for drawing the rectangle to the screen.
func (r *Rectangle) Draw() {
	gl.UseProgram(r.shaderProgram)
	r.buildVao()
	// draw the stuff.
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(r.vao.Get())/6))
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (r *Rectangle) DrawWithUniforms(view, projection mgl32.Mat4) {
	gl.UseProgram(r.shaderProgram)

	viewLocation := gl.GetUniformLocation(r.shaderProgram, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
	projectionLocation := gl.GetUniformLocation(r.shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	modelLocation := gl.GetUniformLocation(r.shaderProgram, gl.Str("model\x00"))

	M := mgl32.Translate3D(r.points[0].X(), r.points[0].Y(), r.points[0].Z())
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])

	r.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
