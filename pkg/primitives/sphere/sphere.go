package sphere

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
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
	DrawTriangles(int32)
	Close(int)
	VertexAttribPointer(uint32, int32, int32, int)
	BindVertexArray()
	BindBufferData([]float32)
}

type Sphere struct {
	precision int
	vao       *vao.VAO
	shader    Shader

	center mgl32.Vec3
	radius float32
	color  mgl32.Vec3

	direction mgl32.Vec3
	speed     float32
	// rotation parameters
	// angle has to be in radian
	angle float32
	axis  mgl32.Vec3

	material *material.Material
	drawMode int
}

func New(center, color mgl32.Vec3, radius float32, shader Shader) *Sphere {
	return &Sphere{
		precision: 10,
		vao:       vao.NewVAO(),
		shader:    shader,

		center: center,
		radius: radius,
		color:  color,

		direction: mgl32.Vec3{0, 0, 0},
		speed:     0,

		angle:    float32(0.0),
		axis:     mgl32.Vec3{0, 0, 0},
		material: material.New(color, color, color, 36.0),
		drawMode: DRAW_MODE_COLOR,
	}
}

// SetRadius updates the radius of the sphere
func (s *Sphere) Log() string {
	logString := "Sphere:\n"
	logString += " - Center : Coordinate: Vector{" + trans.Vec3ToString(s.center) + "}, radius: " + trans.Float32ToString(s.radius) + ", color: Vector{" + trans.Vec3ToString(s.color) + "}\n"
	logString += " - Movement : Direction: Vector{" + trans.Vec3ToString(s.direction) + "}, speed: " + trans.Float32ToString(s.speed) + "\n"
	logString += " - Rotation : Axis: Vector{" + trans.Vec3ToString(s.axis) + "}, angle: " + trans.Float32ToString(s.angle) + "\n"
	logString += s.material.Log() + "\n"
	return logString
}

// SetPrecision updates the precision of the rectangle
func (s *Sphere) SetPrecision(p int) {
	s.precision = p
}

// SetCenter updates the center of the sphere
func (s *Sphere) SetCenter(c mgl32.Vec3) {
	s.center = c
}

// GetCenter returns the center of the sphere
func (s *Sphere) GetCenter() mgl32.Vec3 {
	return s.center
}

// SetColor updates the color of the sphere
func (s *Sphere) SetColor(c mgl32.Vec3) {
	s.color = c
}

// GetColor returns the color of the sphere
func (s *Sphere) GetColor() mgl32.Vec3 {
	return s.color
}

// SetRadius updates the radius of the sphere
func (s *Sphere) SetRadius(r float32) {
	s.radius = r
}

// GetRadius returns the radius of the sphere
func (s *Sphere) GetRadius() float32 {
	return s.radius
}

// GetDirection returns the direction vector.
func (s *Sphere) GetDirection() mgl32.Vec3 {
	return s.direction
}

// SetDirection updates the direction vector.
func (s *Sphere) SetDirection(dir mgl32.Vec3) {
	s.direction = dir
}

// SetIndexDirection updates the direction vector.
func (s *Sphere) SetIndexDirection(index int, value float32) {
	s.direction[index] = value
}

// SetSpeed updates the speed.
func (s *Sphere) SetSpeed(speed float32) {
	s.speed = speed
}

// GetAngle returns the angle of the sphere
func (s *Sphere) GetAngle() float32 {
	return s.angle
}

// GetAxis returns the axis vector.
func (s *Sphere) GetAxis() mgl32.Vec3 {
	return s.axis
}

// SetAngle updates the angle of the sphere
func (s *Sphere) SetAngle(angle float32) {
	s.angle = angle
}

// SetAxis updatess the axis vector.
func (s *Sphere) SetAxis(axis mgl32.Vec3) {
	s.axis = axis
}

// SetMaterial updates the material of the sphere.
func (s *Sphere) SetMaterial(mat *material.Material) {
	s.material = mat
}

