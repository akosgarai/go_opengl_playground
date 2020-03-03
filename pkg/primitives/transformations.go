package primitives

import (
	"math"
)

var (
	EPSILON = 0.0001
)

// returns the given degree in radian
func DegToRad(deg float64) float64 {
	return (deg * math.Pi / 180)
}

// returns the given radian in degree
func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// https://stackoverflow.com/questions/8115352/glmperspective-explanation
// ProjectionNewSolution.
// https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/project.go - based on this, it's fine.
func Perspective(angle, ratio, near, far float32) *Matrix4x4 {
	// degree to radian formula: n deg = n * PI / 180 rad
	slopey := float32(math.Tan(float64(angle * math.Pi / 180)))
	result := NullMatrix4x4()
	result.Points[0] = 1 / slopey / ratio
	result.Points[5] = 1 / slopey
	result.Points[10] = -((far + near) / (far - near))
	result.Points[11] = -1
	result.Points[14] = -(2 * far * near / (far - near))
	return result
}

// https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/project.go#L48-L61
// updates based on the link above.
func LookAt_v4(eye, center, up Vector) *Matrix4x4 {
	f := center.Subtract(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)
	M := Matrix4x4{[16]float32{
		float32(s.X), float32(u.X), -float32(f.X), 0,
		float32(s.Y), float32(u.Y), -float32(f.Y), 0,
		float32(s.Z), float32(u.Z), -float32(f.Z), 0,
		0, 0, 0, 1,
	}}
	return M.Mul4(TranslationMatrix4x4(float32(-eye.X), float32(-eye.Y), float32(-eye.Z)).TransposeMatrix())
}
