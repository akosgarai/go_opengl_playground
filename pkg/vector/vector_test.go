package vector

import (
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
func TestMultiplyVector(t *testing.T) {
	testData := []struct {
		Vector1  Vector
		Vector2  Vector
		Expected Vector
	}{
		{Vector{0, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{1, 0, 0}, Vector{0, 0, 0}, Vector{0, 0, 0}},
		{Vector{0, 1, 0}, Vector{0, 1, 0}, Vector{0, 1, 0}},
		{Vector{2, 1, 3}, Vector{2, 0, 1}, Vector{4, 0, 3}},
	}

	for _, tt := range testData {
		value := tt.Vector1.MultiplyVector(tt.Vector2)
		if value.X != tt.Expected.X || value.Y != tt.Expected.Y || value.Z != tt.Expected.Z {
			t.Error("Invalid output for MultiplyVector")
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
		{Vector{2, 2, 2}, 4},
	}
	for _, tt := range testData {
		value := tt.v.Length()
		if value != tt.Expected {
			t.Error("Invalid output for Length")
		}
	}
}
