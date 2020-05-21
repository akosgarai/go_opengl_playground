package export

import (
	"os"
	"testing"

	"github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/akosgarai/opengl_playground/pkg/interfaces"
	"github.com/akosgarai/opengl_playground/pkg/mesh"
	"github.com/akosgarai/opengl_playground/pkg/primitives/cuboid"
	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/rectangle"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	"github.com/akosgarai/opengl_playground/pkg/primitives/vertex"
	"github.com/akosgarai/opengl_playground/pkg/testhelper"
	"github.com/akosgarai/opengl_playground/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	glWrapper testhelper.GLWrapperMock
)

func TestNew(t *testing.T) {
	var meshes []interfaces.Mesh
	square := rectangle.NewSquare()
	v, i := square.MeshInput()
	materialMesh := mesh.NewMaterialMesh(v, i, material.Jade, glWrapper)
	meshes = append(meshes, materialMesh)
	exporter := New(meshes)
	if len(exporter.meshes) != 1 {
		t.Error("Invalid mesh length")
	}
}
func TestMaterialExport(t *testing.T) {
	exporter := &Export{
		materials: []Mtl{},
	}
	result := exporter.materialExport()
	if result != "" {
		t.Error("Without materials it shouldn't generate output.")
	}
	mat := Mtl{
		Ka:   [3]float32{0, 1, 0},
		Kd:   [3]float32{0, 1, 0},
		Ks:   [3]float32{1, 1, 1},
		Ns:   float32(32),
		Name: "TestMaterial",
	}
	exporter.materials = append(exporter.materials, mat)
	result = exporter.materialExport()
	if result != "newmtl TestMaterial\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\n\n" {
		t.Error("Invalid material string")
	}
	mat2 := Mtl{
		Ka:   [3]float32{0, 1, 0},
		Kd:   [3]float32{0, 1, 0},
		Ks:   [3]float32{1, 1, 1},
		Ns:   float32(32),
		Name: "TestMaterial 2",
	}
	exporter.materials = append(exporter.materials, mat2)
	result = exporter.materialExport()
	if result != "newmtl TestMaterial\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\n\nnewmtl TestMaterial 2\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\n\n" {
		t.Error("Invalid material string")
	}
}
func TestMaterialTextureExport(t *testing.T) {
	exporter := &Export{
		materials: []Mtl{},
	}
	result := exporter.materialExport()
	if result != "" {
		t.Error("Without materials it shouldn't generate output.")
	}
	mat := Mtl{
		Ka:    [3]float32{0, 1, 0},
		Kd:    [3]float32{0, 1, 0},
		Ks:    [3]float32{1, 1, 1},
		Ns:    float32(32),
		Name:  "TestMaterial",
		MapKa: "ambient.jpg",
		MapKd: "diffuse.jpg",
		MapKs: "specular.jpg",
	}
	exporter.materials = append(exporter.materials, mat)
	result = exporter.materialExport()
	if result != "newmtl TestMaterial\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka ambient.jpg\nmap_Kd diffuse.jpg\nmap_Ks specular.jpg\n\n" {
		t.Error("Invalid material string")
	}
	mat2 := Mtl{
		Ka:    [3]float32{0, 1, 0},
		Kd:    [3]float32{0, 1, 0},
		Ks:    [3]float32{1, 1, 1},
		Ns:    float32(32),
		Name:  "TestMaterial 2",
		MapKa: "ambient.jpg",
		MapKd: "diffuse.jpg",
		MapKs: "specular.jpg",
	}
	exporter.materials = append(exporter.materials, mat2)
	result = exporter.materialExport()
	if result != "newmtl TestMaterial\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka ambient.jpg\nmap_Kd diffuse.jpg\nmap_Ks specular.jpg\n\nnewmtl TestMaterial 2\nKa 0.0000000000 1.0000000000 0.0000000000\nKd 0.0000000000 1.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka ambient.jpg\nmap_Kd diffuse.jpg\nmap_Ks specular.jpg\n\n" {
		t.Error("Invalid material string")
	}
}
func TestObjectExportRectangle(t *testing.T) {
	exporter := &Export{
		objects: []Obj{},
	}
	result := exporter.objectExport()
	if result != "" {
		t.Error("Without objects it shouldn't generate output.")
	}
	obj := Obj{
		Name: "testObject",
		V: [][3]float32{
			[3]float32{-0.5, 0, -0.5},
			[3]float32{-0.5, 0, 0.5},
			[3]float32{0.5, 0, 0.5},
			[3]float32{0.5, 0, -0.5},
		},
		Normal: [][3]float32{
			[3]float32{0, 1, 0},
		},
		Indices: []string{
			"1//1",
			"2//1",
			"3//1",
			"4//1",
			"2//1",
			"1//1",
		},
		HasFaces:     true,
		MaterialName: "TestMaterial",
	}
	exporter.objects = append(exporter.objects, obj)

	result = exporter.objectExport()
	if result != "o testObject\nv -0.5000000000 0.0000000000 -0.5000000000\nv -0.5000000000 0.0000000000 0.5000000000\nv 0.5000000000 0.0000000000 0.5000000000\nv 0.5000000000 0.0000000000 -0.5000000000\nvn 0.0000000000 1.0000000000 0.0000000000\nusemtl TestMaterial\ns off\nf 1//1 2//1 3//1\nf 4//1 2//1 1//1\n\n" {
		t.Error("Invalid object string")
	}
}
func TestExportTexturedColoredMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	colors := []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
		mgl32.Vec3{1.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 0.0},
		mgl32.Vec3{0.0, 1.0, 1.0},
		mgl32.Vec3{0.0, 0.0, 1.0},
		mgl32.Vec3{1.0, 0.0, 1.0},
	}
	cube := cuboid.NewCube()
	v, i := cube.TexturedColoredMeshInput(colors)
	tcMesh := mesh.NewTexturedColoredMesh(v, i, tex, colors, glWrapper)
	meshes = append(meshes, tcMesh)
	exporter := New(meshes)
	result := exporter.Export("./tests")
	if result != nil {
		t.Error("Textured colored mesh should be handled as textured colored mesh")
	}
	os.Remove("./tests/material.mat")
	os.Remove("./tests/object.obj")
}
func TestProcessTexturedColoredMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	colors := []mgl32.Vec3{
		mgl32.Vec3{1.0, 0.0, 0.0},
	}
	cube := cuboid.NewCube()
	v, i := cube.TexturedColoredMeshInput(colors)
	tcMesh := mesh.NewTexturedColoredMesh(v, i, tex, colors, glWrapper)
	meshes = append(meshes, tcMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processTexturedColorMesh(tcMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Textured_Color_Material_0\nKa 1.0000000000 0.0000000000 0.0000000000\nKd 1.0000000000 0.0000000000 0.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka test-image.jpg\nmap_Kd test-image.jpg\nmap_Ks test-image.jpg\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Texture_Color_Material_Object_0\nv -0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 -0.5000000000\nv -0.5000000000 0.5000000000 -0.5000000000\nvt 0.0000000000 1.0000000000\nvt 1.0000000000 1.0000000000\nvt 1.0000000000 0.0000000000\nvt 0.0000000000 0.0000000000\nusemtl Textured_Color_Material_0\ns off\nf 1/1 2/2 3/3\nf 1/1 3/3 4/4\nf 5/1 6/2 7/3\nf 5/1 7/3 8/4\nf 8/1 7/2 2/3\nf 8/1 2/3 1/4\nf 4/1 3/2 6/3\nf 4/1 6/3 5/4\nf 8/1 1/2 4/3\nf 8/1 4/3 5/4\nf 2/1 7/2 6/3\nf 2/1 6/3 3/4\n\n" {
		t.Error("Invalid object string")
	}

}
func TestExportTexturedMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	spherePrimitive := sphere.New(20)
	v, i := spherePrimitive.TexturedMeshInput()
	texturedMesh := mesh.NewTexturedMesh(v, i, tex, glWrapper)
	meshes = append(meshes, texturedMesh)
	exporter := New(meshes)
	result := exporter.Export("./tests")
	if result != nil {
		t.Error("Textured mesh should be handled as textured mesh")
	}
	os.Remove("./tests/material.mat")
	os.Remove("./tests/object.obj")
}
func TestProcessTexturedGoodNamesMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	cube := cuboid.NewCube()
	v, i := cube.MeshInput()
	texturedMesh := mesh.NewTexturedMesh(v, i, tex, glWrapper)
	meshes = append(meshes, texturedMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processTextureMesh(texturedMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Texture_Material_0\nKa 1.0000000000 1.0000000000 1.0000000000\nKd 1.0000000000 1.0000000000 1.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka test-image.jpg\nmap_Kd test-image.jpg\nmap_Ks test-image.jpg\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Material_Object_0\nv -0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 -0.5000000000\nv -0.5000000000 0.5000000000 -0.5000000000\nvn 0.0000000000 -1.0000000000 0.0000000000\nvn 0.0000000000 1.0000000000 0.0000000000\nvn 0.0000000000 0.0000000000 -1.0000000000\nvn 0.0000000000 0.0000000000 1.0000000000\nvn -1.0000000000 0.0000000000 0.0000000000\nvn 1.0000000000 0.0000000000 0.0000000000\nvt 0.0000000000 1.0000000000\nvt 1.0000000000 1.0000000000\nvt 1.0000000000 0.0000000000\nvt 0.0000000000 0.0000000000\nusemtl Texture_Material_0\ns off\nf 1/1/1 2/2/1 3/3/1\nf 1/1/1 3/3/1 4/4/1\nf 5/1/2 6/2/2 7/3/2\nf 5/1/2 7/3/2 8/4/2\nf 8/1/3 7/2/3 2/3/3\nf 8/1/3 2/3/3 1/4/3\nf 4/1/4 3/2/4 6/3/4\nf 4/1/4 6/3/4 5/4/4\nf 8/1/5 1/2/5 4/3/5\nf 8/1/5 4/3/5 5/4/5\nf 2/1/6 7/2/6 6/3/6\nf 2/1/6 6/3/6 3/4/6\n\n" {
		t.Error("Invalid object string")
	}

}
func TestProcessTexturedNoSpecularMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	cube := cuboid.NewCube()
	v, i := cube.MeshInput()
	texturedMesh := mesh.NewTexturedMesh(v, i, tex, glWrapper)
	meshes = append(meshes, texturedMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processTextureMesh(texturedMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Texture_Material_0\nKa 1.0000000000 1.0000000000 1.0000000000\nKd 1.0000000000 1.0000000000 1.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka test-image.jpg\nmap_Kd test-image.jpg\nmap_Ks test-image.jpg\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Material_Object_0\nv -0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 -0.5000000000\nv -0.5000000000 0.5000000000 -0.5000000000\nvn 0.0000000000 -1.0000000000 0.0000000000\nvn 0.0000000000 1.0000000000 0.0000000000\nvn 0.0000000000 0.0000000000 -1.0000000000\nvn 0.0000000000 0.0000000000 1.0000000000\nvn -1.0000000000 0.0000000000 0.0000000000\nvn 1.0000000000 0.0000000000 0.0000000000\nvt 0.0000000000 1.0000000000\nvt 1.0000000000 1.0000000000\nvt 1.0000000000 0.0000000000\nvt 0.0000000000 0.0000000000\nusemtl Texture_Material_0\ns off\nf 1/1/1 2/2/1 3/3/1\nf 1/1/1 3/3/1 4/4/1\nf 5/1/2 6/2/2 7/3/2\nf 5/1/2 7/3/2 8/4/2\nf 8/1/3 7/2/3 2/3/3\nf 8/1/3 2/3/3 1/4/3\nf 4/1/4 3/2/4 6/3/4\nf 4/1/4 6/3/4 5/4/4\nf 8/1/5 1/2/5 4/3/5\nf 8/1/5 4/3/5 5/4/5\nf 2/1/6 7/2/6 6/3/6\nf 2/1/6 6/3/6 3/4/6\n\n" {
		t.Error("Invalid object string")
	}

}
func TestProcessTexturedNoDiffuseMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	var tex texture.Textures
	tex.AddTexture("test-image.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	cube := cuboid.NewCube()
	v, i := cube.MeshInput()
	texturedMesh := mesh.NewTexturedMesh(v, i, tex, glWrapper)
	meshes = append(meshes, texturedMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processTextureMesh(texturedMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Texture_Material_0\nKa 1.0000000000 1.0000000000 1.0000000000\nKd 1.0000000000 1.0000000000 1.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\nmap_Ka test-image.jpg\nmap_Kd test-image.jpg\nmap_Ks test-image.jpg\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Material_Object_0\nv -0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 -0.5000000000\nv 0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 -0.5000000000 0.5000000000\nv -0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 0.5000000000\nv 0.5000000000 0.5000000000 -0.5000000000\nv -0.5000000000 0.5000000000 -0.5000000000\nvn 0.0000000000 -1.0000000000 0.0000000000\nvn 0.0000000000 1.0000000000 0.0000000000\nvn 0.0000000000 0.0000000000 -1.0000000000\nvn 0.0000000000 0.0000000000 1.0000000000\nvn -1.0000000000 0.0000000000 0.0000000000\nvn 1.0000000000 0.0000000000 0.0000000000\nvt 0.0000000000 1.0000000000\nvt 1.0000000000 1.0000000000\nvt 1.0000000000 0.0000000000\nvt 0.0000000000 0.0000000000\nusemtl Texture_Material_0\ns off\nf 1/1/1 2/2/1 3/3/1\nf 1/1/1 3/3/1 4/4/1\nf 5/1/2 6/2/2 7/3/2\nf 5/1/2 7/3/2 8/4/2\nf 8/1/3 7/2/3 2/3/3\nf 8/1/3 2/3/3 1/4/3\nf 4/1/4 3/2/4 6/3/4\nf 4/1/4 6/3/4 5/4/4\nf 8/1/5 1/2/5 4/3/5\nf 8/1/5 4/3/5 5/4/5\nf 2/1/6 7/2/6 6/3/6\nf 2/1/6 6/3/6 3/4/6\n\n" {
		t.Error("Invalid object string")
	}

}
func TestExportColorMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	square := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{1, 1, 1}}
	v, i := square.ColoredMeshInput(col)
	colorMesh := mesh.NewColorMesh(v, i, col, glWrapper)
	meshes = append(meshes, colorMesh)
	exporter := New(meshes)
	result := exporter.Export("./tests")
	if result != nil {
		t.Error("Color mesh should be handled as color mesh")
	}
	os.Remove("./tests/material.mat")
	os.Remove("./tests/object.obj")
}
func TestProcessColorMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	square := rectangle.NewSquare()
	col := []mgl32.Vec3{mgl32.Vec3{1, 1, 1}}
	v, i := square.ColoredMeshInput(col)
	colorMesh := mesh.NewColorMesh(v, i, col, glWrapper)
	meshes = append(meshes, colorMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processColorMesh(colorMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Color_Material_0\nKa 1.0000000000 1.0000000000 1.0000000000\nKd 1.0000000000 1.0000000000 1.0000000000\nKs 1.0000000000 1.0000000000 1.0000000000\nNs 32.0000000000\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Color_Material_Object_0\nv -0.5000000000 0.0000000000 -0.5000000000\nv 0.5000000000 0.0000000000 -0.5000000000\nv 0.5000000000 0.0000000000 0.5000000000\nv -0.5000000000 0.0000000000 0.5000000000\nusemtl Color_Material_0\ns off\nf 1 2 3\nf 1 3 4\n\n" {
		t.Error("Invalid object string")
	}
}
func TestExportMaterialMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	square := rectangle.NewSquare()
	v, i := square.MeshInput()
	materialMesh := mesh.NewMaterialMesh(v, i, material.Jade, glWrapper)
	meshes = append(meshes, materialMesh)
	exporter := New(meshes)
	result := exporter.Export("./tests")
	if result != nil {
		t.Error("Material mesh should be handled as material mesh")
	}
	os.Remove("./tests/material.mat")
	os.Remove("./tests/object.obj")
}
func TestProcessMaterialMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	square := rectangle.NewSquare()
	v, i := square.MeshInput()
	materialMesh := mesh.NewMaterialMesh(v, i, material.Jade, glWrapper)
	meshes = append(meshes, materialMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processMaterialMesh(materialMesh)
	if len(exporter.materials) != 1 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.materialExport()
	if result != "newmtl Material_0\nKa 0.1350000054 0.2224999964 0.1574999988\nKd 0.5400000215 0.8899999857 0.6299999952\nKs 0.3162280023 0.3162280023 0.3162280023\nNs 0.1000000015\n\n" {
		t.Error("Invalid material string")
	}
	result = exporter.objectExport()
	if result != "o Material_Object_0\nv -0.5000000000 0.0000000000 -0.5000000000\nv 0.5000000000 0.0000000000 -0.5000000000\nv 0.5000000000 0.0000000000 0.5000000000\nv -0.5000000000 0.0000000000 0.5000000000\nvn 0.0000000000 -1.0000000000 0.0000000000\nusemtl Material_0\ns off\nf 1//1 2//1 3//1\nf 1//1 3//1 4//1\n\n" {
		t.Error("Invalid object string")
	}
}
func TestExportPointMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	pointMesh := mesh.NewPointMesh(glWrapper)
	meshes = append(meshes, pointMesh)
	exporter := New(meshes)
	result := exporter.Export("./tests")
	if result != nil {
		t.Error("Point mesh should be handled as point mesh")
	}
	os.Remove("./tests/material.mat")
	os.Remove("./tests/object.obj")
}
func TestProcessPointMesh(t *testing.T) {
	var meshes []interfaces.Mesh
	pointMesh := mesh.NewPointMesh(glWrapper)
	vert := vertex.Vertex{Position: mgl32.Vec3{1, 0, 0}}
	pointMesh.AddVertex(vert)
	meshes = append(meshes, pointMesh)
	exporter := New(meshes)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 0 {
		t.Error("Invalid object length")
	}
	exporter.processPointMesh(pointMesh)
	if len(exporter.materials) != 0 {
		t.Error("Invalid material length")
	}
	if len(exporter.objects) != 1 {
		t.Error("Invalid object length")
	}
	result := exporter.objectExport()
	if result != "o Point_Object_0\nv 1.0000000000 0.0000000000 0.0000000000\np 1\n\n" {
		t.Error("Invalid object string")

	}
}
