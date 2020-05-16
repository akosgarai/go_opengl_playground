package sphere

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultRadius    = float32(2.0)
	DefaultColor     = mgl32.Vec3{0, 0, 1}
	DefaultCenter    = mgl32.Vec3{3, 3, 5}
	DefaultSpeed     = float32(0.0)
	DefaultDirection = mgl32.Vec3{0, 0, 0}
	DefaultAngle     = float32(0.0)
	DefaultAxis      = mgl32.Vec3{0, 0, 0}
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
