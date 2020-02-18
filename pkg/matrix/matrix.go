// It represents the 4x4 matrixes
package matrix

import (
	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

type Matrix struct {
	Points [16]float64
}

// GetPoints returns the points in float32 format.
func (m *Matrix) GetPoints() [16]float32 {
	var result [16]float32
	for i := 0; i < 16; i++ {
		result[i] = float32(m.Points[i])
	}
	return result
}

// Clear makes nullmatrix from this matrix.
func (m *Matrix) Clear() {
	for i := 0; i < 16; i++ {
		m.Points[i] = 0
	}
}

// LoadIdentity makes a unit matrix from this matrix.
func (m *Matrix) LoadIdentity() {
	m.Clear()
	m.Points[0] = 1
	m.Points[4] = 1
	m.Points[8] = 1
	m.Points[12] = 1
}

// Add returns a new matrix, it doesn't update this matrix.
func (m *Matrix) Add(m2 Matrix) Matrix {
	var result Matrix
	for i := 0; i < 16; i++ {
		result.Points[i] = m.Points[i] + m2.Points[i]
	}
	return result
}

// Dot returns a new matrix, it doesn't update this matrix. It constructs the multiplication of the given matrixes.
func (m *Matrix) Dot(m2 Matrix) *Matrix {
	var result Matrix
	result.Clear()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result.Points[4*i+j] = m.Points[4*i+0]*m2.Points[4*0+j] +
				m.Points[4*i+1]*m2.Points[4*1+j] +
				m.Points[4*i+2]*m2.Points[4*2+j] +
				m.Points[4*i+3]*m2.Points[4*3+j]
		}
	}
	return &result
}

// MultiVector returns a new Vector. this is the multiplication of a vector - matrix element.
func (m *Matrix) MultiVector(v V.Vector) *V.Vector {
	Xh := m.Points[0]*v.X + m.Points[1]*v.Y + m.Points[2]*v.Z + m.Points[3]
	Yh := m.Points[4]*v.X + m.Points[5]*v.Y + m.Points[6]*v.Z + m.Points[7]
	Zh := m.Points[8]*v.X + m.Points[9]*v.Y + m.Points[10]*v.Z + m.Points[11]
	h := m.Points[12]*v.X + m.Points[13]*v.Y + m.Points[14]*v.Z + m.Points[15]
	return &V.Vector{Xh / h, Yh / h, Zh / h}
}
