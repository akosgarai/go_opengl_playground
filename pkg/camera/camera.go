package camera

import (
	M "github.com/akosgarai/opengl_playground/pkg/matrix"
	V "github.com/akosgarai/opengl_playground/pkg/vector"
)

type Camera struct {
	Position    V.Vector
	LookAt      V.Vector
	UpDirection V.Vector
	xAxis       V.Vector
	yAxis       V.Vector
	zAxis       V.Vector
}

// New returns a camera instance with the initial setup.
func New(pos, lookAt, upDir V.Vector) *Camera {
	c := &Camera{}
	c.Position = pos
	c.LookAt = lookAt
	c.UpDirection = upDir
	c.setZAxis()
	c.setXAxis()
	c.setYAxis()
	return c
}
func (c *Camera) setZAxis() {
	c.zAxis = (c.Position.Add(c.LookAt.MultiplyScalar(-1))).Normalize()
}
func (c *Camera) setXAxis() {
	c.xAxis = ((c.UpDirection.Normalize()).Cross(c.zAxis)).Normalize()
}
func (c *Camera) setYAxis() {
	c.yAxis = c.zAxis.Cross(c.xAxis)
}

// GetTransformation returns the calculated camera transformation in M.Matrix format
func (c *Camera) GetTransformation() *M.Matrix {

	cameraTranslationMatrix := M.Translation(c.Position.MultiplyScalar(-1))
	var cameraRotationMatrix M.Matrix
	cameraTranslationMatrix.LoadIdentity()
	cameraRotationMatrix.Points[0] = c.xAxis.X
	cameraRotationMatrix.Points[4] = c.xAxis.Y
	cameraRotationMatrix.Points[8] = c.xAxis.Z
	cameraRotationMatrix.Points[1] = c.yAxis.X
	cameraRotationMatrix.Points[5] = c.yAxis.Y
	cameraRotationMatrix.Points[9] = c.yAxis.Z
	cameraRotationMatrix.Points[2] = c.zAxis.X
	cameraRotationMatrix.Points[6] = c.zAxis.Y
	cameraRotationMatrix.Points[10] = c.zAxis.Z
	return cameraRotationMatrix.Dot(*cameraTranslationMatrix)
}
