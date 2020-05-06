package rectangle

import (
	"github.com/go-gl/mathgl/mgl32"

	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

const (
	DRAW_MODE_COLOR = 0
	DRAW_MODE_LIGHT = 1
)

type Shader interface {
	Use()
	SetUniformMat4(string, mgl32.Mat4)
	SetUniform3f(string, float32, float32, float32)
	SetUniform1f(string, float32)
	DrawTriangles(int32)
	Close(int)
	VertexAttribPointer(uint32, int32, int32, int)
	BindVertexArray()
	BindBufferData([]float32)
	HasTexture() bool
}

type Rectangle struct {
	precision int
	vao       *vao.VAO
	shader    Shader

	colors [4]mgl32.Vec3
	points [4]mgl32.Vec3

	// movement paramteres
	direction mgl32.Vec3
	speed     float32
	// rotation parameters
	// angle has to be in radian
	angle float32
	axis  mgl32.Vec3
	// for the lighting stuff we have to distinguse the need of the color or the normal vector.
	// it also means, that we need to know the normal vector of the rectangle.-> GetNormal() function.
	// drawMode map:
	// 0 - normal draw with colors
	// 1 - draw with normal vectors.
	drawMode int
}

func New(points, color [4]mgl32.Vec3, shader Shader) *Rectangle {
	return &Rectangle{
		precision: 1,
		vao:       vao.NewVAO(),
		shader:    shader,
		colors:    color,
		points:    points,
		direction: mgl32.Vec3{0, 0, 0},
		speed:     0,
		angle:     0,
		axis:      mgl32.Vec3{0, 0, 0},
		drawMode:  DRAW_MODE_COLOR,
	}
}

// NewSquare returns a square. The inputs are 2 diagonal points of the square, it's normal vector, and the color vector.
// IT also needs a shader as input.
// Let's define the points of the square as 'A', 'B', 'C', 'D'.
// Let's imagine, that we have 'A' and 'C' as inputs. (or we have 'B', and 'D'. these are the only options for the diagonals)
// We can calulate 'B' and 'D' by rotating the other points with 90 deg on the normal axis.
func NewSquare(p1, p3, normal, color mgl32.Vec3, shader Shader) *Rectangle {
	rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(90.0), normal)
	p2 := mgl32.TransformCoordinate(p1, rotationMatrix)
	p0 := mgl32.TransformCoordinate(p3, rotationMatrix)
	points := [4]mgl32.Vec3{p0, p1, p2, p3}
	colors := [4]mgl32.Vec3{color, color, color, color}
	return New(points, colors, shader)
}

// Coordinates returns the points of the rectangle. It's necessary for the cuboid construction.
func (r *Rectangle) Coordinates() [4]mgl32.Vec3 {
	return r.points
}

// Colors returns the colors of the rectangle. It's necessary for the cuboid construction.
func (r *Rectangle) Colors() [4]mgl32.Vec3 {
	return r.colors
}

// Log returns the string representation of this object.
func (r *Rectangle) Log() string {
	logString := "Rectangle:\n"
	logString += " - A : Coordinate: Vector{" + trans.Vec3ToString(r.points[0]) + "}, color: Vector{" + trans.Vec3ToString(r.colors[0]) + "}\n"
	logString += " - B : Coordinate: Vector{" + trans.Vec3ToString(r.points[1]) + "}, color: Vector{" + trans.Vec3ToString(r.colors[1]) + "}\n"
	logString += " - C : Coordinate: Vector{" + trans.Vec3ToString(r.points[2]) + "}, color: Vector{" + trans.Vec3ToString(r.colors[2]) + "}\n"
	logString += " - D : Coordinate: Vector{" + trans.Vec3ToString(r.points[3]) + "}, color: Vector{" + trans.Vec3ToString(r.colors[3]) + "}\n"
	logString += " - Normal : Coordinate: Vector{" + trans.Vec3ToString(r.axis) + "}, DrawMode: " + trans.IntegerToString(r.drawMode) + "\n"
	logString += " - Movement : Direction: Vector{" + trans.Vec3ToString(r.direction) + "}, speed: " + trans.Float32ToString(r.speed) + "}\n"
	logString += " - Rotation : Axis: Vector{" + trans.Vec3ToString(r.axis) + "}, angle: " + trans.Float32ToString(r.angle) + "}\n"
	return logString
}

// SetColor updates every color with the given one.
func (r *Rectangle) SetColor(color mgl32.Vec3) {
	for i := 0; i < 4; i++ {
		r.colors[i] = color
	}
}

// SetIndexColor updates the color of the given index.
func (r *Rectangle) SetIndexColor(index int, color mgl32.Vec3) {
	r.colors[index] = color
}

// SetDirection updates the direction vector.
func (r *Rectangle) SetDirection(dir mgl32.Vec3) {
	r.direction = dir
}

// SetIndexDirection updates the direction vector.
func (r *Rectangle) SetIndexDirection(index int, value float32) {
	r.direction[index] = value
}

// SetSpeed updates the speed.
func (r *Rectangle) SetSpeed(speed float32) {
	r.speed = speed
}

// SetPrecision updates the precision of the rectangle
func (r *Rectangle) SetPrecision(p int) {
	r.precision = p
}

// SetAngle updates the angle.
// Input has to be radian.
func (r *Rectangle) SetAngle(angle float32) {
	r.angle = angle
}

