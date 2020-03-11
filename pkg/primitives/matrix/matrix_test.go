package matrix

import (
	"testing"

	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
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
		var m Matrix
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
		var m Matrix
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
		Input    Matrix
		Expected Matrix
	}{
		{
			Matrix{
				[16]float32{
					1.0, 2.0, 3.0, 4.0,
					5.0, 6.0, 7.0, 8.0,
					9.0, 10.0, 11.0, 12.0,
					13.0, 14.0, 15.0, 16.0},
			},
			Matrix{
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
func TestDot(t *testing.T) {
	testData := []struct {
		m1Input  float32
		m2Input  float32
		Expected float32
	}{
		{0, 0, 0},
	}
	for _, tt := range testData {
		var m1 Matrix
		var m2 Matrix
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
func TestNullMatrix(t *testing.T) {
	nullMatrix := NullMatrix()
	for i := 0; i < 16; i++ {
		if nullMatrix.Points[i] != 0.0 {
			t.Error("Invalid value - NullMatrix")
		}
	}

}
func TestUnitMatrix(t *testing.T) {
	unitMatrix := UnitMatrix()
	for i := 0; i < 16; i++ {
		expected := 0.0
		if i == 0 || i == 5 || i == 10 || i == 15 {
			expected = 1.0
		}
		if unitMatrix.Points[i] != float32(expected) {
			t.Error("Invalid value - UnitMatrix")
		}
	}

}
func TestMultiVector(t *testing.T) {
	testData := []struct {
		M        *Matrix
		V        vec.Vector
		Expected vec.Vector
	}{
		{
			UnitMatrix(),
			vec.Vector{0.0, 0.0, 0.0},
			vec.Vector{0.0, 0.0, 0.0},
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
		M1 *Matrix
		M2 *Matrix
	}{
		{
			UnitMatrix(),
			UnitMatrix(),
		},
		{
			UnitMatrix(),
			NullMatrix(),
		},
		{
			NullMatrix(),
			UnitMatrix(),
		},
		{
			UnitMatrix(),
			&Matrix{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
		},
		{
			&Matrix{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			UnitMatrix(),
		},
		{
			&Matrix{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			&Matrix{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
		},
		{
			&Matrix{[16]float32{0.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
			&Matrix{[16]float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0}},
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
