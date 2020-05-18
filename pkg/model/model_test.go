package main

import (
	"testing"
)

const (
	InvalidFilename = "not-existing-file.obj"
	ValidFilename   = "testdata/test_cube.obj"
)

var ()

func TestNewModelFromFile(t *testing.T) {
	//m := NewModelFromFile("InvalidFilePath")
	//m := NewModelFromFile("ValidFilePath")
	t.Skip("Unimplemented")

}
func TestDraw(t *testing.T) {
	t.Skip("Unimplemented")
}
func Test_loadModel_InvalidFilename(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("loadModel should have panicked due to the missing file!")
			}
		}()
		m := &Model{}
		m.loadModel(InvalidFilename)
	}()
}
func Test_loadModel_ValidFilename(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
				t.Errorf("loadModel shouldn't have panicked!")
			}
		}()
		glWrapper.InitOpenGL()
		m := &Model{}
		m.loadModel(ValidFilename)
	}()
}
