package primitives

import (
	"math"
	"testing"
)

func TestGetMatrix(t *testing.T) {
	testData := []struct {
		Input    float32
		Expected float32
	}{
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, tt := range testData {
		var m Matrix4x4
		for i := 0; i < 16; i++ {
			m.Points[i] = tt.Input
		}
		output := m.GetMatrix()
		for i := 0; i < 16; i++ {
			if output[i] != tt.Expected {
				t.Error("Value mismatch - GetMatrix")
			}
		}
	}
}
func TestGetTransposeMatrix(t *testing.T) {
	testData := []struct {
		Input    float32
		Expected float32
	}{
		{1, 1},
		{2, 2},
		{3, 3},
	}
	for _, tt := range testData {
		var m Matrix4x4
		for i := 0; i < 16; i++ {
			m.Points[i] = tt.Input
		}
		output := m.GetTransposeMatrix()
		for i := 0; i < 16; i++ {
			if output[i] != tt.Expected {
				t.Error("Value mismatch - GetTransposeMatrix")
			}
		}
	}
}
func TestTransposeMatrix(t *testing.T) {
	testData := []struct {
		Input    Matrix4x4
		Expected Matrix4x4
	}{
		{
			Matrix4x4{
				[16]float32{
					1.0, 2.0, 3.0, 4.0,
					5.0, 6.0, 7.0, 8.0,
					9.0, 10.0, 11.0, 12.0,
					13.0, 14.0, 15.0, 16.0},
			},
			Matrix4x4{
				[16]float32{
					1.0, 5.0, 9.0, 13.0,
					2.0, 6.0, 10.0, 14.0,
					3.0, 7.0, 11.0, 15.0,
					4.0, 8.0, 12.0, 16.0},
			},
		},
	}
	for _, tt := range testData {
		transposed := tt.Input.TransposeMatrix()
		for i := 0; i < 16; i++ {
			if transposed.Points[i] != tt.Expected.Points[i] {
				t.Log(transposed)
				t.Log(tt.Expected)
				t.Error("Value mismatch - TransposeMatrix")
			}
		}
	}
}
func TestTranslationMatrix4x4(t *testing.T) {
	testData := []struct {
		Translation Vector
		Expected    Matrix4x4
	}{
		{
			Vector{0, 0, 0},
			*UnitMatrix4x4(),
		},
		{
			Vector{1, 2, 3},
			Matrix4x4{
				[16]float32{
					1.0, 0.0, 0.0, 1.0,
					0.0, 1.0, 0.0, 2.0,
					0.0, 0.0, 1.0, 3.0,
					0.0, 0.0, 0.0, 1.0},
			},
		},
	}
	for _, tt := range testData {
		translated := TranslationMatrix4x4(float32(tt.Translation.X), float32(tt.Translation.Y), float32(tt.Translation.Z))
		for i := 0; i < 16; i++ {
			if translated.Points[i] != tt.Expected.Points[i] {
				t.Log(translated)
				t.Log(tt.Expected)
				t.Error("Value mismatch - TranslationMatrix4x4")
			}
		}
	}
}
func Test_Dot(t *testing.T) {
	testData := []struct {
		m1Input  float32
		m2Input  float32
		Expected float32
	}{
		{0, 0, 0},
	}
	for _, tt := range testData {
		var m1 Matrix4x4
		var m2 Matrix4x4
		for i := 0; i < 16; i++ {
			m1.Points[i] = tt.m1Input
			m2.Points[i] = tt.m2Input
		}
		result := m1.Dot(&m2)
		for i := 0; i < 16; i++ {
			if result.Points[i] != tt.Expected {
				t.Error("Invalid value - Dot")
			}
		}
	}
}
func TestNullMatrix4x4(t *testing.T) {
	nullMatrix := NullMatrix4x4()
	for i := 0; i < 16; i++ {
		if nullMatrix.Points[i] != 0.0 {
			t.Error("Invalid value - NullMatrix4x4")
		}
	}

}
func TestUnitMatrix4x4(t *testing.T) {
	unitMatrix := UnitMatrix4x4()
	for i := 0; i < 16; i++ {
		expected := 0.0
		if i == 0 || i == 5 || i == 10 || i == 15 {
			expected = 1.0
		}
		if unitMatrix.Points[i] != float32(expected) {
			t.Error("Invalid value - UnitMatrix4x4")
		}
	}

}
func TestMultiVector(t *testing.T) {
	testData := []struct {
		M        *Matrix4x4
		V        Vector
		Expected Vector
	}{
		{
			UnitMatrix4x4(),
			Vector{0.0, 0.0, 0.0},
			Vector{0.0, 0.0, 0.0},
		},
	}
	for _, tt := range testData {
		result := tt.M.MultiVector(tt.V)
		if result.X != tt.Expected.X || result.Y != tt.Expected.Y || result.Z != tt.Expected.Z {
			t.Log(result)
			t.Log(tt.Expected)
			t.Error("Invalid result - MultiVector")
		}
	}
}
func TestDotVsMul4Transformation(t *testing.T) {
	testData := []struct {
		M1 *Matrix4x4
		M2 *Matrix4x4
	}{
		{
			UnitMatrix4x4(),
			UnitMatrix4x4(),
		},
		{
			UnitMatrix4x4(),
			NullMatrix4x4(),
		},
		{
			NullMatrix4x4(),
			UnitMatrix4x4(),
		},
		{
			UnitMatrix4x4(),
			&Matrix4x4{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
		},
		{
			&Matrix4x4{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			UnitMatrix4x4(),
		},
		{
			&Matrix4x4{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			&Matrix4x4{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
		},
		{
			&Matrix4x4{[16]float32{0.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			&Matrix4x4{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
		},
	}
	for _, tt := range testData {
		dotProduct := tt.M1.Dot(tt.M2)
		mulProduct := tt.M1.Mul4(tt.M2)
		for i := 0; i < 16; i++ {
			if dotProduct.Points[i] != mulProduct.Points[i] {
				t.Log(tt.M1)
				t.Log(tt.M2)
				t.Log(dotProduct)
				t.Log(mulProduct)
				t.Error("Dot vs Mul4 difference", i)
			}
		}
	}
}
func TestRotationTransformationXAxis(t *testing.T) {
	testData := []struct {
		Point         *Vector
		RotationAngle float64
		RotatedPoint  *Vector
	}{
		{&Vector{0, 0, 0}, 0.0, &Vector{0, 0, 0}},
		{&Vector{1, 0, 0}, 0.0, &Vector{1, 0, 0}},
		{&Vector{1, 1, 0}, 0.0, &Vector{1, 1, 0}},
		{&Vector{1, 1, 1}, 0.0, &Vector{1, 1, 1}},
		{&Vector{0, 0, 0}, 90.0, &Vector{0, 0, 0}},
		{&Vector{1, 0, 0}, 90.0, &Vector{1, 0, 0}},
		{&Vector{0, 1, 0}, math.Pi / 2, &Vector{0, 0, 1}},
		{&Vector{0, 0, 1}, math.Pi / 2, &Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationXMatrix4x4(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != 1.0 ||
			points[1] != 0.0 ||
			points[2] != 0.0 ||
			points[3] != 0.0 ||
			points[4] != 0.0 ||
			points[5] != cosAngle ||
			points[6] != sinAngle ||
			points[7] != 0.0 ||
			points[8] != 0.0 ||
			points[9] != -sinAngle ||
			points[10] != cosAngle ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationX matrix", points, sinAngle, cosAngle)
		}
	}
}
func TestRotationTransformationYAxis(t *testing.T) {
	testData := []struct {
		Point         *Vector
		RotationAngle float64
		RotatedPoint  *Vector
	}{
		{&Vector{0, 0, 0}, 0.0, &Vector{0, 0, 0}},
		{&Vector{1, 0, 0}, 0.0, &Vector{1, 0, 0}},
		{&Vector{1, 1, 0}, 0.0, &Vector{1, 1, 0}},
		{&Vector{1, 1, 1}, 0.0, &Vector{1, 1, 1}},
		{&Vector{0, 0, 0}, 90.0, &Vector{0, 0, 0}},
		{&Vector{0, 1, 0}, 90.0, &Vector{0, 1, 0}},
		{&Vector{1, 0, 0}, math.Pi / 2, &Vector{0, 0, 1}},
		{&Vector{0, 0, 1}, math.Pi / 2, &Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationYMatrix4x4(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != cosAngle ||
			points[1] != 0.0 ||
			points[2] != -sinAngle ||
			points[3] != 0.0 ||
			points[4] != 0.0 ||
			points[5] != 1.0 ||
			points[6] != 0.0 ||
			points[7] != 0.0 ||
			points[8] != sinAngle ||
			points[9] != 0.0 ||
			points[10] != cosAngle ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationY matrix", points, sinAngle, cosAngle)
		}
	}
}
func TestRotationTransformationZAxis(t *testing.T) {
	testData := []struct {
		Point         *Vector
		RotationAngle float64
		RotatedPoint  *Vector
	}{
		{&Vector{0, 0, 0}, 0.0, &Vector{0, 0, 0}},
		{&Vector{1, 0, 0}, 0.0, &Vector{1, 0, 0}},
		{&Vector{1, 1, 0}, 0.0, &Vector{1, 1, 0}},
		{&Vector{1, 1, 1}, 0.0, &Vector{1, 1, 1}},
		{&Vector{0, 0, 0}, 90.0, &Vector{0, 0, 0}},
		{&Vector{0, 0, 1}, 90.0, &Vector{0, 0, 1}},
		{&Vector{1, 0, 0}, math.Pi / 2, &Vector{0, 1, 0}},
		{&Vector{0, 1, 0}, math.Pi / 2, &Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationZMatrix4x4(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != cosAngle ||
			points[1] != sinAngle ||
			points[2] != 0.0 ||
			points[3] != 0.0 ||
			points[4] != -sinAngle ||
			points[5] != cosAngle ||
			points[6] != 0.0 ||
			points[7] != 0.0 ||
			points[8] != 0.0 ||
			points[9] != 0.0 ||
			points[10] != 1.0 ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationY matrix", points, sinAngle, cosAngle)
		}
	}
}
