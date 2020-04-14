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
}

func NewCubeByPoints(a, b, c, d, e, f, g, h Point) *Cube {
	return &Cube{a, b, c, d, e, f, g, h}
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
	return &Cube{a, b, c, d, e, f, g, h}
}

func (c *Cube) appendPointToVao(currentVao []float32, p Point) []float32 {
	currentVao = append(currentVao, p.Coordinate.X())
	currentVao = append(currentVao, p.Coordinate.Y())
	currentVao = append(currentVao, p.Coordinate.Z())
	currentVao = append(currentVao, p.Color.X())
	currentVao = append(currentVao, p.Color.Y())
	currentVao = append(currentVao, p.Color.Z())
	return currentVao
}
func (c *Cube) sideByPointToVao(currentVao []float32, pa, pb, pc, pd Point) []float32 {
	currentVao = c.appendPointToVao(currentVao, pa)
	currentVao = c.appendPointToVao(currentVao, pb)
	currentVao = c.appendPointToVao(currentVao, pc)
	currentVao = c.appendPointToVao(currentVao, pa)
	currentVao = c.appendPointToVao(currentVao, pc)
	currentVao = c.appendPointToVao(currentVao, pd)
	return currentVao
}
func (c *Cube) setupVao() []float32 {
	var points []float32

	// back
	points = c.sideByPointToVao(points, c.A, c.B, c.C, c.D)
	// right
	points = c.sideByPointToVao(points, c.B, c.G, c.F, c.C)
	// top
	points = c.sideByPointToVao(points, c.C, c.F, c.E, c.D)
	// front
	points = c.sideByPointToVao(points, c.G, c.F, c.E, c.H)
	// left
	points = c.sideByPointToVao(points, c.E, c.D, c.A, c.H)
	// bottom
	points = c.sideByPointToVao(points, c.A, c.B, c.G, c.H)

	return points
}
func (c *Cube) setupVaoWithColor() []float32 {
	var points []float32

	// back - 1,0,0
	red := mgl32.Vec3{1.0, 0.0, 0.0}
	c.A.Color = red
	c.B.Color = red
	c.C.Color = red
	c.D.Color = red
	points = c.sideByPointToVao(points, c.A, c.B, c.C, c.D)
	// right - 0,1,0
	green := mgl32.Vec3{0.0, 1.0, 0.0}
	c.B.Color = green
	c.G.Color = green
	c.F.Color = green
	c.C.Color = green
	points = c.sideByPointToVao(points, c.B, c.G, c.F, c.C)
	// top - 0,0,1
	blue := mgl32.Vec3{0.0, 0.0, 1.0}
	c.C.Color = blue
	c.F.Color = blue
	c.E.Color = blue
	c.D.Color = blue
	points = c.sideByPointToVao(points, c.C, c.F, c.E, c.D)
	// front 0,1,1
	redI := mgl32.Vec3{0.0, 1.0, 1.0}
	c.G.Color = redI
	c.F.Color = redI
	c.E.Color = redI
	c.H.Color = redI
	points = c.sideByPointToVao(points, c.G, c.F, c.E, c.H)
	// left 1,0,1
	greenI := mgl32.Vec3{1.0, 0.0, 1.0}
	c.E.Color = greenI
	c.D.Color = greenI
	c.A.Color = greenI
	c.H.Color = greenI
	points = c.sideByPointToVao(points, c.E, c.D, c.A, c.H)
	// bottom 1,1,0
	blueI := mgl32.Vec3{1.0, 1.0, 0.0}
	c.A.Color = blueI
	c.B.Color = blueI
	c.G.Color = blueI
	c.H.Color = blueI
	points = c.sideByPointToVao(points, c.A, c.B, c.G, c.H)

	return points
}

func (c *Cube) buildVao() uint32 {
	points := c.setupVaoWithColor()

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

	return vertexArrayObject
}

func (c *Cube) Draw() {
	vertexArrayObject := c.buildVao()
	gl.BindVertexArray(vertexArrayObject)
	gl.DrawArrays(gl.TRIANGLES, 0, 3*12)
}