// SetAxis updates the axis.
func (r *Rectangle) SetAxis(axis mgl32.Vec3) {
	r.axis = axis
}

// GetDirection returns the direction of the rectangle
func (r *Rectangle) GetDirection() mgl32.Vec3 {
	return r.direction
}
func (r *Rectangle) appendRectangleToVao(coordinates, data2 [4]mgl32.Vec3) {
	indicies := [6]int{0, 1, 2, 0, 2, 3}
	if r.shader.HasTexture() {
		textureCoords := [4]mgl32.Vec2{
			{0.0, 1.0},
			{1.0, 1.0},
			{1.0, 0.0},
			{0.0, 0.0},
		}
		for i := 0; i < 6; i++ {
			r.vao.AppendTextureVectors(coordinates[indicies[i]], data2[indicies[i]], textureCoords[indicies[i]])
		}
	} else {
		for i := 0; i < 6; i++ {
			r.vao.AppendVectors(coordinates[indicies[i]], data2[indicies[i]])
		}
	}
}

func (r *Rectangle) insertEverythingToVao() {
	verticalStep := (r.points[1].Sub(r.points[0])).Mul(1.0 / float32(r.precision))
	horisontalStep := (r.points[3].Sub(r.points[0])).Mul(1.0 / float32(r.precision))

	for horisontalLoopIndex := 0; horisontalLoopIndex < r.precision; horisontalLoopIndex++ {
		for verticalLoopIndex := 0; verticalLoopIndex < r.precision; verticalLoopIndex++ {
			a := r.points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			b := r.points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			c := r.points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			d := r.points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			if r.drawMode == DRAW_MODE_COLOR {
				r.appendRectangleToVao([4]mgl32.Vec3{a, b, c, d}, r.colors)
			} else if r.drawMode == DRAW_MODE_LIGHT {
				normals := [4]mgl32.Vec3{r.GetNormal(), r.GetNormal(), r.GetNormal(), r.GetNormal()}
				r.appendRectangleToVao([4]mgl32.Vec3{a, b, c, d}, normals)
			}
		}
	}
}
func (r *Rectangle) setupVao() {
	r.vao.Clear()
	r.insertEverythingToVao()
}
func (r *Rectangle) SetupExternalVao(v *vao.VAO) *vao.VAO {
	tmpVao := r.vao
	r.vao = v
	r.insertEverythingToVao()
	v = r.vao
	r.vao = tmpVao
	return v
}

func (r *Rectangle) buildVaoWithTexture() {
	// Create the vao object
	r.setupVao()
	r.shader.BindBufferData(r.vao.Get())
	r.shader.BindVertexArray()
	// setup points
	r.shader.VertexAttribPointer(0, 3, 4*8, 0)
	// setup color
	r.shader.VertexAttribPointer(1, 3, 4*8, 4*3)
	// setup texture
	r.shader.VertexAttribPointer(2, 2, 4*8, 4*6)
}
func (r *Rectangle) buildVaoWithoutTexture() {
	// Create the vao object
	r.setupVao()

	r.shader.BindBufferData(r.vao.Get())

	r.shader.BindVertexArray()
	// setup points
	r.shader.VertexAttribPointer(0, 3, 4*6, 0)
	// setup color
	r.shader.VertexAttribPointer(1, 3, 4*6, 4*3)
}

// Draw is for drawing the rectangle to the screen.
func (r *Rectangle) Draw() {
	r.shader.Use()
	if !r.shader.HasTexture() {
		r.shader.SetUniformMat4("MVP", r.modelTransformation())
		r.drawWithoutTextures()
	} else {
		r.drawWithTextures()
	}
}
func (r *Rectangle) drawWithTextures() {
	r.buildVaoWithTexture()
	r.shader.DrawTriangles(int32(len(r.vao.Get()) / 8))
	r.shader.Close(2)
}
func (r *Rectangle) drawWithoutTextures() {
	r.buildVaoWithoutTexture()
	r.shader.DrawTriangles(int32(len(r.vao.Get()) / 6))
	r.shader.Close(1)
}
func (r *Rectangle) modelTransformation() mgl32.Mat4 {
	return mgl32.HomogRotate3D(r.angle, r.axis)
}

// DrawWithUniforms is for drawing the rectangle to the screen. It setups the
func (r *Rectangle) DrawWithUniforms(view, projection mgl32.Mat4) {
	r.shader.Use()
	r.shader.SetUniformMat4("view", view)
	r.shader.SetUniformMat4("projection", projection)
	r.shader.SetUniformMat4("model", r.modelTransformation())
	if !r.shader.HasTexture() {
		r.drawWithoutTextures()
	} else {
		r.drawWithTextures()
	}
}
func (r *Rectangle) Update(dt float64) {
	delta := float32(dt)
	motionVector := r.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * r.speed)
	}
	for i := 0; i < 4; i++ {
		r.points[i] = (r.points[i]).Add(motionVector)
	}
}

// GetNormal returns the normal vector of the square.
func (r *Rectangle) GetNormal() mgl32.Vec3 {
	v1 := r.points[1].Sub(r.points[0])
	v2 := r.points[3].Sub(r.points[0])
	return v1.Cross(v2).Normalize()
}

// DrawMode updates the draw mode after validation. If it fails, it keeps the original value.
func (s *Rectangle) DrawMode(mode int) {
	if mode == DRAW_MODE_COLOR || mode == DRAW_MODE_LIGHT {
		s.drawMode = mode
	}
}
