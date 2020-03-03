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

// YPR aka YawPitchRoll aka GetMatrixUsingYawPitchRoll creates a matrix that represents the transformations
// the input parameters are in radians.
func YPR(yaw, pitch, roll float64) *Matrix4x4 {
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	cosPitch := math.Cos(pitch)
	sinPitch := math.Sin(pitch)

	cosRoll := math.Cos(roll)
	sinRoll := math.Sin(roll)

	result := NullMatrix4x4()
	result.Points[0] = float32(cosYaw*cosRoll + sinYaw*sinPitch*sinRoll)
	result.Points[1] = float32(sinRoll * cosPitch)
	result.Points[2] = float32(-sinYaw*cosRoll + cosYaw*sinPitch*sinRoll)
	result.Points[4] = float32(-cosYaw*sinRoll + sinYaw*sinPitch*cosRoll)
	result.Points[5] = float32(cosRoll * cosPitch)
	result.Points[6] = float32(sinRoll*sinYaw + cosYaw*sinPitch*cosRoll)
	result.Points[8] = float32(sinYaw * cosPitch)
	result.Points[9] = float32(-sinPitch)
	result.Points[10] = float32(cosYaw * cosPitch)
	result.Points[15] = 1

	return result
}

// implementation of the gluLookAt function (based on the Szirmay graph book. page: 194)
func LookAtV(right, up, forward, position Vector) *Matrix4x4 {
	result := NullMatrix4x4()
	result.Points[0] = float32(right.X)
	result.Points[1] = float32(up.X)
	result.Points[2] = float32(forward.X)
	result.Points[4] = float32(right.Y)
	result.Points[5] = float32(up.Y)
	result.Points[6] = float32(forward.Y)
	result.Points[8] = float32(right.Z)
	result.Points[9] = float32(up.Z)
	result.Points[10] = float32(forward.Z)

	result.Points[12] = float32(-right.X*position.X - right.Y*position.Y - right.Z*position.Z)
	result.Points[13] = float32(-up.X*position.X - up.Y*position.Y - up.Z*position.Z)
	result.Points[14] = float32(-forward.X*position.X - forward.Y*position.Y - forward.Z*position.Z)

	result.Points[15] = 1
	return result
	//translationMatrix := TranslationMatrix4x4(-float32(position.X), -float32(position.Y), -float32(position.Z))
	//return translationMatrix.Dot(result)
}
func LookAt(eye, lookAt, worldUp Vector) *Matrix4x4 {
	// u - vizszintes,  v - fuggoleges, w - elore
	w := eye.Subtract(lookAt)
	wNorm := w.Length()
	if wNorm > EPSILON {
		w = w.DivideScalar(wNorm)
	} else {
		w.Z = 1
		w.X = 0
		w.Y = 0
	}
	u := worldUp.Cross(w)
	uNorm := u.Length()
	if uNorm > EPSILON {
		u = u.DivideScalar(uNorm)
	} else {
		u.X = 1
		u.Y = 0
		u.Z = 0
	}
	v := w.Cross(u)
	result := NullMatrix4x4()
	// u, v, w, 0
	result.Points[0] = float32(u.X)
	result.Points[1] = float32(v.X)
	result.Points[2] = float32(w.X)
	result.Points[4] = float32(u.Y)
	result.Points[5] = float32(v.Y)
	result.Points[6] = float32(w.Y)
	result.Points[8] = float32(u.Z)
	result.Points[9] = float32(v.Z)
	result.Points[10] = float32(w.Z)
	result.Points[15] = 1
	translationMatrix := TranslationMatrix4x4(-float32(eye.X), -float32(eye.Y), -float32(eye.Z))
	return translationMatrix.Dot(result)
}

// Helper function to transform the mouse coordinates.
func ConvertMouseCoordinates(mousePositionX, mousePositionY, windowWidth, windowHeight float64) (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (mousePositionX - halfWidth) / (halfWidth)
	y := (halfHeight - mousePositionY) / (halfHeight)
	return x, y
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
