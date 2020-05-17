package cuboid

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	Indicies = []uint32{
		0, 1, 2, 0, 2, 3, // bottom
		4, 5, 6, 4, 6, 7, // top
		8, 9, 10, 8, 10, 11, // front
		12, 13, 14, 12, 14, 15, // back
		16, 17, 18, 16, 18, 19, // left
		20, 21, 22, 20, 22, 23, // right
	}
	Color = []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
	}
)

func TestNew(t *testing.T) {
	testData := []struct {
		sideLengths         [3]float32
		expectedCoordinates [3]float32
	}{
		{[3]float32{1, 1, 1}, [3]float32{0.5, 0.5, 0.5}},
		{[3]float32{2, 1, 1}, [3]float32{0.5, 0.25, 0.25}},
		{[3]float32{1, 2, 1}, [3]float32{0.25, 0.25, 0.5}},
		{[3]float32{1, 1, 2}, [3]float32{0.25, 0.5, 0.25}},
		{[3]float32{1, 1, 4}, [3]float32{0.125, 0.5, 0.125}},
	}
	for _, tt := range testData {
		cube := New(tt.sideLengths[0], tt.sideLengths[1], tt.sideLengths[2])
		// bottom
		a := mgl32.Vec3{-tt.expectedCoordinates[0], -tt.expectedCoordinates[1], -tt.expectedCoordinates[2]}
		b := mgl32.Vec3{tt.expectedCoordinates[0], -tt.expectedCoordinates[1], -tt.expectedCoordinates[2]}
		c := mgl32.Vec3{tt.expectedCoordinates[0], -tt.expectedCoordinates[1], tt.expectedCoordinates[2]}
		d := mgl32.Vec3{-tt.expectedCoordinates[0], -tt.expectedCoordinates[1], tt.expectedCoordinates[2]}
		// top
		e := mgl32.Vec3{-tt.expectedCoordinates[0], tt.expectedCoordinates[1], -tt.expectedCoordinates[2]}
		f := mgl32.Vec3{tt.expectedCoordinates[0], tt.expectedCoordinates[1], -tt.expectedCoordinates[2]}
		g := mgl32.Vec3{tt.expectedCoordinates[0], tt.expectedCoordinates[1], tt.expectedCoordinates[2]}
		h := mgl32.Vec3{-tt.expectedCoordinates[0], tt.expectedCoordinates[1], tt.expectedCoordinates[2]}
		points := [24]mgl32.Vec3{
			// bottom
			a, b, c, d,
			// top
			h, g, f, e,
			// front
			e, f, b, a,
			// back
			d, c, g, h,
			// left
			e, a, d, h,
			// right
			b, f, g, c,
		}
		if cube.Points != points {
			t.Error("Invalid points")
			t.Log(cube.Points)
		}
	}
}

func TestNewCube(t *testing.T) {
	cube := NewCube()
	expectedNormals := [6]mgl32.Vec3{
		mgl32.Vec3{0, -1, 0}, // bottom
		mgl32.Vec3{0, 1, 0},  // top
		mgl32.Vec3{0, 0, -1}, // front
		mgl32.Vec3{0, 0, 1},  // back
		mgl32.Vec3{-1, 0, 0}, // left
		mgl32.Vec3{1, 0, 0},  // right
	}
	// bottom
	a := mgl32.Vec3{-0.5, -0.5, -0.5}
	b := mgl32.Vec3{0.5, -0.5, -0.5}
	c := mgl32.Vec3{0.5, -0.5, 0.5}
	d := mgl32.Vec3{-0.5, -0.5, 0.5}
	// top
	e := mgl32.Vec3{-0.5, 0.5, -0.5}
	f := mgl32.Vec3{0.5, 0.5, -0.5}
	g := mgl32.Vec3{0.5, 0.5, 0.5}
	h := mgl32.Vec3{-0.5, 0.5, 0.5}
	points := [24]mgl32.Vec3{
		// bottom
		a, b, c, d,
		// top
		h, g, f, e,
		// front
		e, f, b, a,
		// back
		d, c, g, h,
		// left
		e, a, d, h,
		// right
		b, f, g, c,
	}
	if cube.Normals != expectedNormals {
		t.Error("Invalid normal vectors")
	}
	if cube.Points != points {
		t.Error("Invalid points")
	}
	if !reflect.DeepEqual(cube.Indicies, Indicies) {
		t.Error("Invalid indicies")
	}
}
func TestMeshInput(t *testing.T) {
	cube := NewCube()
	vert, ind := cube.MeshInput()
	if !reflect.DeepEqual(ind, Indicies) {
		t.Error("Invalid indicies")
	}
	if len(vert) != 24 {
		t.Error("Invalid verticies size")
	}
}
func TestColoredMeshInput(t *testing.T) {
	cube := NewCube()
	vert, ind := cube.ColoredMeshInput(Color)
	if !reflect.DeepEqual(ind, Indicies) {
		t.Error("Invalid indicies")
	}
	if len(vert) != 24 {
		t.Error("Invalid verticies size")
	}
}
func TestTexturedColoredMeshInput(t *testing.T) {
	cube := NewCube()
	vert, ind := cube.TexturedColoredMeshInput(Color)
	if !reflect.DeepEqual(ind, Indicies) {
		t.Error("Invalid indicies")
	}
	if len(vert) != 24 {
		t.Error("Invalid verticies size")
	}
}
