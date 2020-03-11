package vector

import (
	"math"
	"strconv"
)

type Vector struct {
	X, Y, Z float64
}

var UnitVector = Vector{1, 1, 1}

// ToString returns the string representation of the vector
func (v Vector) ToString() string {
	x := strconv.FormatFloat(v.X, 'f', 6, 64)
	y := strconv.FormatFloat(v.Y, 'f', 6, 64)
	z := strconv.FormatFloat(v.Z, 'f', 6, 64)
	return "X : " + x + ", Y : " + y + ", Z : " + z
}

// Length returns the length of the vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns the normalized vector. This operand makes the length of the vector 1.
func (v Vector) Normalize() Vector {
	return v.DivideScalar(v.Length())
}

// Dot aka SquaredLength returns the squared length of the vector. (length^2)
func (v Vector) Dot(ov Vector) float64 {
	return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z
}

// Cross returns the cross product vector of the given vectors.
func (v Vector) Cross(ov Vector) Vector {
	return Vector{
		v.Y*ov.Z - v.Z*ov.Y,
		v.Z*ov.X - v.X*ov.Z,
		v.X*ov.Y - v.Y*ov.X,
	}
}

// Add returns the sum of the 2 vectors.
func (v Vector) Add(ov Vector) Vector {
	return Vector{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z}
}

// Subtract returns the substracted vector.
func (v Vector) Subtract(ov Vector) Vector {
	return Vector{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z}
}

// Multiply returns the multiplied vector.
func (v Vector) Multiply(ov Vector) Vector {
	return Vector{v.X * ov.X, v.Y * ov.Y, v.Z * ov.Z}
}

// Devide returns the divided vector.
func (v Vector) Divide(ov Vector) Vector {
	return Vector{v.X / ov.X, v.Y / ov.Y, v.Z / ov.Z}
}

// AddScalar returns the incremented vector.
// v.AddScalar(t) = v.Add(vector{t,t,t})
func (v Vector) AddScalar(t float64) Vector {
	return Vector{v.X + t, v.Y + t, v.Z + t}
}

// SubtractScalar returns the decremented vector.
// v.SubtractScalar(t) = v.Subtract(vector{t,t,t})
func (v Vector) SubtractScalar(t float64) Vector {
	return Vector{v.X - t, v.Y - t, v.Z - t}
}

// MultiplyScalar returns the multiplied vector.
// v.MultiplyScalar(t) = v.Multiply(vector{t,t,t})
func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

// DivideScalar returns the divided vector.
// v.DivideScalar(t) = v.Divide(vector{t,t,t})
func (v Vector) DivideScalar(t float64) Vector {
	return Vector{v.X / t, v.Y / t, v.Z / t}
}

/*
func (v Vector) Reflect(ov Vector) Vector {
	b := 2 * v.Dot(ov)
	return v.Subtract(ov.MultiplyScalar(b))
}

func (v Vector) Refract(ov Vector, n float64) (bool, Vector) {
	uv := v.Normalize()
	uo := ov.Normalize()
	dt := uv.Dot(uo)
	discriminant := 1.0 - (n * n * (1 - dt*dt))
	if discriminant > 0 {
		a := uv.Subtract(ov.MultiplyScalar(dt)).MultiplyScalar(n)
		b := ov.MultiplyScalar(math.Sqrt(discriminant))
		return true, a.Subtract(b)
	}
	return false, Vector{}
}

func VectorInUnitSphere(rnd *rand.Rand) Vector {
	for {
		r := Vector{rnd.Float64(), rnd.Float64(), rnd.Float64()}
		p := r.MultiplyScalar(2.0).Subtract(UnitVector)
		if p.Dot() >= 1.0 {
			return p
		}
	}
}
*/
