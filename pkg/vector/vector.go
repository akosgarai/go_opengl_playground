package vector

import (
	"math"
)

type Vector struct {
	X, Y, Z float64
}

// Add returns the sum of the 2 vector.
func (v Vector) Add(ov Vector) Vector {
	return Vector{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z}
}

// MultiplyScalar returns the vector of the multiplication of the given vector and scalar.
func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

// MultiplyVector return the vector of the scalar multiplication of the given vectors.
// aka scalar multiplication.
func (v Vector) MultiplyVector(ov Vector) Vector {
	return Vector{v.X * ov.X, v.Y * ov.Y, v.Z * ov.Z}
}

// Cross return the vector of the vectorial multiplication of the given vectors.
// aka vectorial multiplication
func (v Vector) Cross(ov Vector) Vector {
	return Vector{
		v.Y*ov.Z - v.Z*ov.Y,
		v.Z*ov.X - v.X*ov.Z,
		v.X*ov.Y - v.Y*ov.X,
	}
}

// Length returns the length of the given vector.
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns the normalized vector
func (v Vector) Normalize() Vector {
	if v.X == 0 && v.Y == 0 && v.Z == 0 {
		return v
	}
	return v.MultiplyScalar(1 / v.Length())
}
