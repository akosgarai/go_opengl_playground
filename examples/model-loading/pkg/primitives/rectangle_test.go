package primitives

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewSquare(t *testing.T) {
	square := NewSquare()

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Error("Unexpected points.")
	}
	if square.Normal != ExpectedNormal {
		t.Error("Invalid normal vector")
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Error("Invalid normal vs calculated normal")
	}
}
func TestNewOneAsScale(t *testing.T) {
	width := float32(1)
	height := float32(1)
	square := New(width, height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Error("Unexpected points.")
	}
	if square.Normal != ExpectedNormal {
		t.Error("Invalid normal vector")
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Error("Invalid normal vs calculated normal")
	}
}
func TestNewLowScale(t *testing.T) {
	width := float32(2)
	height := float32(1)
	square := New(width, height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.25},
		mgl32.Vec3{0.5, 0, -0.25},
		mgl32.Vec3{0.5, 0, 0.25},
		mgl32.Vec3{-0.5, 0, 0.25},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Error("Unexpected points.")
	}
	if square.Normal != ExpectedNormal {
		t.Error("Invalid normal vector")
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Error("Invalid normal vs calculated normal")
	}
}
func TestNewHighScale(t *testing.T) {
	width := float32(1)
	height := float32(2)
	square := New(width, height)
	t.Log(width / height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.25, 0, -0.5},
		mgl32.Vec3{0.25, 0, -0.5},
		mgl32.Vec3{0.25, 0, 0.5},
		mgl32.Vec3{-0.25, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Error("Unexpected points.")
	}
	if square.Normal != ExpectedNormal {
		t.Error("Invalid normal vector")
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Error("Invalid normal vs calculated normal")
	}
}
