package vector

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected Vector
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, Vector{1, 0, 0}},
		{Vector{0, 1, 0}, Vector{0, 0, 0}, Vector{0, 1, 0}},
		{Vector{0, 1, 0}, Vector{0, 0, 1}, Vector{0, 1, 1}},
	}

	for _, tt := range testData {
		value := tt.Vector1.Add(tt.Vector2)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for Add")
		}
	}
}
func TestSubtract(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected Vector
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, Vector{1, 0, 0}},
		{Vector{0, 1, 0}, Vector{0, 0, 0}, Vector{0, 1, 0}},
		{Vector{0, 1, 0}, Vector{0, 0, 1}, Vector{0, 1, -1}},
	}

	for _, tt := range testData {
		value := tt.Vector1.Subtract(tt.Vector2)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for Subtract")
		}
	}
}
func TestMultiplyScalar(t *testing.T) {
	testData := []struct {
		v        Vector
		s        float64
		Expected Vector
	}{
		{Vector{0, 0, 0}, 1, Vector{0, 0, 0}},
		{Vector{1, 1, 1}, 0, Vector{0, 0, 0}},
		{Vector{1, 2, 3}, 1, Vector{1, 2, 3}},
		{Vector{1, 2, 3}, 2, Vector{2, 4, 6}},
	}
	for _, tt := range testData {
		value := tt.v.MultiplyScalar(tt.s)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for MultiplyScalar")
		}
	}
}
func TestAddScalar(t *testing.T) {
	testData := []struct {
		v        Vector
		s        float64
		Expected Vector
	}{
		{Vector{0, 0, 0}, 1, Vector{1, 1, 1}},
		{Vector{1, 1, 1}, 0, Vector{1, 1, 1}},
		{Vector{1, 2, 3}, 1, Vector{2, 3, 4}},
		{Vector{1, 2, 3}, 2, Vector{3, 4, 5}},
	}
	for _, tt := range testData {
		value := tt.v.AddScalar(tt.s)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for AddScalar")
		}
	}
}
func TestCross(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected Vector
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{0, 1, 0}, Vector{0, 1, 0}, Vector{0, 0, 0}},
		{Vector{0, 1, 0}, Vector{1, 0, 0}, Vector{0, 0, -1}},
		{Vector{1, 1, 1}, Vector{1, 1, 1}, Vector{0, 0, 0}},
	}

	for _, tt := range testData {
		value := tt.Vector1.Cross(tt.Vector2)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for Cross")
		}
	}
}
func TestLength(t *testing.T) {
	testData := []struct {
		v        Vector
		Expected float64
	}{
		{Vector{0, 0, 0}, 0},
		{Vector{4, 4, 4}, math.Sqrt(3 * 4 * 4)},
	}
	for _, tt := range testData {
		value := tt.v.Length()
		if value != tt.Expected {
			t.Error("Invalid output for Length")
		}
	}
}
func TestNormalize(t *testing.T) {
	testData := []struct {
		v        Vector
		Expected Vector
	}{
		{Vector{2, 0, 0}, Vector{1, 0, 0}},
		{Vector{0, 5, 0}, Vector{0, 1, 0}},
		{Vector{0, 0, 9}, Vector{0, 0, 1}},
	}
	for _, tt := range testData {
		value := tt.v.Normalize()
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Log(value)
			t.Log(tt.Expected)
			t.Error("Invalid output for Normalize")
		}
	}
}
func TestMultiply(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected Vector
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{0, 1, 0}, Vector{0, 1, 0}, Vector{0, 1, 0}},
		{Vector{0, 1, 0}, Vector{1, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 1, 1}, Vector{1, 1, 1}, Vector{1, 1, 1}},
		{Vector{1, 2, 3}, Vector{5, 7, 11}, Vector{5, 14, 33}},
	}

	for _, tt := range testData {
		value := tt.Vector1.Multiply(tt.Vector2)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for Multiply")
		}
	}
}
func TestDot(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected float64
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, 0.0},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, 0.0},
		{Vector{0, 1, 0}, Vector{0, 1, 0}, 1.0},
		{Vector{0, 1, 0}, Vector{1, 0, 0}, 0.0},
		{Vector{1, 1, 1}, Vector{1, 1, 1}, 3.0},
		{Vector{1, 2, 3}, Vector{5, 7, 11}, 52},
	}

	for _, tt := range testData {
		value := tt.Vector1.Dot(tt.Vector2)
		if value != tt.Expected {
			t.Error("Invalid output for Dot")
		}
	}
}
func TestToString(t *testing.T) {
	testData := []struct {
		v        Vector
		Expected string
	}{
		{Vector{2, 0, 0}, "X : 2.000000, Y : 0.000000, Z : 0.000000"},
		{Vector{0, 5, 0}, "X : 0.000000, Y : 5.000000, Z : 0.000000"},
		{Vector{0, 0, 9}, "X : 0.000000, Y : 0.000000, Z : 9.000000"},
		{Vector{1, 2, 3}, "X : 1.000000, Y : 2.000000, Z : 3.000000"},
	}
	for _, tt := range testData {
		value := tt.v.ToString()
		if value != tt.Expected {
			t.Log(value)
			t.Error("Invalid output for ToString")
		}
	}
}
