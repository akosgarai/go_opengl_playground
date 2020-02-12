package primitives

type Matrix4x4 struct {
	Points [16]float32
}

func NullMatrix4x4() *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
		},
	}
}
func UnitMatrix4x4() *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, 1.0, 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

func (m *Matrix4x4) Dot(m2 *Matrix4x4) *Matrix4x4 {
	result := NullMatrix4x4()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result.Points[4*i+j] = m.Points[4*i+0]*m2.Points[4*0+j] +
				m.Points[4*i+1]*m2.Points[4*1+j] +
				m.Points[4*i+2]*m2.Points[4*2+j] +
				m.Points[4*i+3]*m2.Points[4*3+j]
		}
	}
	return result
}
