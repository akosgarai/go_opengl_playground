package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Cube struct {
	A Point // bottom-left-back
	B Point // bottom-right-back
	C Point // top-right-back
	D Point // top-left-back
	E Point // top-left-front
	F Point // top-right-front
	G Point // bottom-right-front
	H Point // bottom-left-front

	vao           *VAO
	shaderProgram uint32
}

func NewCubeByPoints(a, b, c, d, e, f, g, h Point) *Cube {
	return &Cube{
		A:   a,
		B:   b,
		C:   c,
		D:   d,
		E:   e,
		F:   f,
		G:   g,
		H:   h,
		vao: NewVAO(),
	}
}
func NewCubeByVectorAndLength(in mgl32.Vec3, sideLength float32) *Cube {
	color1 := mgl32.Vec3{1, 0, 0}
	color2 := mgl32.Vec3{0, 1, 0}
	a := Point{in, color1}
	b := Point{a.Coordinate.Add(mgl32.Vec3{sideLength, 0, 0}), color1}
	c := Point{b.Coordinate.Add(mgl32.Vec3{0, sideLength, 0}), color1}
	d := Point{a.Coordinate.Add(mgl32.Vec3{0, sideLength, 0}), color1}
	e := Point{d.Coordinate.Add(mgl32.Vec3{0, 0, sideLength}), color2}
	f := Point{e.Coordinate.Add(mgl32.Vec3{sideLength, 0, 0}), color2}
	g := Point{b.Coordinate.Add(mgl32.Vec3{0, 0, sideLength}), color2}
	h := Point{a.Coordinate.Add(mgl32.Vec3{0, 0, sideLength}), color2}
	return &Cube{
		A:   a,
		B:   b,
		C:   c,
		D:   d,
		E:   e,
		F:   f,
		G:   g,
		H:   h,
		vao: NewVAO(),
	}
}

// Log returns the string representation of this object.
func (c *Cube) Log() string {
	logString := "Square:\n"
	logString += " - A : Vector{" + Vec3ToString(c.A.Coordinate) + " }\n"
	logString += " - B : Vector{" + Vec3ToString(c.B.Coordinate) + " }\n"
	logString += " - C : Vector{" + Vec3ToString(c.C.Coordinate) + " }\n"
	logString += " - D : Vector{" + Vec3ToString(c.D.Coordinate) + " }\n"
	logString += " - E : Vector{" + Vec3ToString(c.E.Coordinate) + " }\n"
	logString += " - F : Vector{" + Vec3ToString(c.F.Coordinate) + " }\n"
	logString += " - G : Vector{" + Vec3ToString(c.G.Coordinate) + " }\n"
	logString += " - H : Vector{" + Vec3ToString(c.H.Coordinate) + " }\n"
	return logString
}

// SetShaderProgram updates the shaderProgram of the sphere.
func (c *Cube) SetShaderProgram(p uint32) {
	c.shaderProgram = p
}

func (c *Cube) setupVao() {
	// back
	c.vao.AppendSquarePoints(c.A, c.B, c.C, c.D)
	// right
	c.vao.AppendSquarePoints(c.B, c.G, c.F, c.C)
	// top
	c.vao.AppendSquarePoints(c.C, c.F, c.E, c.D)
	// front
	c.vao.AppendSquarePoints(c.G, c.F, c.E, c.H)
	// left
	c.vao.AppendSquarePoints(c.E, c.D, c.A, c.H)
	// bottom
	c.vao.AppendSquarePoints(c.A, c.B, c.G, c.H)
}
func (c *Cube) setupVaoWithColor() {
	// back - 1,0,0
	red := mgl32.Vec3{1.0, 0.0, 0.0}
	c.A.Color = red
	c.B.Color = red
	c.C.Color = red
	c.D.Color = red
	c.vao.AppendSquarePoints(c.A, c.B, c.C, c.D)
	// right - 0,1,0
	green := mgl32.Vec3{0.0, 1.0, 0.0}
	c.B.Color = green
	c.G.Color = green
	c.F.Color = green
	c.C.Color = green
	c.vao.AppendSquarePoints(c.B, c.G, c.F, c.C)
	// top - 0,0,1
	blue := mgl32.Vec3{0.0, 0.0, 1.0}
	c.C.Color = blue
	c.F.Color = blue
	c.E.Color = blue
	c.D.Color = blue
	c.vao.AppendSquarePoints(c.C, c.F, c.E, c.D)
	// front 0,1,1
	redI := mgl32.Vec3{0.0, 1.0, 1.0}
	c.G.Color = redI
	c.F.Color = redI
	c.E.Color = redI
	c.H.Color = redI
	c.vao.AppendSquarePoints(c.G, c.F, c.E, c.H)
	// left 1,0,1
	greenI := mgl32.Vec3{1.0, 0.0, 1.0}
	c.E.Color = greenI
	c.D.Color = greenI
	c.A.Color = greenI
	c.H.Color = greenI
	c.vao.AppendSquarePoints(c.E, c.D, c.A, c.H)
	// bottom 1,1,0
	blueI := mgl32.Vec3{1.0, 1.0, 0.0}
	c.A.Color = blueI
	c.B.Color = blueI
	c.G.Color = blueI
	c.H.Color = blueI
	c.vao.AppendSquarePoints(c.A, c.B, c.G, c.H)
}

func (c *Cube) buildVao() uint32 {
	c.vao.Clear()
	c.setupVaoWithColor()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(c.vao.Get()), gl.Ptr(c.vao.Get()), gl.STATIC_DRAW)

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

	return vertexArrayObject
}

func (c *Cube) Draw() {
	gl.UseProgram(c.shaderProgram)
	c.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
func (c *Cube) DrawWithUniforms(view, projection mgl32.Mat4) {
	gl.UseProgram(c.shaderProgram)

	viewLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
	projectionLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	modelLocation := gl.GetUniformLocation(c.shaderProgram, gl.Str("model\x00"))

	M := mgl32.Translate3D(c.H.Coordinate.X(), c.H.Coordinate.Y(), c.H.Coordinate.Z())
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])

	c.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
func (c *Cube) Update(delta float64) {
}
