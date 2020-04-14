package primitives

func MouseCoordinates(currentX, currentY, windowWidth, windowHeight float64) (float64, float64) {
	halfWidth := windowWidth / 2.0
	halfHeight := windowHeight / 2.0
	x := (currentX - halfWidth) / (halfWidth)
	y := (halfHeight - currentY) / (halfHeight)
	return x, y
}
