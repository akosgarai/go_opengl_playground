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
func LookAt(eye, lookAt, worldUp Vector) *Matrix4x4 {
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
	return result.Dot(translationMatrix)
}
