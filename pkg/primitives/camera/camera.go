package primitives

import (
	"math"
	"strconv"

	mat "github.com/akosgarai/opengl_playground/pkg/primitives/matrix"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"
)

type Camera struct {
	// Camera options

	// Eular Angles
	pitch float64
	yaw   float64

	// Camera attributes
	cameraPosition       vec.Vector
	cameraFrontDirection vec.Vector
	cameraUpDirection    vec.Vector
	cameraRightDirection vec.Vector
	worldUp              vec.Vector
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
	logString := "cameraPosition: Vector{" + c.cameraPosition.ToString() + "}\n"
	logString += "worldUp: Vector{" + c.worldUp.ToString() + "}\n"
	logString += "cameraFrontDirection: Vector{" + c.cameraFrontDirection.ToString() + "}\n"
	logString += "cameraUpDirection: Vector{" + c.cameraUpDirection.ToString() + "}\n"
	logString += "cameraRightDirection: Vector{" + c.cameraRightDirection.ToString() + "}\n"
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
func NewCamera(position, worldUp vec.Vector, yaw, pitch float64) *Camera {
	cam := Camera{
		pitch:             pitch,
		yaw:               yaw,
		cameraPosition:    position,
		cameraUpDirection: vec.Vector{0, 1, 0},
		worldUp:           worldUp,
	}

	cam.updateVectors()
	return &cam
}

// Walk updates the position (forward, back directions)
func (c *Camera) Walk(amount float64) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraFrontDirection.MultiplyScalar(amount))
	c.updateVectors()
}

// Strafe updates the position (left, right directions)
func (c *Camera) Strafe(amount float64) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraFrontDirection.Cross(c.cameraUpDirection).Normalize().MultiplyScalar(amount))
	c.updateVectors()
}

// Lift updates the position (up, down directions)
func (c *Camera) Lift(amount float64) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraRightDirection.Cross(c.cameraFrontDirection).Normalize().MultiplyScalar((amount)))
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
func (c *Camera) GetProjectionMatrix() *mat.Matrix {
	return Perspective(float32(c.projectionOptions.fov), float32(c.projectionOptions.aspectRatio), float32(c.projectionOptions.near), float32(c.projectionOptions.far))
}

// GetViewMatrix gets the matrix to transform from world coordinates to
// this camera's coordinates.
// GetViewMatrix returns the viewMatrix of the camera
func (c *Camera) GetViewMatrix() *mat.Matrix {
	return LookAt(c.cameraPosition, c.cameraPosition.Add(c.cameraFrontDirection), c.cameraUpDirection)
}
func (c *Camera) updateVectors() {
	c.cameraFrontDirection = vec.Vector{
		math.Cos(DegToRad(c.pitch)) * math.Cos(DegToRad(c.yaw)),
		math.Sin(DegToRad(c.pitch)),
		math.Cos(DegToRad(c.pitch)) * math.Sin(DegToRad(c.yaw)),
	}
	c.cameraFrontDirection = c.cameraFrontDirection.Normalize()
	// Gram-Schmidt process to figure out right and up vectors
	c.cameraRightDirection = c.worldUp.Cross(c.cameraFrontDirection).Normalize()
	c.cameraUpDirection = c.cameraRightDirection.Cross(c.cameraFrontDirection).Normalize()
}

func (c *Camera) UpdateDirection(amountX, amountY float64) {
	c.pitch = math.Mod(c.pitch+amountY, 360)
	c.yaw = math.Mod(c.yaw+amountX, 360)
	c.updateVectors()
}
