package primitives

import (
	"math"
)

type Camera struct {
	lowerLeft, horizontal, vertical, origin, u, v, w Vector
	lensRadius                                       float64
}

func NewCamera(lookFrom, lookAt, vUp Vector, vFov, aspect, aperture float64) *Camera {
	c := Camera{}

	c.origin = lookFrom
	c.lensRadius = aperture / 2

	theta := vFov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	w := lookFrom.Subtract(lookAt).Normalize()
	u := vUp.Cross(w).Normalize()
	v := w.Cross(u)

	focusDist := lookFrom.Subtract(lookAt).Length()

	x := u.MultiplyScalar(halfWidth * focusDist)
	y := v.MultiplyScalar(halfHeight * focusDist)

	c.lowerLeft = c.origin.Subtract(x).Subtract(y).Subtract(w.MultiplyScalar(focusDist))
	c.horizontal = x.MultiplyScalar(2)
	c.vertical = y.MultiplyScalar(2)

	c.w = w
	c.u = u
	c.v = v

	return &c
}
