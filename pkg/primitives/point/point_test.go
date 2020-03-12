package point

import (
	"testing"

	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

func TestSetColor(t *testing.T) {
	testData := []struct {
		Color vec.Vector
	}{
		{vec.Vector{0, 0, 0}},
		{vec.Vector{1, 0, 0}},
		{vec.Vector{0, 1, 0}},
		{vec.Vector{0, 1, 0}},
		{vec.Vector{1, 1, 1}},
	}

	for _, tt := range testData {
		var point Point
		point.SetColor(tt.Color)
		if point.Color.X != tt.Color.X || point.Color.Y != tt.Color.Y || point.Color.Z != tt.Color.Z {
			t.Error("Invalid color - SetColor")
		}
	}
}
