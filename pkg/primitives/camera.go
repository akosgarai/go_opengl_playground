package primitives

// based on the opengl dev. cookbook solution

import (
	"math"
)

type Camera struct {
	// camera position
	position Vector
	// look at position (the camera is looking at lookat from it's position.)
	lookAt Vector
	// world up direction.
	UpDirection Vector

	fov         float64
	aspectRatio float64

	cameraUpDirection    Vector
	cameraRightDirection Vector
	// view matrix
	viewMatrix *Matrix4x4
	// projection matrix
	projectionMatrix *Matrix4x4

	freeCameraOptions struct {
		translation Vector
		yaw         float64
		pitch       float64
		roll        float64
	}
	targetCameraOptions struct {
		target      Vector
		yaw         float64
		pitch       float64
		roll        float64
		minRy       float64
		maxRy       float64
		distance    float64
		minDistance float64
		maxDistance float64
	}
}

func NewCamera() *Camera {
	var cam Camera
	cam.freeCameraOptions.translation = Vector{0, 0, 0}
	cam.targetCameraOptions.minRy = -60
	cam.targetCameraOptions.maxRy = 60
	cam.targetCameraOptions.minDistance = 0.1
	cam.targetCameraOptions.maxDistance = 1000
	cam.fov = 45
	cam.aspectRatio = 1
	return &cam
}

// GetAspectRatio returns the aspect ratio of the camera
func (c *Camera) GetAspectRatio() float64 {
	return c.aspectRatio
}

// GetFOV returns the fov of the camera
func (c *Camera) GetFOV() float64 {
	return c.fov
}

// GetViewMatrix returns the viewMatrix of the camera
func (c *Camera) GetViewMatrix() *Matrix4x4 {
	return c.viewMatrix
}

// GetProjectionMatrix returns the viewMatrix of the camera
func (c *Camera) GetProjectionMatrix() *Matrix4x4 {
	return c.projectionMatrix
}

// GetPosition returns the position of the camera
func (c *Camera) GetPosition() Vector {
	return c.position
}

// SetPosition updates the position of the camera
func (c *Camera) SetPosition(newPos Vector) {
	c.position = newPos
}

// SetupProjection creates the projection matrix and setups the fow and the aspect_ration
func (c *Camera) SetupProjection(fov, aspRatio float64) {
	c.fov = fov
	c.aspectRatio = aspRatio
	c.projectionMatrix = Perspective(float32(fov), float32(aspRatio), 0.1, 1000)
}

// FreeCameraRotate sets up the YPR values for the transformation
func (c *Camera) FreeCameraRotate(yaw, pitch, roll float64) {
	c.freeCameraOptions.yaw = yaw
	c.freeCameraOptions.pitch = pitch
	c.freeCameraOptions.roll = roll
}

// FreeCameraWalk updates the translation for the transformation (forward, back directions)
func (c *Camera) FreeCameraWalk(amount float64) {
	c.freeCameraOptions.translation = c.freeCameraOptions.translation.Add(
		c.lookAt.MultiplyScalar(amount))
}

// FreeCameraStrafe updates the translation for the transformation (left, right directions)
func (c *Camera) FreeCameraStrafe(amount float64) {
	c.freeCameraOptions.translation = c.freeCameraOptions.translation.Add(
		c.cameraRightDirection.MultiplyScalar(amount))
}

// FreeCameraLift updates the translation for the transformation (up, down directions)
func (c *Camera) FreeCameraLift(amount float64) {
	c.freeCameraOptions.translation = c.freeCameraOptions.translation.Add(
		c.cameraUpDirection.MultiplyScalar(amount))
}

// FreeCameraUpdate setup the viewMatrix based on the given freeCameraOptions
func (c *Camera) FreeCameraUpdate() {
	rotation := YPR(
		c.freeCameraOptions.yaw,
		c.freeCameraOptions.pitch,
		c.freeCameraOptions.roll)
	c.position = c.position.Add(c.freeCameraOptions.translation)
	c.freeCameraOptions.translation = Vector{0, 0, 0}

	c.lookAt = *rotation.MultiVector(Vector{0, 0, 1})
	target := c.position.Add(c.lookAt)
	c.cameraUpDirection = *rotation.MultiVector(Vector{0, 1, 0})
	c.cameraRightDirection = c.lookAt.Cross(c.cameraUpDirection)
	// glm::lookAt(position, tgt, up)
	c.viewMatrix = LookAt(c.position, target, c.cameraUpDirection)
}