// DrawMode updates the draw mode after validation. If it fails, it keeps the original value.
func (s *Sphere) DrawMode(mode int) {
	if mode != DRAW_MODE_COLOR && mode != DRAW_MODE_LIGHT {
		return
	}
	s.drawMode = mode
}
func (s *Sphere) triangleToVao(pa, pb, pc mgl32.Vec3) {
	s.vao.AppendVectors(pa, s.color)
	s.vao.AppendVectors(pb, s.color)
	s.vao.AppendVectors(pc, s.color)
}
func (s *Sphere) setupVao() {
	// the coordinates will be set as a following: origo as center, 1 as radius, for drawing, the translation and scale could be done later in the model transformation.
	// Sphere top: center + v{0,radius,0}, bottom: center + v{0,-radius,0}, left: center + v{-radius,0,0}, right: center + v{radius,0,0}
	// Idea : start drawing triangles from both direction (top, bottom). step the coordinates and calculate the triangles, add them to vao.
	// - step for y coord, : radius * 2 / numOfRows
	RefPoint := mgl32.Vec3{0, 1, 0}
	step := -mgl32.DegToRad(float32(360.0) / float32(s.precision))
	for i := 0; i < s.precision; i++ {
		i_Rotation := mgl32.HomogRotate3DZ(float32(i) * step)
		i1_Rotation := mgl32.HomogRotate3DZ(float32(i+1) * step)
		for j := 0; j < s.precision; j++ {
			j1_Rotation := mgl32.HomogRotate3DY(float32(j+1) * step)
			j_Rotation := mgl32.HomogRotate3DY(float32(j) * step)
			if i == 0 {
				p2 := mgl32.TransformCoordinate(RefPoint, j_Rotation.Mul4(i1_Rotation))
				p3 := mgl32.TransformCoordinate(RefPoint, j1_Rotation.Mul4(i1_Rotation))
				s.triangleToVao(RefPoint, p2, p3)
			} else {
				p1 := mgl32.TransformCoordinate(RefPoint, j_Rotation.Mul4(i_Rotation))
				p2 := mgl32.TransformCoordinate(RefPoint, j1_Rotation.Mul4(i_Rotation))
				p3 := mgl32.TransformCoordinate(RefPoint, j1_Rotation.Mul4(i1_Rotation))
				p4 := mgl32.TransformCoordinate(RefPoint, j_Rotation.Mul4(i1_Rotation))
				s.triangleToVao(p1, p2, p3)
				s.triangleToVao(p1, p3, p4)
			}
		}
	}
}
func (s *Sphere) buildVao() {
	s.setupVao()

	s.shader.BindBufferData(s.vao.Get())

	s.shader.BindVertexArray()
	// setup points
	s.shader.VertexAttribPointer(0, 3, 4*6, 0)
	// setup color
	s.shader.VertexAttribPointer(1, 3, 4*6, 4*3)

}
func (s *Sphere) modelTransformation() mgl32.Mat4 {
	return mgl32.Translate3D(
		s.center.X(),
		s.center.Y(),
		s.center.Z()).Mul4(mgl32.HomogRotate3D(s.angle, s.axis)).Mul4(mgl32.Scale3D(
		s.radius,
		s.radius,
		s.radius))
}
func (s *Sphere) DrawWithUniforms(view, projection mgl32.Mat4) {
	s.shader.Use()
	s.shader.SetUniformMat4("view", view)
	s.shader.SetUniformMat4("projection", projection)
	M := s.modelTransformation()

	s.shader.SetUniformMat4("model", M)
	s.draw()
}

func (s *Sphere) Draw() {
	s.shader.Use()
	s.draw()
}
func (s *Sphere) draw() {
	s.buildVao()
	s.shader.DrawTriangles(int32(len(s.vao.Get()) / 6))
	s.shader.Close(1)
}
func (s *Sphere) Update(dt float64) {
	delta := float32(dt)
	motionVector := s.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * s.speed)
	}
	s.center = (s.center).Add(motionVector)
}
