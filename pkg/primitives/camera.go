package primitives

// based on the opengl dev. cookbook solution

type Camera struct {
	// camera position
	position Vector
	// look at position (the camera is looking at lookat from it's position.)
	LookAt Vector
	// world up direction.
	UpDirection Vector

	fov         float64
	aspectRatio float64

	cameraUpDirection    Vector
	cameraRightDirection Vector
	// view matrix
	viewMatrix Matrix4x4
	// projection matrix
	projectionMatrix Matrix4x4

	freeCameraOptions struct {
	    translation Vector
	    yaw float64
	    pitch float64
	    roll float64
	}
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
func (c *Camera) GetViewMatrix() Matrix4x4 {
	return c.viewMatrix
}

// GetProjectionMatrix returns the viewMatrix of the camera
func (c *Camera) GetProjectionMatrix() Matrix4x4 {
	return c.projectionMatrix
}

// GetPosition returns the position of the camera
func (c *Camera) GetPosition() Vector {
	return c.position
}
// SetPosition updates the position of the camera
func (c *Camera) SetPosition(newPos Vector) {
	return c.position = newPos
}

// SetupProjection creates the projection matrix and setups the fow and the aspect_ration
func (c *Camera) SetupProjection(fov, aspRatio float64) {
	c.fov = fov
	c.aspectRatio = aspRatio
	c.projectionMatrix = Perspective(float32(fov), float32(aspRatio), 0.1, 1000)
}

// YPR aka YawPitchRoll creates a matrix that represents the transformations
func (c *Camera) YPR(yaw, pitch, roll float64) {
    cosYaw = math.Cos(yaw)
    sinYaw = math.Sin(yaw)

    cosPitch = math.Cos(pitch)
    sinPitch = math.Sin(pitch)

    cosRoll = math.Cos(roll)
    sinRoll = math.Sin(roll)

    result := NullMatrix4x4()
    result.Points[0] = float32(cosYaw * cosRoll + sinYaw * sinPitch * sinRoll)
    result.Points[1] = float32(sinRoll * cosPitch)
    result.Points[2] = float32(- sinYaw * cosRoll + cosYaw * sinPitch * sinRoll)
    result.Points[4] = float32(- cosYaw * sinRoll + sinYaw * sinPitch * cosRoll)
    result.Points[5] = float32(cosRoll * cosPitch)
    result.Points[6] = float32(sinRoll * sinYaw + cosYaw * sinPitch * cosRoll)
    result.Points[8] = float32(sinYaw * cosPitch)
    result.Points[9] = float32(- sinPitch)
    result.Points[10] = float32(cosYaw * cosPitch)
    result.Points[15] = 1

    return result
}
