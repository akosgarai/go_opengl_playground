package primitives

import (
	"strconv"
)

type CameraImpr struct {
	// camera position
	position Vector
	// look at position (the camera is looking at lookat from it's position.)
	lookAt Vector
	// world up direction.
	UpDirection Vector

	projectionOptions struct {
		fov         float64
		aspectRatio float64

		far  float64
		near float64
	}
	// projection matrix
	projectionMatrix *Matrix4x4

	// camera coordinate system directions
	cameraUpDirection    Vector
	cameraRightDirection Vector
	cameraFrontDirection Vector
	// view matrix
	viewMatrix *Matrix4x4

	cameraOptions struct {
		// Yaw, pitch, roll are stored as radiants.
		yaw         float64
		pitch       float64
		roll        float64
		translation Vector
	}
}

// Log returns the string representation of this object.
func (c *CameraImpr) Log() string {
	logString := "Position: Vector{" + c.position.ToString() + "}\n"
	logString += "LookAt: Vector{" + c.lookAt.ToString() + "}\n"
	logString += "UpDirection: Vector{" + c.UpDirection.ToString() + "}\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + strconv.FormatFloat(c.projectionOptions.fov, 'f', 6, 64) + "\n"
	logString += " - aspectRatio : " + strconv.FormatFloat(c.projectionOptions.aspectRatio, 'f', 6, 64) + "\n"
	logString += " - far : " + strconv.FormatFloat(c.projectionOptions.far, 'f', 6, 64) + "\n"
	logString += " - near : " + strconv.FormatFloat(c.projectionOptions.near, 'f', 6, 64) + "\n"
	logString += "CameraUpDirection: Vector{" + c.cameraUpDirection.ToString() + "}\n"
	logString += "CameraRightDirection: Vector{" + c.cameraRightDirection.ToString() + "}\n"
	logString += "CameraForwardDirection: Vector{" + c.cameraFrontDirection.ToString() + "}\n"
	logString += "CameraOptions:\n"
	logString += " - yaw : " + strconv.FormatFloat(c.cameraOptions.yaw, 'f', 6, 64) + "\n"
	logString += " - pitch : " + strconv.FormatFloat(c.cameraOptions.pitch, 'f', 6, 64) + "\n"
	logString += " - roll : " + strconv.FormatFloat(c.cameraOptions.roll, 'f', 6, 64) + "\n"
	logString += " - translation: Vector{" + c.cameraOptions.translation.ToString() + "}\n"
	return logString
}

func NewCameraImpr() *CameraImpr {
	var cam CameraImpr
	cam.projectionOptions.fov = 45
	cam.projectionOptions.aspectRatio = 1
	cam.projectionOptions.near = 0.1
	cam.projectionOptions.far = 1000
	cam.cameraOptions.translation = Vector{0, 0, 0}
	return &cam
}

// GetViewMatrix returns the viewMatrix of the camera
func (c *CameraImpr) GetViewMatrix() *Matrix4x4 {
	return c.viewMatrix
}

// GetProjectionMatrix returns the viewMatrix of the camera
func (c *CameraImpr) GetProjectionMatrix() *Matrix4x4 {
	return c.projectionMatrix
}

// SetPosition updates the position of the camera
func (c *CameraImpr) SetPosition(newPos Vector) {
	c.position = newPos
}
func (c *CameraImpr) GetPosition() Vector {
	return c.position
}

// SetTranslation updates the cameraOptions.translation of the camera
func (c *CameraImpr) SetTranslation(newPos Vector) {
	c.cameraOptions.translation = newPos
}
func (c *CameraImpr) GetTranslation() Vector {
	return c.cameraOptions.translation
}

// SetupProjection creates the projection matrix and setups the fow and the aspect_ration
func (c *CameraImpr) SetupProjection(fov, aspRatio float64) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionMatrix = Perspective(float32(c.projectionOptions.fov), float32(c.projectionOptions.aspectRatio), float32(c.projectionOptions.near), float32(c.projectionOptions.far))
}

// Update setup the viewMatrix based on the given cameraOptions
func (c *CameraImpr) Update() {
	rotation := YPR(
		c.cameraOptions.yaw,
		c.cameraOptions.pitch,
		c.cameraOptions.roll)
	c.position = c.position.Add(c.cameraOptions.translation)
	c.cameraOptions.translation = Vector{0, 0, 0}
	// camera has to look to -z dir.
	c.cameraFrontDirection = *(rotation.MultiVector(Vector{0, 0, -1}))
	// worldUp * front.
	c.cameraRightDirection = c.cameraFrontDirection.Cross(c.UpDirection)
	c.cameraUpDirection = c.cameraFrontDirection.Cross(c.cameraRightDirection)
	c.viewMatrix = LookAtV(c.cameraRightDirection, c.cameraUpDirection, c.cameraFrontDirection, c.position)
}

// Walk updates the translation for the transformation (forward, back directions)
func (c *CameraImpr) Walk(amount float64) {
	c.cameraOptions.translation = c.cameraOptions.translation.Add(
		c.cameraFrontDirection.MultiplyScalar(amount))
	c.Update()
}

// Strafe updates the translation for the transformation (left, right directions)
func (c *CameraImpr) Strafe(amount float64) {
	c.cameraOptions.translation = c.cameraOptions.translation.Add(
		c.cameraRightDirection.MultiplyScalar(amount))
	c.Update()
}

// Lift updates the translation for the transformation (up, down directions)
func (c *CameraImpr) Lift(amount float64) {
	c.cameraOptions.translation = c.cameraOptions.translation.Add(
		c.cameraUpDirection.MultiplyScalar(amount))
	c.Update()
}

// Rotate updates the ypr values.
func (c *CameraImpr) Rotate(y, p, r float64) {
	c.cameraOptions.yaw = DegToRad(y)
	c.cameraOptions.pitch = DegToRad(p)
	c.cameraOptions.roll = DegToRad(r)
	c.Update()
}
