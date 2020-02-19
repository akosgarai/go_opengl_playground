package primitives

type Point struct {
	Coordinate Vector
	Color      Vector
}

func (p *Point) SetColor(color Vector) {
	p.Color = color
}
