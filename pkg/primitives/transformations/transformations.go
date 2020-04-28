package transformations

import (
	"math"
	"strconv"

	mat "github.com/akosgarai/opengl_playground/pkg/primitives/matrix"
	vec "github.com/akosgarai/opengl_playground/pkg/primitives/vector"

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

// https://stackoverflow.com/questions/8115352/glmperspective-explanation
// ProjectionNewSolution.
// https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/project.go - based on this, it's fine.
func Perspective(angle, ratio, near, far float32) *mat.Matrix {
	// degree to radian formula: n deg = n * PI / 180 rad
	slopey := float32(math.Tan(float64(angle * math.Pi / 180)))
	result := mat.NullMatrix()
	result.Points[0] = 1 / slopey / ratio
	result.Points[5] = 1 / slopey
	result.Points[10] = -((far + near) / (far - near))
	result.Points[11] = -1
	result.Points[14] = -(2 * far * near / (far - near))
	return result
}

// https://github.com/go-gl/mathgl/blob/95de7b3a016a8324097da95ad4417cc2caccb071/mgl32/project.go#L48-L61
// updates based on the link above.
func LookAt(eye, center, up vec.Vector) *mat.Matrix {
	f := center.Subtract(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)
	M := mat.Matrix{[16]float32{
		float32(s.X), float32(u.X), -float32(f.X), 0,
		float32(s.Y), float32(u.Y), -float32(f.Y), 0,
		float32(s.Z), float32(u.Z), -float32(f.Z), 0,
		0, 0, 0, 1,
	}}
	return M.Dot(TranslationMatrix(float32(-eye.X), float32(-eye.Y), float32(-eye.Z)).TransposeMatrix())
}

func MouseCoordinates(currentX, currentY, windowWidth, windowHeight float64) (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (currentX - halfWidth) / (halfWidth)
	y := (halfHeight - currentY) / (halfHeight)
	return x, y
}
func ScaleMatrix(scaleX, scaleY, scaleZ float32) *mat.Matrix {
	return &mat.Matrix{
		[16]float32{
			scaleX, 0.0, 0.0, 0.0,
			0.0, scaleY, 0.0, 0.0,
			0.0, 0.0, scaleZ, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}
func TranslationMatrix(translationX, translationY, translationZ float32) *mat.Matrix {
	return &mat.Matrix{
		[16]float32{
			1.0, 0.0, 0.0, translationX,
			0.0, 1.0, 0.0, translationY,
			0.0, 0.0, 1.0, translationZ,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

// RotationXMatrix is the transformation matrix for the 'x' axis based rotation.
// The input 'rotationAngle' is in radiant.
func RotationXMatrix(rotationAngle float64) *mat.Matrix {
	return &mat.Matrix{
		[16]float32{
			1.0, 0.0, 0.0, 0.0,
			0.0, float32(math.Cos(rotationAngle)), float32(-math.Sin(rotationAngle)), 0.0,
			0.0, float32(math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

// RotationYMatrix is the transformation matrix for the 'y' axis based rotation.
// The input 'rotationAngle' is in radiant.
func RotationYMatrix(rotationAngle float64) *mat.Matrix {
	return &mat.Matrix{
		[16]float32{
			float32(math.Cos(rotationAngle)), 0.0, float32(math.Sin(rotationAngle)), 0.0,
			0.0, 1.0, 0.0, 0.0,
			float32(-math.Sin(rotationAngle)), 0.0, float32(math.Cos(rotationAngle)), 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

// RotationZMatrix is the transformation matrix for the 'z' axis based rotation.
// The input 'rotationAngle' is in radiant.
func RotationZMatrix(rotationAngle float64) *mat.Matrix {
	return &mat.Matrix{
		[16]float32{
			float32(math.Cos(rotationAngle)), float32(-math.Sin(rotationAngle)), 0.0, 0.0,
			float32(math.Sin(rotationAngle)), float32(math.Cos(rotationAngle)), 0.0, 0.0,
			0.0, 0.0, 1.0, 0.0,
			0.0, 0.0, 0.0, 1.0,
		},
	}
}

func ProjectionMatrix(angleOfView, near, far float64) *mat.Matrix {
	scale := float32(1 / math.Tan(angleOfView*0.5*math.Pi/180))
	projection := mat.UnitMatrix()
	projection.Points[0] = scale
	projection.Points[5] = scale
	projection.Points[10] = float32(-far / (far - near))
	projection.Points[14] = float32(-far * near / (far - near))
	projection.Points[11] = -1
	projection.Points[15] = 0
	return projection
}

// Vec3ToString helper function for the string representation of a vector. It is for the log.
func Vec3ToString(v mgl32.Vec3) string {
	x := strconv.FormatFloat(float64(v.X()), 'f', 6, 32)
	y := strconv.FormatFloat(float64(v.Y()), 'f', 6, 32)
	z := strconv.FormatFloat(float64(v.Z()), 'f', 6, 32)
	return "X : " + x + ", Y : " + y + ", Z : " + z
}

// Float64ToString returns the given float number in string format.
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 32)
}

// Float32ToString returns the given float number in string format.
func Float32ToString(num float32) string {
	return Float64ToString(float64(num))
}

// IntegerToString returns the string representation of the given integer
func IntegerToString(num int) string {
	return strconv.Itoa(num)
}
