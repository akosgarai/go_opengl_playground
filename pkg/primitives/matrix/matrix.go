package matrix

import (
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

type Matrix struct {
	Points [16]float32
}

func NullMatrix() *Matrix {
	return &Matrix{
		[16]float32{
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
			0.0, 0.0, 0.0, 0.0,
		},
	}
}
func UnitMatrix() *Matrix {
	return &Matrix{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, 1.0, 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

// GetMatrix returns the points of the matrix
func (m *Matrix) GetMatrix() [16]float32 {
	return m.Points
}

// GetMatrix returns the points of the transpose matrix
func (m *Matrix) GetTransposeMatrix() [16]float32 {
	var result [16]float32
	// i: col, j: row.
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[4*i+j] = m.Points[4*j+i]
		}
	}
	return result
}

// TransposeMatrix returns the transposed matrix
func (m *Matrix) TransposeMatrix() *Matrix {
	return &Matrix{m.GetTransposeMatrix()}
}
func (m *Matrix) Dot(m2 *Matrix) *Matrix {
	result := NullMatrix()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result.Points[4*i+j] = m2.Points[4*i+0]*m.Points[4*0+j] +
				m2.Points[4*i+1]*m.Points[4*1+j] +
				m2.Points[4*i+2]*m.Points[4*2+j] +
				m2.Points[4*i+3]*m.Points[4*3+j]
		}
	}
	return result
}

// MultiVector returns a new Vector. this is the multiplication of a vector - matrix element.
func (m *Matrix) MultiVector(v vec.Vector) *vec.Vector {
	Xh := float64(m.Points[0])*v.X + float64(m.Points[4])*v.Y + float64(m.Points[8])*v.Z + float64(m.Points[12])
	Yh := float64(m.Points[1])*v.X + float64(m.Points[5])*v.Y + float64(m.Points[9])*v.Z + float64(m.Points[13])
	Zh := float64(m.Points[2])*v.X + float64(m.Points[6])*v.Y + float64(m.Points[10])*v.Z + float64(m.Points[14])
	h := float64(m.Points[3])*v.X + float64(m.Points[7])*v.Y + float64(m.Points[11])*v.Z + float64(m.Points[15])
	return &vec.Vector{Xh / h, Yh / h, Zh / h}
}

// Solution from here : https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/matrix.go#L2143-L2166
func (m1 *Matrix) Mul4(m2 *Matrix) *Matrix {
	return &Matrix{
		[16]float32{
			m1.Points[0]*m2.Points[0] + m1.Points[4]*m2.Points[1] + m1.Points[8]*m2.Points[2] + m1.Points[12]*m2.Points[3],
			m1.Points[1]*m2.Points[0] + m1.Points[5]*m2.Points[1] + m1.Points[9]*m2.Points[2] + m1.Points[13]*m2.Points[3],
			m1.Points[2]*m2.Points[0] + m1.Points[6]*m2.Points[1] + m1.Points[10]*m2.Points[2] + m1.Points[14]*m2.Points[3],
			m1.Points[3]*m2.Points[0] + m1.Points[7]*m2.Points[1] + m1.Points[11]*m2.Points[2] + m1.Points[15]*m2.Points[3],
			m1.Points[0]*m2.Points[4] + m1.Points[4]*m2.Points[5] + m1.Points[8]*m2.Points[6] + m1.Points[12]*m2.Points[7],
			m1.Points[1]*m2.Points[4] + m1.Points[5]*m2.Points[5] + m1.Points[9]*m2.Points[6] + m1.Points[13]*m2.Points[7],
			m1.Points[2]*m2.Points[4] + m1.Points[6]*m2.Points[5] + m1.Points[10]*m2.Points[6] + m1.Points[14]*m2.Points[7],
			m1.Points[3]*m2.Points[4] + m1.Points[7]*m2.Points[5] + m1.Points[11]*m2.Points[6] + m1.Points[15]*m2.Points[7],
			m1.Points[0]*m2.Points[8] + m1.Points[4]*m2.Points[9] + m1.Points[8]*m2.Points[10] + m1.Points[12]*m2.Points[11],
			m1.Points[1]*m2.Points[8] + m1.Points[5]*m2.Points[9] + m1.Points[9]*m2.Points[10] + m1.Points[13]*m2.Points[11],
			m1.Points[2]*m2.Points[8] + m1.Points[6]*m2.Points[9] + m1.Points[10]*m2.Points[10] + m1.Points[14]*m2.Points[11],
			m1.Points[3]*m2.Points[8] + m1.Points[7]*m2.Points[9] + m1.Points[11]*m2.Points[10] + m1.Points[15]*m2.Points[11],
			m1.Points[0]*m2.Points[12] + m1.Points[4]*m2.Points[13] + m1.Points[8]*m2.Points[14] + m1.Points[12]*m2.Points[15],
			m1.Points[1]*m2.Points[12] + m1.Points[5]*m2.Points[13] + m1.Points[9]*m2.Points[14] + m1.Points[13]*m2.Points[15],
			m1.Points[2]*m2.Points[12] + m1.Points[6]*m2.Points[13] + m1.Points[10]*m2.Points[14] + m1.Points[14]*m2.Points[15],
			m1.Points[3]*m2.Points[12] + m1.Points[7]*m2.Points[13] + m1.Points[11]*m2.Points[14] + m1.Points[15]*m2.Points[15],
		},
	}
}
