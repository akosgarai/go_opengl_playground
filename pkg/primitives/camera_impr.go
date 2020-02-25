package primitives

import (
	"math"
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

	cameraUpDirection    Vector
	cameraRightDirection Vector
	// view matrix
	viewMatrix *Matrix4x4

	cameraOptions struct {
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

func (c *CameraImpr) Log() string {
	logString := "Position: Vector{" + c.vectorToString(c.position) + "}\n"
	logString += "LookAt: Vector{" + c.vectorToString(c.lookAt) + "}\n"
	logString += "UpDirection: Vector{" + c.vectorToString(c.UpDirection) + "}\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + strconv.FormatFloat(c.projectionOptions.fov, 'f', 6, 64)
	logString += " - aspectRatio : " + strconv.FormatFloat(c.projectionOptions.aspectRatio, 'f', 6, 64)
	logString += " - far : " + strconv.FormatFloat(c.projectionOptions.far, 'f', 6, 64)
	logString += " - near : " + strconv.FormatFloat(c.projectionOptions.near, 'f', 6, 64)
	return logString
}
func (c *CameraImpr) vectorToString(v Vector) string {
	x := strconv.FormatFloat(v.X, 'f', 6, 64)
	y := strconv.FormatFloat(v.Y, 'f', 6, 64)
	z := strconv.FormatFloat(v.Z, 'f', 6, 64)
	return "X : " + x + ", Y : " + y + ", Z : " + z

}

func NewCameraImpr() *CameraImpr {
	var cam CameraImpr
	cam.cameraOptions.minRy = -60
	cam.cameraOptions.maxRy = 60
	cam.cameraOptions.minDistance = 0.1
	cam.cameraOptions.maxDistance = 1000
	cam.projectionOptions.fov = 45
	cam.projectionOptions.aspectRatio = 1
	cam.projectionOptions.near = 0.1
	cam.projectionOptions.far = 1000
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

// SetupProjection creates the projection matrix and setups the fow and the aspect_ration
func (c *CameraImpr) SetupProjection(fov, aspRatio float64) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionMatrix = Perspective(float32(c.projectionOptions.fov), float32(c.projectionOptions.aspectRatio), float32(c.projectionOptions.near), float32(c.projectionOptions.far))
}

// TargetCameraSetTarget updates the camera based on the new targetPoint
// It updates the lookAt vector (target - position normalized)
// It calculates the distance
// Sets the viewMatrix
// Sets the yaw, pitch values
func (c *CameraImpr) TargetCameraSetTarget(target Vector) {
	c.lookAt = (target.Subtract(c.position)).Normalize()
	c.cameraOptions.distance = (c.position.Subtract(target)).Length()
	c.cameraOptions.distance = math.Max(
		c.cameraOptions.minDistance, math.Min(
			c.cameraOptions.distance, c.cameraOptions.maxDistance))
	c.viewMatrix = LookAt(c.position, c.lookAt, c.cameraUpDirection)
	c.cameraOptions.yaw = 0
	c.cameraOptions.pitch = 0
	if c.viewMatrix.Points[0] < 0 {
		c.cameraOptions.yaw = RadToDeg(math.Pi - math.Asin(-float64(c.viewMatrix.Points[8])))
	} else {
		c.cameraOptions.yaw = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[8])))
	}
	c.cameraOptions.pitch = RadToDeg(math.Asin(-float64(c.viewMatrix.Points[6])))
}

// TargetCameraUpdate setup the viewMatrix based on the given cameraOptions
func (c *CameraImpr) TargetCameraUpdate() {
	rotation := YPR(
		DegToRad(c.cameraOptions.yaw),
		DegToRad(c.cameraOptions.pitch),
		DegToRad(c.cameraOptions.roll))
	T := rotation.MultiVector(Vector{0, 0, c.cameraOptions.distance})
	c.position = c.lookAt.Add(*T)
	c.lookAt = (c.lookAt.Subtract(c.position)).Normalize()
	c.cameraUpDirection = *(rotation.MultiVector(c.UpDirection))
	c.cameraRightDirection = c.lookAt.Cross(c.cameraUpDirection)
	c.viewMatrix = LookAt(c.position, c.lookAt, c.cameraUpDirection)
}
func (c *CameraImpr) TargetCameraMove(dX, dY float64) {
	x := c.cameraRightDirection.MultiplyScalar(dX)
	y := c.lookAt.MultiplyScalar(dY)
	c.position = c.position.Add(x.Add(y))
	c.lookAt = c.lookAt.Add(x.Add(y))
	c.TargetCameraUpdate()
}
