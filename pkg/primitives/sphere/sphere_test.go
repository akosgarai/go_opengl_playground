package sphere

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultColor = []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
)

func TestNew(t *testing.T) {
	sphere := New(2)
	if len(sphere.Points) != 9 {
		t.Error("Invalid number of points")
	}
	if len(sphere.TexCoords) != 9 {
		t.Error("Invalid number of tex coords")
	}
	if len(sphere.Indicies) != 12 {
		t.Error("Invalid number of indicies")
	}
}
func TestMaterialMeshInput(t *testing.T) {
	sphere := New(2)
	vert, ind := sphere.MaterialMeshInput()
	if len(vert) != 9 {
		t.Error("Invalid number of verticies")
	}
	if len(ind) != 12 {
		t.Error("Invalid number of indicies")
	}
}
func TestColoredMeshInput(t *testing.T) {
	sphere := New(2)
	vert, ind := sphere.ColoredMeshInput(DefaultColor)
	if len(vert) != 9 {
		t.Error("Invalid number of verticies")
	}
	if len(ind) != 12 {
		t.Error("Invalid number of indicies")
	}
}
func TestTexturedMeshInput(t *testing.T) {
	sphere := New(2)
	vert, ind := sphere.TexturedMeshInput()
	if len(vert) != 9 {
		t.Error("Invalid number of verticies")
	}
	if len(ind) != 12 {
		t.Error("Invalid number of indicies")
	}
}
