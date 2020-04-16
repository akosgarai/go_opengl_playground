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
func (v *VAO) Get() []float32 {
	return v.vao
}
func (v *VAO) Clear() {
	v.vao = []float32{}
}
