package transformations

import (
	"math"
	"testing"

	mat "github.com/akosgarai/opengl_playground/pkg/primitives/matrix"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

func TestTranslationMatrix(t *testing.T) {
	testData := []struct {
		Translation vec.Vector
		Expected    mat.Matrix
	}{
		{
			vec.Vector{0, 0, 0},
			*mat.UnitMatrix(),
		},
		{
			vec.Vector{1, 2, 3},
			mat.Matrix{
				[16]float32{
					1.0, 0.0, 0.0, 1.0,
					0.0, 1.0, 0.0, 2.0,
					0.0, 0.0, 1.0, 3.0,
					0.0, 0.0, 0.0, 1.0},
			},
		},
	}
	for _, tt := range testData {
		translated := TranslationMatrix(float32(tt.Translation.X), float32(tt.Translation.Y), float32(tt.Translation.Z))
		for i := 0; i < 16; i++ {
			if translated.Points[i] != tt.Expected.Points[i] {
				t.Log(translated)
				t.Log(tt.Expected)
				t.Error("Value mismatch - TranslationMatrix")
			}
		}
	}
}
func TestRotationTransformationXAxis(t *testing.T) {
	testData := []struct {
		Point         *vec.Vector
		RotationAngle float64
		RotatedPoint  *vec.Vector
	}{
		{&vec.Vector{0, 0, 0}, 0.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{1, 0, 0}, 0.0, &vec.Vector{1, 0, 0}},
		{&vec.Vector{1, 1, 0}, 0.0, &vec.Vector{1, 1, 0}},
		{&vec.Vector{1, 1, 1}, 0.0, &vec.Vector{1, 1, 1}},
		{&vec.Vector{0, 0, 0}, 90.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{1, 0, 0}, 90.0, &vec.Vector{1, 0, 0}},
		{&vec.Vector{0, 1, 0}, math.Pi / 2, &vec.Vector{0, 0, 1}},
		{&vec.Vector{0, 0, 1}, math.Pi / 2, &vec.Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationXMatrix(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != 1.0 ||
			points[1] != 0.0 ||
			points[2] != 0.0 ||
			points[3] != 0.0 ||
			points[4] != 0.0 ||
			points[5] != cosAngle ||
			points[6] != sinAngle ||
			points[7] != 0.0 ||
			points[8] != 0.0 ||
			points[9] != -sinAngle ||
			points[10] != cosAngle ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationX matrix", points, sinAngle, cosAngle)
		}
	}
}
func TestRotationTransformationYAxis(t *testing.T) {
	testData := []struct {
		Point         *vec.Vector
		RotationAngle float64
		RotatedPoint  *vec.Vector
	}{
		{&vec.Vector{0, 0, 0}, 0.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{1, 0, 0}, 0.0, &vec.Vector{1, 0, 0}},
		{&vec.Vector{1, 1, 0}, 0.0, &vec.Vector{1, 1, 0}},
		{&vec.Vector{1, 1, 1}, 0.0, &vec.Vector{1, 1, 1}},
		{&vec.Vector{0, 0, 0}, 90.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{0, 1, 0}, 90.0, &vec.Vector{0, 1, 0}},
		{&vec.Vector{1, 0, 0}, math.Pi / 2, &vec.Vector{0, 0, 1}},
		{&vec.Vector{0, 0, 1}, math.Pi / 2, &vec.Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationYMatrix(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != cosAngle ||
			points[1] != 0.0 ||
			points[2] != -sinAngle ||
			points[3] != 0.0 ||
			points[4] != 0.0 ||
			points[5] != 1.0 ||
			points[6] != 0.0 ||
			points[7] != 0.0 ||
			points[8] != sinAngle ||
			points[9] != 0.0 ||
			points[10] != cosAngle ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationY matrix", points, sinAngle, cosAngle)
		}
	}
}
func TestRotationTransformationZAxis(t *testing.T) {
	testData := []struct {
		Point         *vec.Vector
		RotationAngle float64
		RotatedPoint  *vec.Vector
	}{
		{&vec.Vector{0, 0, 0}, 0.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{1, 0, 0}, 0.0, &vec.Vector{1, 0, 0}},
		{&vec.Vector{1, 1, 0}, 0.0, &vec.Vector{1, 1, 0}},
		{&vec.Vector{1, 1, 1}, 0.0, &vec.Vector{1, 1, 1}},
		{&vec.Vector{0, 0, 0}, 90.0, &vec.Vector{0, 0, 0}},
		{&vec.Vector{0, 0, 1}, 90.0, &vec.Vector{0, 0, 1}},
		{&vec.Vector{1, 0, 0}, math.Pi / 2, &vec.Vector{0, 1, 0}},
		{&vec.Vector{0, 1, 0}, math.Pi / 2, &vec.Vector{1, 0, 0}},
	}
	for _, tt := range testData {
		rotationMatrix := RotationZMatrix(tt.RotationAngle).TransposeMatrix()
		sinAngle, cosAngle := float32(math.Sin(tt.RotationAngle)), float32(math.Cos(tt.RotationAngle))
		points := rotationMatrix.Points
		if points[0] != cosAngle ||
			points[1] != sinAngle ||
			points[2] != 0.0 ||
			points[3] != 0.0 ||
			points[4] != -sinAngle ||
			points[5] != cosAngle ||
			points[6] != 0.0 ||
			points[7] != 0.0 ||
			points[8] != 0.0 ||
			points[9] != 0.0 ||
			points[10] != 1.0 ||
			points[11] != 0.0 ||
			points[12] != 0.0 ||
			points[13] != 0.0 ||
			points[14] != 0.0 ||
			points[15] != 1.0 {
			t.Error("Invalid rotationY matrix", points, sinAngle, cosAngle)
		}
	}
}
