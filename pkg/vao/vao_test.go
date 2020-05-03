package vao

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewVAO(t *testing.T) {
	vao := NewVAO()
	// should be empty
	if len(vao.vao) != 0 {
		t.Error("The new vao should be empty")
	}
}
func TestAppendVector(t *testing.T) {
	vao := NewVAO()
	vector := mgl32.Vec3{1, 2, 3}
	vao.appendVector(vector)
	if len(vao.vao) != 3 {
		t.Error("appendVector should add 3 element to the vao.")
	}
	if vao.vao[0] != 1 {
		t.Error("appendVector should add 1 as first element to the vao.")
	}
	if vao.vao[1] != 2 {
		t.Error("appendVector should add 2 as second element to the vao.")
	}
	if vao.vao[2] != 3 {
		t.Error("appendVector should add 3 as third element to the vao.")
	}
}
func TestAppendVectors(t *testing.T) {
	vao := NewVAO()
	vector1 := mgl32.Vec3{1, 2, 3}
	vector2 := mgl32.Vec3{4, 5, 6}
	vao.AppendVectors(vector1, vector2)
	if len(vao.vao) != 6 {
		t.Error("AppendVectors should add 6 element to the vao.")
	}
	if vao.vao[0] != 1 {
		t.Error("AppendVectors should add 1 as first element to the vao.")
	}
	if vao.vao[1] != 2 {
		t.Error("AppendVectors should add 2 as second element to the vao.")
	}
	if vao.vao[2] != 3 {
		t.Error("AppendVectors should add 3 as third element to the vao.")
	}
	if vao.vao[3] != 4 {
		t.Error("AppendVectors should add 4 as 4. element to the vao.")
	}
	if vao.vao[4] != 5 {
		t.Error("AppendVectors should add 5 as 5. element to the vao.")
	}
	if vao.vao[5] != 6 {
		t.Error("AppendVectors should add 6 as 6. element to the vao.")
	}
}
func TestGet(t *testing.T) {
	vao := NewVAO()
	vector1 := mgl32.Vec3{1, 2, 3}
	vector2 := mgl32.Vec3{4, 5, 6}
	vao.AppendVectors(vector1, vector2)
	if len(vao.Get()) != 6 {
		t.Error("Get should return 6 items.")
	}
	result := vao.Get()
	if result[0] != 1 {
		t.Error("Get should return 1 as 1. item.")
	}
	if result[1] != 2 {
		t.Error("Get should return 2 as 2. item.")
	}
	if result[2] != 3 {
		t.Error("Get should return 3 as 3. item.")
	}
	if result[3] != 4 {
		t.Error("Get should return 4 as 4. item.")
	}
	if result[4] != 5 {
		t.Error("Get should return 5 as 5. item.")
	}
	if result[5] != 6 {
		t.Error("Get should return 6 as 6. item.")
	}
}
func TestClear(t *testing.T) {
	vao := NewVAO()
	vector1 := mgl32.Vec3{1, 2, 3}
	vector2 := mgl32.Vec3{4, 5, 6}
	vao.AppendVectors(vector1, vector2)
	if len(vao.Get()) != 6 {
		t.Error("Get should return 6 items.")
	}
	vao.Clear()
	if len(vao.Get()) != 0 {
		t.Error("Clear should remove every item from the vao.")
	}
}
func TestAppendPoint(t *testing.T) {
	vao := NewVAO()
	vector1 := mgl32.Vec3{1, 2, 3}
	vector2 := mgl32.Vec3{4, 5, 6}
	size := float32(7.0)
	vao.AppendPoint(vector1, vector2, size)
	if len(vao.vao) != 7 {
		t.Error("AppendVectors should add 7 element to the vao.")
	}
	if vao.vao[0] != 1 {
		t.Error("AppendVectors should add 1 as first element to the vao.")
	}
	if vao.vao[1] != 2 {
		t.Error("AppendVectors should add 2 as second element to the vao.")
	}
	if vao.vao[2] != 3 {
		t.Error("AppendVectors should add 3 as third element to the vao.")
	}
	if vao.vao[3] != 4 {
		t.Error("AppendVectors should add 4 as 4. element to the vao.")
	}
	if vao.vao[4] != 5 {
		t.Error("AppendVectors should add 5 as 5. element to the vao.")
	}
	if vao.vao[5] != 6 {
		t.Error("AppendVectors should add 6 as 6. element to the vao.")
	}
	if vao.vao[6] != 7 {
		t.Error("AppendVectors should add 7 as 7. element to the vao.")
	}
}
func TestAppendTextureVectors(t *testing.T) {
}
