package bug

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/primitives/material"
	"github.com/akosgarai/opengl_playground/pkg/primitives/sphere"
	trans "github.com/akosgarai/opengl_playground/pkg/primitives/transformations"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Bug struct {
	vao    *vao.VAO
	shader sphere.Shader

	// the initial position of the bug. The shapes are calculated based on this.
	position mgl32.Vec3
	size     float32
	shapes   [4]*sphere.Sphere

	direction mgl32.Vec3
	speed     float32
	// rotation parameters
	// angle has to be in radian
	angle float32
	axis  mgl32.Vec3
}

// SetRadius updates the radius of the sphere
func (b *Bug) Log() string {
	logString := "Bug:\n"
	logString += " - Position : Coordinate: Vector{" + trans.Vec3ToString(b.position) + "}\n"
	logString += " - Movement : Direction: Vector{" + trans.Vec3ToString(b.direction) + "}, speed: " + trans.Float32ToString(b.speed) + "\n"
	logString += " - Rotation : Axis: Vector{" + trans.Vec3ToString(b.axis) + "}, angle: " + trans.Float32ToString(b.angle) + "\n"
	for i := 0; i < 4; i++ {
		logString += b.shapes[i].Log()
	}
	return logString
}

// SetPosition updates the position of the bug
func (b *Bug) SetCenter(p mgl32.Vec3) {
	b.position = p
}

// GetCenterPoint returns the position of the Bug
func (b *Bug) GetCenterPoint() mgl32.Vec3 {
	return b.position
}

// GetDirection returns the direction vector.
func (b *Bug) GetDirection() mgl32.Vec3 {
	return b.direction
}

// SetDirection updates the direction vector.
func (b *Bug) SetDirection(dir mgl32.Vec3) {
	b.direction = dir
	for i, _ := range b.shapes {
		b.shapes[i].SetDirection(dir)
	}
}

// SetSpeed updates the speed.
func (b *Bug) SetSpeed(speed float32) {
	b.speed = speed
	for i, _ := range b.shapes {
		b.shapes[i].SetSpeed(speed)
	}
}

// GetAngle returns the angle of the bug
func (b *Bug) GetAngle() float32 {
	return b.angle
}

// GetAxis returns the axis vector.
func (b *Bug) GetAxis() mgl32.Vec3 {
	return b.axis
}

// SetAngle updates the angle of the Bug
func (b *Bug) SetAngle(angle float32) {
	b.angle = angle
	for i, _ := range b.shapes {
		b.shapes[i].SetAngle(angle)
	}
}

// SetAxis updatess the axis vector.
func (b *Bug) SetAxis(axis mgl32.Vec3) {
	b.axis = axis
	for i, _ := range b.shapes {
		b.shapes[i].SetAxis(axis)
	}
}

// Draw calls the Draw function of the shapes.
func (b *Bug) Draw() {
	for i, _ := range b.shapes {
		b.shapes[i].Draw()
	}
}

// DrawWithUniforms calls the DrawWithUniforms function of the shapes.
func (b *Bug) DrawWithUniforms(view, projection mgl32.Mat4) {
	for i, _ := range b.shapes {
		b.shapes[i].DrawWithUniforms(view, projection)
	}
}

// DrawMode updates the draw mode after validation. Currently it only supports the `DRAW_MODE_LIGHT`.
func (b *Bug) DrawMode(mode int) {
	return
}

// Update calls the Update function of the shapes.
func (b *Bug) Update(dt float64) {
	for i, _ := range b.shapes {
		b.shapes[i].Update(dt)
	}
}

func Firefly(position mgl32.Vec3, size float32, materials [3]*material.Material, shaderProgram sphere.Shader) *Bug {
	FakeColor := mgl32.Vec3{1, 1, 1}
	bottom := sphere.New(position, FakeColor, size, shaderProgram)
	bottom.SetMaterial(materials[0])
	bottom.SetPrecision(15)
	bottom.DrawMode(sphere.DRAW_MODE_LIGHT)

	body := sphere.New(mgl32.Vec3{position.X(), position.Y(), position.Z() + size*2}, FakeColor, size*2, shaderProgram)
	body.SetMaterial(materials[1])
	body.SetPrecision(15)
	body.DrawMode(sphere.DRAW_MODE_LIGHT)

	leftEye := sphere.New(mgl32.Vec3{position.X() + size, position.Y(), position.Z() + size*3.5}, FakeColor, size/2, shaderProgram)
	leftEye.SetMaterial(materials[2])
	leftEye.SetPrecision(15)
	leftEye.DrawMode(sphere.DRAW_MODE_LIGHT)

	rightEye := sphere.New(mgl32.Vec3{position.X() - size, position.Y(), position.Z() + size*3.5}, FakeColor, size/2, shaderProgram)
	rightEye.SetMaterial(material.Ruby)
	rightEye.SetPrecision(15)
	rightEye.DrawMode(sphere.DRAW_MODE_LIGHT)

	shapes := [4]*sphere.Sphere{bottom, body, leftEye, rightEye}

	return &Bug{
		vao:    vao.NewVAO(),
		shader: shaderProgram,

		position: position,
		size:     size,
		shapes:   shapes,

		direction: mgl32.Vec3{0, 0, 0},
		speed:     float32(0.0),

		angle: float32(0.0),
		axis:  mgl32.Vec3{0, 0, 0},
	}
}
