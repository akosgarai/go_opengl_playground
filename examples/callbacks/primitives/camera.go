package primitives

import (
	"math"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	// Camera options

	// Eular Angles
	pitch float32
	yaw   float32

	// Camera attributes
	cameraPosition       mgl32.Vec3
	cameraFrontDirection mgl32.Vec3
	cameraUpDirection    mgl32.Vec3
	cameraRightDirection mgl32.Vec3
	worldUp              mgl32.Vec3
	// Projection options.
	projectionOptions struct {
		fov         float32
		aspectRatio float32

		far  float32
		near float32
	}
}

// Vec3ToString helper function for the string representation of a vector. It is for the log.
func Vec3ToString(v mgl32.Vec3) string {
	x := strconv.FormatFloat(v.X, 'f', 6, 32)
	y := strconv.FormatFloat(v.Y, 'f', 6, 32)
	z := strconv.FormatFloat(v.Z, 'f', 6, 32)
	return "X : " + x + ", Y : " + y + ", Z : " + z
}

// Log returns the string representation of this object.
func (c *Camera) Log() string {
	logString := "cameraPosition: Vector{" + Vec3ToString(c.cameraPosition) + "}\n"
	logString += "worldUp: Vector{" + Vec3ToString(c.worldUp.ToString() + "}\n"
	logString += "cameraFrontDirection: Vector{" + Vec3ToString(c.cameraFrontDirection) + "}\n"
	logString += "cameraUpDirection: Vector{" + Vec3ToString(c.cameraUpDirection) + "}\n"
	logString += "cameraRightDirection: Vector{" + Vec3ToString(c.cameraRightDirection) + "}\n"
	logString += "yaw : " + strconv.FormatFloat(c.yaw, 'f', 6, 32) + "\n"
	logString += "pitch : " + strconv.FormatFloat(c.pitch, 'f', 6, 32) + "\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + strconv.FormatFloat(c.projectionOptions.fov, 'f', 6, 32) + "\n"
	logString += " - aspectRatio : " + strconv.FormatFloat(c.projectionOptions.aspectRatio, 'f', 6, 32) + "\n"
	logString += " - far : " + strconv.FormatFloat(c.projectionOptions.far, 'f', 6, 32) + "\n"
	logString += " - near : " + strconv.FormatFloat(c.projectionOptions.near, 'f', 6, 32) + "\n"
	return logString
}

// Returns a new camera with the given setup
// position - the camera or eye position
// worldUp - the up direction in the world coordinate system
// yaw - the rotation in z
// pitch - the rotation in y
func NewCamera(position, worldUp mgl32.Vec3, yaw, pitch float32) *Camera {
	cam := Camera{
		pitch:             pitch,
		yaw:               yaw,
		cameraPosition:    position,
		cameraUpDirection: mgl32.Vec3{0, 1, 0},
		worldUp:           worldUp,
	}

	cam.updateVectors()
	return &cam
}

// Walk updates the position (forward, back directions)
func (c *Camera) Walk(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraFrontDirection.Mul(amount))
	c.updateVectors()
}

// Strafe updates the position (left, right directions)
func (c *Camera) Strafe(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraFrontDirection.Cross(c.cameraUpDirection).Normalize().Mul(amount))
	c.updateVectors()
}

// Lift updates the position (up, down directions)
func (c *Camera) Lift(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraRightDirection.Cross(c.cameraFrontDirection).Normalize().Mul((amount)))
	c.updateVectors()
}

// SetupProjection sets the projection related variables
// fov - field of view
// aspectRation - windowWidth/windowHeight
// near - near clip plane
// far - far clip plane
func (c *Camera) SetupProjection(fov, aspRatio, near, far float32) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionOptions.near = near
	c.projectionOptions.far = far
}

// GetProjectionMatrix returns the projectionMatrix of the camera
func (c *Camera) GetProjectionMatrix() *mat.Matrix {
	return mgl32.Perspective(c.projectionOptions.fov, c.projectionOptions.aspectRatio, c.projectionOptions.near, c.projectionOptions.far)
}

// GetViewMatrix gets the matrix to transform from world coordinates to
// this camera's coordinates.
// GetViewMatrix returns the viewMatrix of the camera
func (c *Camera) GetViewMatrix() *mat.Matrix {
	return mgl32.LookAtV(c.cameraPosition, c.cameraPosition.Add(c.cameraFrontDirection), c.cameraUpDirection)
}

func (c *Camera) updateVectors() {
	c.cameraFrontDirection = vec.Vector{
		math.Cos(mgl32.DegToRad(c.pitch)) * math.Cos(mgl32.DegToRad(c.yaw)),
		math.Sin(mgl32.DegToRad(c.pitch)),
		math.Cos(mgl32.DegToRad(c.pitch)) * math.Sin(mgl32.DegToRad(c.yaw)),
	}
	c.cameraFrontDirection = c.cameraFrontDirection.Normalize()
	// Gram-Schmidt process to figure out right and up vectors
	c.cameraRightDirection = c.worldUp.Cross(c.cameraFrontDirection).Normalize()
	c.cameraUpDirection = c.cameraRightDirection.Cross(c.cameraFrontDirection).Normalize()
}

// UpdateDirection updates the pitch and yaw values.
func (c *Camera) UpdateDirection(amountX, amountY float32) {
	c.pitch = float32(math.Mod(float64(c.pitch+amountY), 360))
	c.yaw = float32(math.Mod(float64(c.yaw+amountX), 360))
	c.updateVectors()
}