// TargetCameraUpdate setup the viewMatrix based on the given targetCameraOptions
func (c *Camera) TargetCameraUpdate() {
	rotation := YPR(
		DegToRad(c.targetCameraOptions.yaw),
		DegToRad(c.targetCameraOptions.pitch),
		DegToRad(c.targetCameraOptions.roll))
	T := rotation.MultiVector(Vector{0, 0, c.targetCameraOptions.distance})
	c.position = c.targetCameraOptions.target.Add(*T)
	c.lookAt = (c.targetCameraOptions.target.Subtract(c.position)).Normalize()
	c.cameraUpDirection = *(rotation.MultiVector(c.UpDirection))
	c.cameraRightDirection = c.lookAt.Cross(c.cameraUpDirection)
	c.viewMatrix = LookAt(c.position, c.targetCameraOptions.target, c.cameraUpDirection)
}
func (c *Camera) TargetCameraSetTarget(target Vector) {
	c.targetCameraOptions.target = target
	// glm::distance(vector, vector) ??
	c.targetCameraOptions.distance = (c.position.Subtract(c.targetCameraOptions.target)).Length()
	c.targetCameraOptions.distance = math.Max(
		c.targetCameraOptions.minDistance, math.Min(
			c.targetCameraOptions.distance, c.targetCameraOptions.maxDistance))
	c.viewMatrix = LookAt(c.position, c.targetCameraOptions.target, c.cameraUpDirection)
	c.targetCameraOptions.yaw = 0
	c.targetCameraOptions.pitch = 0
	if c.viewMatrix.Points[0] < 0 {
		// m_yaw = glm::degrees((float)(M_PI - asinf(-V[2][0])) );
		c.targetCameraOptions.yaw = RadToDeg(math.Pi - math.Asin(-float64(c.viewMatrix.Points[8])))
	} else {
		// m_yaw = glm::degrees(asinf(-V[2][0]));
		c.targetCameraOptions.yaw = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[8])))
	}
	// m_pitch = glm::degrees(asinf(-V[1][2]));
	c.targetCameraOptions.pitch = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[6])))
}
func (c *Camera) TargetCameraGetTarget() Vector {
	return c.targetCameraOptions.target
}

// TargetCameraRotate sets up the YPR values for the transformation
func (c *Camera) TargetCameraRotate(yaw, pitch, roll float64) {
	c.targetCameraOptions.yaw = c.targetCameraOptions.yaw + yaw
	c.targetCameraOptions.pitch = c.targetCameraOptions.pitch + pitch
	c.targetCameraOptions.pitch = math.Min(
		math.Max(c.targetCameraOptions.pitch, c.targetCameraOptions.minRy),
		c.targetCameraOptions.maxRy)
	c.TargetCameraUpdate()
}
func (c *Camera) TargetCameraPan(dX, dY float64) {
	x := c.cameraRightDirection.MultiplyScalar(dX)
	y := c.cameraUpDirection.MultiplyScalar(dY)
	c.position = c.position.Add(x.Add(y))
	c.targetCameraOptions.target = c.targetCameraOptions.target.Add(x.Add(y))
	c.TargetCameraUpdate()
}
func (c *Camera) TargetCameraZoom(amount float64) {
	c.position = c.position.Add(c.lookAt.MultiplyScalar(amount))
	c.targetCameraOptions.distance = (c.position.Subtract(c.targetCameraOptions.target)).Length()
	c.targetCameraOptions.distance = math.Max(
		c.targetCameraOptions.minDistance,
		math.Min(
			c.targetCameraOptions.distance,
			c.targetCameraOptions.maxDistance))
	c.TargetCameraUpdate()
}
func (c *Camera) TargetCameraMove(dX, dY float64) {
	x := c.cameraRightDirection.MultiplyScalar(dX)
	y := c.lookAt.MultiplyScalar(dY)
	c.position = c.position.Add(x.Add(y))
	c.targetCameraOptions.target = c.targetCameraOptions.target.Add(x.Add(y))
	c.TargetCameraUpdate()
}
