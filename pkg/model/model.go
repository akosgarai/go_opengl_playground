package main

import (
	"fmt"

	"github.com/akosgarai/opengl_playground/pkg/interfaces"

	"github.com/udhos/gwob"
)

type Model struct {
	meshes    []interfaces.Mesh
	directory string
}

// NewModelFromFile function gets a path as input. It loads the stuff from the paths to a model
// that the function returns.
func NewModelFromFile(path string) *Model {
	m := &Model{}
	m.loadModel(path)
	return m
}

// Draw function loops over each of the meshes and calls ther Draw function.
func (m *Model) Draw(s interfaces.Shader) {
	for i, _ := range meshes {
		meshes[i].Draw(s)
	}
}
func (m *Model) loadModel(path string) {
	options := &gwob.ObjParserOptions{}

	o, errObj := gwob.NewObjFromFile(path, options)

	if errObj != nil {
		panic(err)
	}

	for _, g := range o.Groups {
		fmt.Println(g)
	}
}
