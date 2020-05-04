package rectangle

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type testShader struct {
	HasTextureValue bool
}

func (t testShader) Use() {
}
func (t testShader) SetUniformMat4(s string, m mgl32.Mat4) {
}
func (t testShader) DrawTriangles(i int32) {
}
func (t testShader) Close(i int) {
}
func (t testShader) VertexAttribPointer(i uint32, c int32, s int32, o int) {
}
func (t testShader) BindVertexArray() {
}
func (t testShader) BindBufferData(d []float32) {
}
func (t testShader) HasTexture() bool {
	return t.HasTextureValue
}

func TestNew(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)

	if square.speed != 0.0 {
		t.Error("Speed should be 0")
	}
	if square.direction.X() != 0.0 || square.direction.Y() != 0.0 || square.direction.Z() != 0.0 {
		t.Error("Direction vector is not 0")
	}

	for i := 0; i < 4; i++ {
		if square.points[i] != points[i] {
			t.Error("Mismatch in the coordinates")
		}
		if square.colors[i] != colors[i] {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestSetColor(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	square.SetColor(newColor)

	for i := 0; i < 4; i++ {
		if square.colors[i] != newColor {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestSetIndexColor(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	square.SetIndexColor(0, newColor)

	if square.colors[0] != newColor {
		t.Error("Mismatch in the new color")
	}
	for i := 1; i < 4; i++ {
		if square.colors[i] != colors[i] {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestLog(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	log := square.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetupVao(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.setupVao()
	if len(square.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVaoWithoutTexture(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.buildVaoWithoutTexture()
	if len(square.vao.Get()) != 36 {
		t.Errorf("Invalid number of items in the vao. Instead of '36', we have '%d'.", len(square.vao.Get()))
	}
}
func TestBuildVaoWithTexture(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = true
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.buildVaoWithTexture()
	if len(square.vao.Get()) != 48 {
		t.Errorf("Invalid number of items in the vao. Instead of '48', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDraw(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.Draw()
	if len(square.vao.Get()) != 36 {
		t.Errorf("Invalid number of items in the vao. Instead of '36', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDrawWithTexture(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = true
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.Draw()
	if len(square.vao.Get()) != 48 {
		t.Errorf("Invalid number of items in the vao. Instead of '48', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDrawWithUniforms(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(square.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawWithUniformsTextures(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = true
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(square.vao.Get()) != 48 {
		t.Errorf("Vao should be 48 long. Instead of it, it's '%d'", len(square.vao.Get()))
	}
}
func TestUpdate(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	square.SetDirection(mgl32.Vec3{1, 0, 0})
	square.SetSpeed(1)
	dt := 10.0
	square.Update(dt)
	if square.points[0].X() != 10.0 || square.points[0].Y() != 0.0 || square.points[0].Z() != 0.0 {
		t.Error("Invalid position for p0")
	}
	if square.points[1].X() != 11.0 || square.points[1].Y() != 0.0 || square.points[1].Z() != 0.0 {
		t.Error("Invalid position for p1")
	}
	if square.points[2].X() != 11.0 || square.points[2].Y() != 1.0 || square.points[2].Z() != 0.0 {
		t.Error("Invalid position for p2")
	}
	if square.points[3].X() != 10.0 || square.points[3].Y() != 1.0 || square.points[3].Z() != 0.0 {
		t.Error("Invalid position for p3")
	}
}
func TestSetDirection(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	newDirection := mgl32.Vec3{1, 1, 0}
	square.SetDirection(newDirection)

	if square.direction != newDirection {
		t.Error("Mismatch in the direction")
	}
}
func TestSetIndexDirection(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	square.SetIndexDirection(0, 1)

	if square.direction.X() != 1.0 || square.direction.Y() != 0.0 || square.direction.Z() != 0.0 {
		t.Error("Mismatch in the direction")
	}
}
func TestSetSpeed(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	square.SetSpeed(10)

	if square.speed != 10.0 {
		t.Error("Mismatch in the speed")
	}
}
func TestCoordinates(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	coords := square.Coordinates()
	if coords != points {
		t.Error("Mismatch in coordinates")
	}
}
func TestColors(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	cols := square.Colors()
	if cols != colors {
		t.Error("Mismatch in colors")
	}
}
func TestSetupExternalVao(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("VAO should be empty")
	}
	outerVao := square.SetupExternalVao(vao.NewVAO())
	if len(square.vao.Get()) != 0 {
		t.Error("VAO should be empty")
	}
	if len(outerVao.Get()) == 0 {
		t.Error("VAO shouldn't be empty")
	}
}
func TestSetPrecision(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	newPrecision := 100
	square.SetPrecision(newPrecision)

	if square.precision != newPrecision {
		t.Error("Mismatch in the precision")
	}
}

func TestSetRotationAngle(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	angle := float32(2.0)
	square.SetAngle(angle)

	if square.angle != angle {
		t.Errorf("Mismatch in the angle instead of '%f', we have '%f'.", angle, square.angle)
	}
}
func TestSetRotationAxis(t *testing.T) {
	points := [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	colors := [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	var shader testShader
	shader.HasTextureValue = false
	square := New(points, colors, shader)
	axis := mgl32.Vec3{0, 1, 0}
	square.SetAxis(axis)

	if square.axis != axis {
		t.Errorf("Mismatch in the axis")
	}
}
