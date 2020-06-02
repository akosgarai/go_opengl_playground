package cylinder

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	testData := []struct {
		input                   mgl32.Vec3
		expectedLengthPoint     int
		expectedLengthNormal    int
		expectedLengthIndices   int
		expectedLengthTexCoords int
	}{
		{mgl32.Vec3{1, 2, 1}, 10, 10, 24, 10},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		if len(cyl.Points) != tt.expectedLengthPoint {
			t.Errorf("Invalid point length. instead of '%d', we have '%d'\n", tt.expectedLengthPoint, len(cyl.Points))
		}
		if len(cyl.Normals) != tt.expectedLengthNormal {
			t.Errorf("Invalid normals length. instead of '%d', we have '%d'\n", tt.expectedLengthNormal, len(cyl.Normals))
		}
		if len(cyl.Indices) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(cyl.Indices))
		}
		if len(cyl.TexCoords) != tt.expectedLengthTexCoords {
			t.Errorf("Invalid tex coords length. instead of '%d', we have '%d'\n", tt.expectedLengthTexCoords, len(cyl.TexCoords))
		}
	}
}
func TestNewHalfCircleBased(t *testing.T) {
	testData := []struct {
		input                   mgl32.Vec3
		expectedLengthPoint     int
		expectedLengthNormal    int
		expectedLengthIndices   int
		expectedLengthTexCoords int
	}{
		{mgl32.Vec3{1, 2, 1}, 10, 10, 24, 10},
	}
	for _, tt := range testData {
		cyl := NewHalfCircleBased(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		if len(cyl.Points) != tt.expectedLengthPoint {
			t.Errorf("Invalid point length. instead of '%d', we have '%d'\n", tt.expectedLengthPoint, len(cyl.Points))
		}
		if len(cyl.Normals) != tt.expectedLengthNormal {
			t.Errorf("Invalid normals length. instead of '%d', we have '%d'\n", tt.expectedLengthNormal, len(cyl.Normals))
		}
		if len(cyl.Indices) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(cyl.Indices))
		}
		if len(cyl.TexCoords) != tt.expectedLengthTexCoords {
			t.Errorf("Invalid tex coords length. instead of '%d', we have '%d'\n", tt.expectedLengthTexCoords, len(cyl.TexCoords))
		}
	}
}
func TestMaterialMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 10, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i := cyl.MaterialMeshInput()
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
func TestColoredMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 10, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i := cyl.ColoredMeshInput([]mgl32.Vec3{mgl32.Vec3{1, 1, 1}})
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
func TestTexturedMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 10, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i := cyl.TexturedMeshInput()
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
