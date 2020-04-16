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

	vao *VAO
}

func NewCubeByPoints(a, b, c, d, e, f, g, h Point) *Cube {
	return &Cube{a, b, c, d, e, f, g, h, NewVAO()}
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
	return &Cube{a, b, c, d, e, f, g, h, NewVAO()}
}

func (c *Cube) sideByPointToVao(pa, pb, pc, pd Point) {
	c.vao.AppendPoint(pa)
	c.vao.AppendPoint(pb)
	c.vao.AppendPoint(pc)
	c.vao.AppendPoint(pa)
	c.vao.AppendPoint(pc)
	c.vao.AppendPoint(pd)
}
func (c *Cube) setupVao() {
	// back
	c.sideByPointToVao(c.A, c.B, c.C, c.D)
	// right
	c.sideByPointToVao(c.B, c.G, c.F, c.C)
	// top
	c.sideByPointToVao(c.C, c.F, c.E, c.D)
	// front
	c.sideByPointToVao(c.G, c.F, c.E, c.H)
	// left
	c.sideByPointToVao(c.E, c.D, c.A, c.H)
	// bottom
	c.sideByPointToVao(c.A, c.B, c.G, c.H)
}
func (c *Cube) setupVaoWithColor() {
	// back - 1,0,0
	red := mgl32.Vec3{1.0, 0.0, 0.0}
	c.A.Color = red
	c.B.Color = red
	c.C.Color = red
	c.D.Color = red
	c.sideByPointToVao(c.A, c.B, c.C, c.D)
	// right - 0,1,0
	green := mgl32.Vec3{0.0, 1.0, 0.0}
	c.B.Color = green
	c.G.Color = green
	c.F.Color = green
	c.C.Color = green
	c.sideByPointToVao(c.B, c.G, c.F, c.C)
	// top - 0,0,1
	blue := mgl32.Vec3{0.0, 0.0, 1.0}
	c.C.Color = blue
	c.F.Color = blue
	c.E.Color = blue
	c.D.Color = blue
	c.sideByPointToVao(c.C, c.F, c.E, c.D)
	// front 0,1,1
	redI := mgl32.Vec3{0.0, 1.0, 1.0}
	c.G.Color = redI
	c.F.Color = redI
	c.E.Color = redI
	c.H.Color = redI
	c.sideByPointToVao(c.G, c.F, c.E, c.H)
	// left 1,0,1
	greenI := mgl32.Vec3{1.0, 0.0, 1.0}
	c.E.Color = greenI
	c.D.Color = greenI
	c.A.Color = greenI
	c.H.Color = greenI
	c.sideByPointToVao(c.E, c.D, c.A, c.H)
	// bottom 1,1,0
	blueI := mgl32.Vec3{1.0, 1.0, 0.0}
	c.A.Color = blueI
	c.B.Color = blueI
	c.G.Color = blueI
	c.H.Color = blueI
	c.sideByPointToVao(c.A, c.B, c.G, c.H)
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
	c.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
