package primitives

import (
	"math"
)

type Matrix4x4 struct {
	Points [16]float32
}

// GetMatrix returns the points of the matrix
func (m *Matrix4x4) GetMatrix() [16]float32 {
	return m.Points
}

// GetMatrix returns the points of the transpose matrix
func (m *Matrix4x4) GetTransposeMatrix() [16]float32 {
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
func (m *Matrix4x4) TransposeMatrix() *Matrix4x4 {
	return &Matrix4x4{m.GetTransposeMatrix()}
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
func ScaleMatrix4x4(scaleX, scaleY, scaleZ float32) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			scaleX, 0.0, 0.0, 0.0,
			0.0, scaleY, 0.0, 0.0,
			0.0, 0.0, scaleZ, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func TranslationMatrix4x4(translationX, translationY, translationZ float32) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, translationX,
			0.0, 1.0, 0.0, translationY,
			0.0, 0.0, 1.0, translationZ,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func RotationXMatrix4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, float32(math.Cos(rotationAngle)), float32(-math.Sin(rotationAngle)), 0.0,
			0.0, float32(math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func RotationYMatrix4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			float32(math.Cos(rotationAngle)), 0.0, float32(math.Sin(rotationAngle)), 0.0,
			0.0, 1.0, 0.0, 0.0,
			float32(-math.Sin(rotationAngle)), 0.0, float32(math.Cos(rotationAngle)), 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func RotationZMatrix4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			float32(math.Cos(rotationAngle)), float32(-math.Sin(rotationAngle)), 0.0, 0.0,
			float32(math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

func ProjectionMatrix4x4(angleOfView, near, far float64) *Matrix4x4 {
	scale := float32(1 / math.Tan(angleOfView*0.5*math.Pi/180))
	projection := UnitMatrix4x4()
	projection.Points[0] = scale
	projection.Points[5] = scale
	projection.Points[10] = float32(-far / (far - near))
	projection.Points[14] = float32(-far * near / (far - near))
	projection.Points[11] = -1
	projection.Points[15] = 0
	return projection
}

func (m *Matrix4x4) Dot(m2 *Matrix4x4) *Matrix4x4 {
	result := NullMatrix4x4()
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
func (m *Matrix4x4) MultiVector(v Vector) *Vector {
	Xh := float64(m.Points[0])*v.X + float64(m.Points[1])*v.Y + float64(m.Points[2])*v.Z + float64(m.Points[3])
	Yh := float64(m.Points[4])*v.X + float64(m.Points[5])*v.Y + float64(m.Points[6])*v.Z + float64(m.Points[7])
	Zh := float64(m.Points[8])*v.X + float64(m.Points[9])*v.Y + float64(m.Points[10])*v.Z + float64(m.Points[11])
	h := float64(m.Points[12])*v.X + float64(m.Points[13])*v.Y + float64(m.Points[14])*v.Z + float64(m.Points[15])
	return &Vector{Xh / h, Yh / h, Zh / h}
}

// Solution from here : https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/matrix.go#L2143-L2166
func (m1 *Matrix4x4) Mul4(m2 *Matrix4x4) *Matrix4x4 {
	return &Matrix4x4{
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
