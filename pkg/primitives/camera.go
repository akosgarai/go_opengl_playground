package primitives

import (
	"math"
	"strconv"
)

type Camera struct {
	// Camera options

	// Eular Angles
	pitch float64
	yaw   float64

	// Camera attributes
	pos     Vector
	front   Vector
	up      Vector
	right   Vector
	worldUp Vector
	// Projection options.
	projectionOptions struct {
		fov         float64
		aspectRatio float64

		far  float64
		near float64
	}
}

// Log returns the string representation of this object.
func (c *Camera) Log() string {
	logString := "Pos: Vector{" + c.pos.ToString() + "}\n"
	logString += "worldUp: Vector{" + c.worldUp.ToString() + "}\n"
	logString += "front: Vector{" + c.front.ToString() + "}\n"
	logString += "up: Vector{" + c.up.ToString() + "}\n"
	logString += "right: Vector{" + c.right.ToString() + "}\n"
	logString += "yaw : " + strconv.FormatFloat(c.yaw, 'f', 6, 64) + "\n"
	logString += "pitch : " + strconv.FormatFloat(c.pitch, 'f', 6, 64) + "\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + strconv.FormatFloat(c.projectionOptions.fov, 'f', 6, 64) + "\n"
	logString += " - aspectRatio : " + strconv.FormatFloat(c.projectionOptions.aspectRatio, 'f', 6, 64) + "\n"
	logString += " - far : " + strconv.FormatFloat(c.projectionOptions.far, 'f', 6, 64) + "\n"
	logString += " - near : " + strconv.FormatFloat(c.projectionOptions.near, 'f', 6, 64) + "\n"
	return logString
}

// Returns a new camera with the given setup
// position - the camera or eye position
// worldUp - the up direction in the world coordinate system
// yaw - the rotation in z
// pitch - the rotatio in y
func NewCamera(position, worldUp Vector, yaw, pitch float64) *Camera {
	cam := Camera{
		pitch:   pitch,
		yaw:     yaw,
		pos:     position,
		up:      Vector{0, 1, 0},
		worldUp: worldUp,
	}

	cam.updateVectors()
	return &cam
}

// Walk updates the position (forward, back directions)
func (c *Camera) Walk(amount float64) {
	c.pos = c.pos.Add(c.front.MultiplyScalar(amount))
	c.updateVectors()
}

// Strafe updates the position (left, right directions)
func (c *Camera) Strafe(amount float64) {
	c.pos = c.pos.Add(c.front.Cross(c.up).Normalize().MultiplyScalar(amount))
	c.updateVectors()
}

// Lift updates the position (up, down directions)
func (c *Camera) Lift(amount float64) {
	c.pos = c.pos.Add(c.right.Cross(c.front).Normalize().MultiplyScalar((amount)))
	c.updateVectors()
}

// SetupProjection sets the projection related variables
// fov - field of view
// aspectRation - windowWidth/windowHeight
// near - near clip plane
// far - far clip plane
func (c *Camera) SetupProjection(fov, aspRatio, near, far float64) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionOptions.near = near
	c.projectionOptions.far = far
}

// GetProjectionMatrix returns the projectionMatrix of the camera
// based on the following solution: https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/project.go
func (c *Camera) GetProjectionMatrix() *Matrix4x4 {
	return Perspective(float32(c.projectionOptions.fov), float32(c.projectionOptions.aspectRatio), float32(c.projectionOptions.near), float32(c.projectionOptions.far))
}

// GetViewMatrix gets the matrix to transform from world coordinates to
// this camera's coordinates.
// GetViewMatrix returns the viewMatrix of the camera
func (c *Camera) GetViewMatrix() *Matrix4x4 {
	return LookAt(c.pos, c.pos.Add(c.front), c.up)
}
func (c *Camera) updateVectors() {
	c.front = Vector{
		math.Cos(DegToRad(c.pitch)) * math.Cos(DegToRad(c.yaw)),
		math.Sin(DegToRad(c.pitch)),
		math.Cos(DegToRad(c.pitch)) * math.Sin(DegToRad(c.yaw)),
	}
	c.front = c.front.Normalize()
	// Gram-Schmidt process to figure out right and up vectors
	c.right = c.worldUp.Cross(c.front).Normalize()
	c.up = c.right.Cross(c.front).Normalize()
}
