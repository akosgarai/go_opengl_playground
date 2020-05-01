package transformations

import (
	"math"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
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

func MouseCoordinates(currentX, currentY, windowWidth, windowHeight float64) (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (currentX - halfWidth) / (halfWidth)
	y := (halfHeight - currentY) / (halfHeight)
	return x, y
}

// Vec3ToString helper function for the string representation of a vector. It is for the log.
func Vec3ToString(v mgl32.Vec3) string {
	x := strconv.FormatFloat(float64(v.X()), 'f', 10, 32)
	y := strconv.FormatFloat(float64(v.Y()), 'f', 10, 32)
	z := strconv.FormatFloat(float64(v.Z()), 'f', 10, 32)
	return "X : " + x + ", Y : " + y + ", Z : " + z
}

// Float64ToString returns the given float number in string format.
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 10, 32)
}

// Float32ToString returns the given float number in string format.
func Float32ToString(num float32) string {
	return Float64ToString(float64(num))
}

// IntegerToString returns the string representation of the given integer
func IntegerToString(num int) string {
	return strconv.Itoa(num)
}
