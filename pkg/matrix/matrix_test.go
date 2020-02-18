package matrix

import (
	"testing"

	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

func TestPoints(t *testing.T) {
	testData := []struct {
		Input    float64
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
		output := m.GetPoints()
		for i := 0; i < 16; i++ {
			if output[i] != tt.Expected {
				t.Error("Value mismatch - Points")
			}
		}
	}
}
func TestClear(t *testing.T) {
	testData := []struct {
		Input    float64
		Expected float64
	}{
		{1, 0},
		{2, 0},
		{3, 0},
	}
	for _, tt := range testData {
		var m Matrix
		for i := 0; i < 16; i++ {
			m.Points[i] = tt.Input
		}
		m.Clear()
		for i := 0; i < 16; i++ {
			if m.Points[i] != tt.Expected {
				t.Error("Value mismatch after clear")
			}
		}
	}
}
func TestLoadIdentity(t *testing.T) {
	testData := []struct {
		Input    float64
		Expected float64
	}{
		{1, 1},
		{2, 1},
		{3, 1},
	}
	for _, tt := range testData {
		var m Matrix
		for i := 0; i < 16; i++ {
			m.Points[i] = tt.Input
		}
		m.LoadIdentity()
		for i := 0; i < 16; i++ {
			if i%4 == 0 {
				if m.Points[i] != tt.Expected {
					t.Error("Value mismatch after clear")
				}
			} else {
				if m.Points[i] != 0 {
					t.Error("Value mismatch after clear")
				}
			}
		}
	}
}
func TestAdd(t *testing.T) {
	testData := []struct {
		m1Input  float64
		m2Input  float64
		Expected float64
	}{
		{1, 0, 1},
		{1, 1, 2},
		{2, 1, 3},
		{3, 1, 4},
		{3, 4, 7},
	}
	for _, tt := range testData {
		var m1 Matrix
		var m2 Matrix
		for i := 0; i < 16; i++ {
			m1.Points[i] = tt.m1Input
			m2.Points[i] = tt.m2Input
		}
		result1 := m1.Add(m2)
		result2 := m2.Add(m1)
		for i := 0; i < 16; i++ {
			if result1.Points[i] != tt.Expected || result2.Points[i] != tt.Expected {
				t.Error("Invalid value - Add")
			}
		}
	}
}
func TestDot(t *testing.T) {
	testData := []struct {
		m1Input  float64
		m2Input  float64
		Expected float64
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
		result := m1.Dot(m2)
		for i := 0; i < 16; i++ {
			if result.Points[i] != tt.Expected {
				t.Error("Invalid value - Dot")
			}
		}
	}
}
func TestMultiVector(t *testing.T) {
	var m Matrix
	v := V.Vector{1, 1, 1}
	m.LoadIdentity()
	result := m.MultiVector(v)
	if result.X != 1 || result.Y != 1 || result.Z != 1 {
		t.Error("Unexpected value - MultiVector")
	}
}
func TestTranslation(t *testing.T) {
	testData := []struct {
		v V.Vector
	}{
		{V.Vector{1, 1, 1}},
		{V.Vector{4, 2, 0}},
		{V.Vector{0, 0, 0}},
	}
	for _, tt := range testData {
		v := tt.v
		result := Translation(v)
		if result.Points[12] != v.X || result.Points[13] != v.Y || result.Points[14] != v.Z {
			t.Error("Unexpected value - Translation")
		}
	}
}
