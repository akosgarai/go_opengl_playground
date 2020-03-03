package primitives

import (
	"testing"
)

func TestSetColor(t *testing.T) {
	testData := []struct {
		Color Vector
	}{
		{Vector{0, 0, 0}},
		{Vector{1, 0, 0}},
		{Vector{0, 1, 0}},
		{Vector{0, 1, 0}},
		{Vector{1, 1, 1}},
	}

	for _, tt := range testData {
		var point Point
		point.SetColor(tt.Color)
		if point.Color.X != tt.Color.X || point.Color.Y != tt.Color.Y || point.Color.Z != tt.Color.Z {
			t.Error("Invalid color - SetColor")
		}
	}
}
