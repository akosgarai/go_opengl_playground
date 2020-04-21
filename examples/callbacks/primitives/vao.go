package primitives

type VAO struct {
	vao []float32
}

func NewVAO() *VAO {
	return &VAO{
		vao: []float32{},
	}
}

func (v *VAO) AppendPoint(p Point) {
	v.vao = append(v.vao, p.Coordinate.X())
	v.vao = append(v.vao, p.Coordinate.Y())
	v.vao = append(v.vao, p.Coordinate.Z())
	v.vao = append(v.vao, p.Color.X())
	v.vao = append(v.vao, p.Color.Y())
	v.vao = append(v.vao, p.Color.Z())
}
func (v *VAO) AppendTrianglePoints(pa, pb, pc Point) {
	v.AppendPoint(pa)
	v.AppendPoint(pb)
	v.AppendPoint(pc)
}
func (v *VAO) AppendSquarePoints(pa, pb, pc, pd Point) {
	v.AppendPoint(pa)
	v.AppendPoint(pb)
	v.AppendPoint(pc)
	v.AppendPoint(pa)
	v.AppendPoint(pc)
	v.AppendPoint(pd)
}
func (v *VAO) Get() []float32 {
	return v.vao
}
func (v *VAO) Clear() {
	v.vao = []float32{}
}
