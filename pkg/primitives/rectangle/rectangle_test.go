package rectangle

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/vao"
)

var (
	DefaultCoordinates = [4]mgl32.Vec3{
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 1, 0},
		mgl32.Vec3{0, 1, 0},
	}
	DefaultColors = [4]mgl32.Vec3{
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
		mgl32.Vec3{1, 0, 0},
	}
	DefaultSpeed     = float32(0.0)
	DefaultDirection = mgl32.Vec3{0, 0, 0}
)

type testShader struct {
	HasTextureValue bool
}

func (t testShader) Use() {
}
func (t testShader) SetUniform1f(s string, f1 float32) {
}
func (t testShader) SetUniform3f(s string, f1, f2, f3 float32) {
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

var shader testShader

func TestNew(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)

	if square.speed != DefaultSpeed {
		t.Errorf("Speed should be '%f, but it is '%f'.", DefaultSpeed, square.speed)
	}
	if square.direction != DefaultDirection {
		t.Error("Direction vector is not 0")
	}

	if square.points != DefaultCoordinates {
		t.Error("Mismatch in the coordinates")
	}
	if square.colors != DefaultColors {
		t.Error("Mismatch in the colors")
	}

}
func TestNewSquare(t *testing.T) {
	shader.HasTextureValue = false

	color := mgl32.Vec3{1, 0, 0}
	testData := []struct {
		ExpectedPoints [4]mgl32.Vec3
		Normal         mgl32.Vec3
	}{
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-0.5, 0, -0.5},
				mgl32.Vec3{-0.5, 0, 0.5},
				mgl32.Vec3{0.5, 0, 0.5},
				mgl32.Vec3{0.5, 0, -0.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-0.5, 5, -0.5},
				mgl32.Vec3{-0.5, 5, 0.5},
				mgl32.Vec3{0.5, 5, 0.5},
				mgl32.Vec3{0.5, 5, -0.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-5.5, 5, -5.5},
				mgl32.Vec3{-5.5, 5, 5.5},
				mgl32.Vec3{5.5, 5, 5.5},
				mgl32.Vec3{5.5, 5, -5.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-3.5, -0.5, -3.5},
				mgl32.Vec3{-3.5, -0.5, -2.5},
				mgl32.Vec3{-2.5, -0.5, -2.5},
				mgl32.Vec3{-2.5, -0.5, -3.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-5.5, 0, -3.5},
				mgl32.Vec3{-5.5, 0, -1.5},
				mgl32.Vec3{-3.5, 0, -1.5},
				mgl32.Vec3{-3.5, 0, -3.5},
			},
			mgl32.Vec3{0, 1, 0}.Normalize(),
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-5.5, 5, -3.5},
				mgl32.Vec3{-5.5, 5, -1.5},
				mgl32.Vec3{-3.5, 5, -1.5},
				mgl32.Vec3{-3.5, 5, -3.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
		{
			[4]mgl32.Vec3{
				mgl32.Vec3{-7.5, -3.5, -5.5},
				mgl32.Vec3{-7.5, -3.5, -4.5},
				mgl32.Vec3{-6.5, -3.5, -4.5},
				mgl32.Vec3{-6.5, -3.5, -5.5},
			},
			mgl32.Vec3{0, 1, 0},
		},
	}
	for index, tt := range testData {
		square := NewSquare(tt.ExpectedPoints[1], tt.ExpectedPoints[3], tt.Normal, color, shader)
		for i := 0; i < 4; i++ {
			// It has to be checked in this way. In the practice, the '-0.5' was calculated as '-0.49999997'.
			if !tt.ExpectedPoints[i].ApproxEqualThreshold(square.points[i], 0.003) {
				rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(90.0), tt.Normal)
				t.Log("RotationMatrix:")
				t.Log(rotationMatrix)
				p2 := mgl32.TransformCoordinate(tt.ExpectedPoints[1], rotationMatrix)
				p0 := mgl32.TransformCoordinate(tt.ExpectedPoints[3], rotationMatrix)
				t.Log("p0")
				t.Log(p0)
				t.Log("p2")
				t.Log(p2)
				t.Log("expected")
				t.Log(tt.ExpectedPoints)
				t.Log("got")
				t.Log(square.points)
				t.Fatalf("TC%d - Vectors are not equal. %d\n", index, i)
			}
		}
	}
}
func TestSetColor(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	square.SetColor(newColor)

	for i := 0; i < 4; i++ {
		if square.colors[i] != newColor {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestSetIndexColor(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	newColor := mgl32.Vec3{1, 1, 0}
	square.SetIndexColor(0, newColor)

	if square.colors[0] != newColor {
		t.Error("Mismatch in the new color")
	}
	for i := 1; i < 4; i++ {
		if square.colors[i] != DefaultColors[i] {
			t.Error("Mismatch in the colors")
		}
	}
}
func TestLog(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	log := square.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestSetupVao(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.setupVao()
	if len(square.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestBuildVaoWithoutTexture(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.buildVaoWithoutTexture()
	if len(square.vao.Get()) != 36 {
		t.Errorf("Invalid number of items in the vao. Instead of '36', we have '%d'.", len(square.vao.Get()))
	}
}
func TestBuildVaoWithTexture(t *testing.T) {
	shader.HasTextureValue = true
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.buildVaoWithTexture()
	if len(square.vao.Get()) != 48 {
		t.Errorf("Invalid number of items in the vao. Instead of '48', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDraw(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.Draw()
	if len(square.vao.Get()) != 36 {
		t.Errorf("Invalid number of items in the vao. Instead of '36', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDrawWithTexture(t *testing.T) {
	shader.HasTextureValue = true
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.Draw()
	if len(square.vao.Get()) != 48 {
		t.Errorf("Invalid number of items in the vao. Instead of '48', we have '%d'.", len(square.vao.Get()))
	}
}
func TestDrawWithUniforms(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(square.vao.Get()) == 0 {
		t.Error("Vao is empty after the first setup.")
	}
}
func TestDrawWithUniformsTextures(t *testing.T) {
	shader.HasTextureValue = true
	square := New(DefaultCoordinates, DefaultColors, shader)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(square.vao.Get()) != 48 {
		t.Errorf("Vao should be 48 long. Instead of it, it's '%d'", len(square.vao.Get()))
	}
}
func TestUpdate(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
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
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	newDirection := mgl32.Vec3{1, 1, 0}
	square.SetDirection(newDirection)

	if square.direction != newDirection {
		t.Error("Mismatch in the direction")
	}
	if square.GetDirection() != newDirection {
		t.Error("Mismatch in the direction")
	}
}
func TestSetIndexDirection(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	square.SetIndexDirection(0, 1)

	if square.direction.X() != 1.0 || square.direction.Y() != 0.0 || square.direction.Z() != 0.0 {
		t.Error("Mismatch in the direction")
	}
}
func TestSetSpeed(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	square.SetSpeed(10)

	if square.speed != 10.0 {
		t.Error("Mismatch in the speed")
	}
}
func TestCoordinates(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	coords := square.Coordinates()
	if coords != DefaultCoordinates {
		t.Error("Mismatch in coordinates")
	}
}
func TestColors(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	cols := square.Colors()
	if cols != DefaultColors {
		t.Error("Mismatch in colors")
	}
}
func TestSetupExternalVao(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
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
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	newPrecision := 100
	square.SetPrecision(newPrecision)

	if square.precision != newPrecision {
		t.Error("Mismatch in the precision")
	}
}

func TestSetRotationAngle(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	angle := float32(2.0)
	square.SetAngle(angle)

	if square.angle != angle {
		t.Errorf("Mismatch in the angle instead of '%f', we have '%f'.", angle, square.angle)
	}
}
func TestSetRotationAxis(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	axis := mgl32.Vec3{0, 1, 0}
	square.SetAxis(axis)

	if square.axis != axis {
		t.Errorf("Mismatch in the axis")
	}
}
func TestGetNormal(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	normal := square.GetNormal()
	expectedNormal := mgl32.Vec3{0, 0, 1}
	if normal != expectedNormal {
		t.Error("Invalid normal vector.")
	}
	testData := []struct {
		coordinates    [4]mgl32.Vec3
		expectedNormal mgl32.Vec3
	}{
		{[4]mgl32.Vec3{mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 0, 0}, mgl32.Vec3{1, 1, 0}, mgl32.Vec3{0, 1, 0}}, mgl32.Vec3{0, 0, 1}},
		{[4]mgl32.Vec3{mgl32.Vec3{-0.5, -1.5, -0.5}, mgl32.Vec3{-0.5, -1.5, 0.5}, mgl32.Vec3{0.5, -1.5, 0.5}, mgl32.Vec3{0.5, -1.5, -0.5}}, mgl32.Vec3{0, 1, 0}},
	}

	for _, item := range testData {
		square = New(item.coordinates, DefaultColors, shader)
		if square.GetNormal() != item.expectedNormal {
			v1 := item.coordinates[1].Sub(item.coordinates[0])
			v2 := item.coordinates[3].Sub(item.coordinates[0])
			t.Log("v1")
			t.Log(v1)
			t.Log("v2")
			t.Log(v2)
			t.Log("v1.Cross(v2)")
			t.Log(v1.Cross(v2))
			t.Log("v1.Cross(v2).Normalize()")
			t.Log(v1.Cross(v2).Normalize())
			t.Log("(v1.Cross(v2)).Normalize()")
			t.Log((v1.Cross(v2)).Normalize())
			t.Log(square.GetNormal())
			t.Log(item.expectedNormal)
			t.Error("Invalid normals")
		}
	}
}
func TestDrawMode(t *testing.T) {
	shader.HasTextureValue = false
	square := New(DefaultCoordinates, DefaultColors, shader)
	if square.drawMode != 0 {
		t.Errorf("Invalid default draw mode. Instead of '0', we got '%d'", square.drawMode)
	}
	square.DrawMode(1)
	if square.drawMode != 1 {
		t.Errorf("Invalid  draw mode. Instead of '1', we got '%d'", square.drawMode)
	}
	square.DrawMode(2)
	if square.drawMode != 1 {
		t.Errorf("Invalid draw mode. Instead of '1', we got '%d'", square.drawMode)
	}
}
func TestDrawWithLight(t *testing.T) {
	square := New(DefaultCoordinates, DefaultColors, shader)
	square.DrawMode(DRAW_MODE_LIGHT)
	if len(square.vao.Get()) != 0 {
		t.Error("Vao is not empty before the first setup.")
	}
	square.DrawWithUniforms(mgl32.Ident4(), mgl32.Ident4())
	if len(square.vao.Get()) != 36 {
		t.Errorf("Vao should be 36 long. Instead of it, it's '%d'", len(square.vao.Get()))
	}
}
