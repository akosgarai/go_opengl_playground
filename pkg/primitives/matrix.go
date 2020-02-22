package primitives

import (
	"math"
)

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
func TranslationMatrixT4x4(translationX, translationY, translationZ float32) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, 1.0, 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			translationX, translationY, translationZ, 1.0,
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
func RotationXMatrixT4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, float32(math.Cos(rotationAngle)), float32(math.Sin(rotationAngle)), 0.0,
			0.0, float32(-math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0,
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
func RotationYMatrixT4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			float32(math.Cos(rotationAngle)), 0.0, float32(-math.Sin(rotationAngle)), 0.0,
			0.0, 1.0, 0.0, 0.0,
			float32(math.Sin(rotationAngle)), 0.0, float32(math.Cos(rotationAngle)), 0.0,
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
func RotationZMatrixT4x4(rotationAngle float64) *Matrix4x4 {
	return &Matrix4x4{
		[16]float32{
			float32(math.Cos(rotationAngle)), float32(math.Sin(rotationAngle)), 0.0, 0.0,
			float32(-math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0, 0.0,
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

// https://stackoverflow.com/questions/8115352/glmperspective-explanation
// ProjectionNewSolution.
func Perspective(angle, ratio, near, far float32) *Matrix4x4 {
	// degree to radian formula: n deg = n * PI / 180 rad
	slopey := float32(math.Tan(float64(angle * math.Pi / 180)))
	result := NullMatrix4x4()
	result.Points[0] = 1 / slopey / ratio
	result.Points[5] = 1 / slopey
	result.Points[10] = -((far + near) / (far - near))
	result.Points[11] = -1
	result.Points[14] = -(2 * far * near / (far - near))
	return result
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

// MultiVector returns a new Vector. this is the multiplication of a vector - matrix element.
func (m *Matrix4x4) MultiVector(v Vector) *Vector {
	Xh := float64(m.Points[0])*v.X + float64(m.Points[1])*v.Y + float64(m.Points[2])*v.Z + float64(m.Points[3])
	Yh := float64(m.Points[4])*v.X + float64(m.Points[5])*v.Y + float64(m.Points[6])*v.Z + float64(m.Points[7])
	Zh := float64(m.Points[8])*v.X + float64(m.Points[9])*v.Y + float64(m.Points[10])*v.Z + float64(m.Points[11])
	h := float64(m.Points[12])*v.X + float64(m.Points[13])*v.Y + float64(m.Points[14])*v.Z + float64(m.Points[15])
	return &Vector{Xh / h, Yh / h, Zh / h}
}
