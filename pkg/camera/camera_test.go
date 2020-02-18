package camera

import (
	"testing"

	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

func TestNew(t *testing.T) {
	testData := []struct {
		Pos    V.Vector
		LookAt V.Vector
		UpDir  V.Vector
	}{
		{V.Vector{0, 0, 0}, V.Vector{0, 0, 0}, V.Vector{0, 0, 0}},
		{V.Vector{1, 0, 0}, V.Vector{0, 0, 0}, V.Vector{1, 0, 0}},
		{V.Vector{0, 1, 0}, V.Vector{0, 0, 0}, V.Vector{0, 1, 0}},
		{V.Vector{0, 1, 0}, V.Vector{0, 0, 1}, V.Vector{0, 1, 1}},
	}

	for _, tt := range testData {
		camera := New(tt.Pos, tt.LookAt, tt.UpDir)
		// check position
		if camera.Position.X != tt.Pos.X ||
			camera.Position.Y != tt.Pos.Y ||
			camera.Position.Z != tt.Pos.Z {
			t.Error("Invalid camera position")
		}
		// check LookAt
		if camera.LookAt.X != tt.LookAt.X ||
			camera.LookAt.Y != tt.LookAt.Y ||
			camera.LookAt.Z != tt.LookAt.Z {
			t.Error("Invalid camera LookAt")
		}
		// check UpDirection
		if camera.UpDirection.X != tt.UpDir.X ||
			camera.UpDirection.Y != tt.UpDir.Y ||
			camera.UpDirection.Z != tt.UpDir.Z {
			t.Error("Invalid camera UpDirection")
		}
		// check zAxis
		predictedZ := (camera.Position.Add(camera.LookAt.MultiplyScalar(-1))).Normalize()
		if camera.zAxis.X != predictedZ.X ||
			camera.zAxis.Y != predictedZ.Y ||
			camera.zAxis.Z != predictedZ.Z {
			t.Error("Invalid camera zAxis")
		}
		// check xAxis
		predictedX := ((camera.UpDirection.Normalize()).Cross(camera.zAxis)).Normalize()
		if camera.xAxis.X != predictedX.X ||
			camera.xAxis.Y != predictedX.Y ||
			camera.xAxis.Z != predictedX.Z {
			t.Error("Invalid camera xAxis")
		}
		// check yAxis
		predictedY := camera.zAxis.Cross(camera.xAxis)
		if camera.yAxis.X != predictedY.X ||
			camera.yAxis.Y != predictedY.Y ||
			camera.yAxis.Z != predictedY.Z {
			t.Error("Invalid camera yAxis")
		}
	}
}
