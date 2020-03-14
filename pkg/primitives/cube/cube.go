package cube

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	P "github.com/akosgarai/opengl_playground/pkg/primitives/point"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

type Cube struct {
	A P.Point // bottom-left-back
	B P.Point // bottom-right-back
	C P.Point // top-right-back
	D P.Point // top-left-back
	E P.Point // top-left-front
	F P.Point // top-right-front
	G P.Point // bottom-right-front
	H P.Point // bottom-left-front
}

func NewCubeByPoints(a, b, c, d, e, f, g, h P.Point) *Cube {
	return &Cube{a, b, c, d, e, f, g, h}
}
func NewCubeByVectorAndLength(in vec.Vector, sideLength float64) *Cube {
	color1 := vec.Vector{1, 0, 0}
	color2 := vec.Vector{0, 1, 0}
	a := P.Point{in, color1}
	b := P.Point{a.Coordinate.Add(vec.Vector{sideLength, 0, 0}), color1}
	c := P.Point{b.Coordinate.Add(vec.Vector{0, sideLength, 0}), color1}
	d := P.Point{a.Coordinate.Add(vec.Vector{0, sideLength, 0}), color1}
	e := P.Point{d.Coordinate.Add(vec.Vector{0, 0, sideLength}), color2}
	f := P.Point{e.Coordinate.Add(vec.Vector{sideLength, 0, 0}), color2}
	g := P.Point{b.Coordinate.Add(vec.Vector{0, 0, sideLength}), color2}
	h := P.Point{a.Coordinate.Add(vec.Vector{0, 0, sideLength}), color2}
	return &Cube{a, b, c, d, e, f, g, h}
}

func (c *Cube) appendPointToVao(currentVao []float32, p P.Point) []float32 {
	currentVao = append(currentVao, float32(p.Coordinate.X))
	currentVao = append(currentVao, float32(p.Coordinate.Y))
	currentVao = append(currentVao, float32(p.Coordinate.Z))
	currentVao = append(currentVao, float32(p.Color.X))
	currentVao = append(currentVao, float32(p.Color.Y))
	currentVao = append(currentVao, float32(p.Color.Z))
	return currentVao
}
func (c *Cube) sideByPointToVao(currentVao []float32, pa, pb, pc, pd P.Point) []float32 {
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
	red := vec.Vector{1.0, 0.0, 0.0}
	c.A.Color = red
	c.B.Color = red
	c.C.Color = red
	c.D.Color = red
	points = c.sideByPointToVao(points, c.A, c.B, c.C, c.D)
	// right - 0,1,0
	green := vec.Vector{0.0, 1.0, 0.0}
	c.B.Color = green
	c.G.Color = green
	c.F.Color = green
	c.C.Color = green
	points = c.sideByPointToVao(points, c.B, c.G, c.F, c.C)
	// top - 0,0,1
	blue := vec.Vector{0.0, 0.0, 1.0}
	c.C.Color = blue
	c.F.Color = blue
	c.E.Color = blue
	c.D.Color = blue
	points = c.sideByPointToVao(points, c.C, c.F, c.E, c.D)
	// front 0,1,1
	redI := vec.Vector{0.0, 1.0, 1.0}
	c.G.Color = redI
	c.F.Color = redI
	c.E.Color = redI
	c.H.Color = redI
	points = c.sideByPointToVao(points, c.G, c.F, c.E, c.H)
	// left 1,0,1
	greenI := vec.Vector{1.0, 0.0, 1.0}
	c.E.Color = greenI
	c.D.Color = greenI
	c.A.Color = greenI
	c.H.Color = greenI
	points = c.sideByPointToVao(points, c.E, c.D, c.A, c.H)
	// bottom 1,1,0
	blueI := vec.Vector{1.0, 1.0, 0.0}
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
