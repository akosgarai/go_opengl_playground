# Cube

It represents a cube. It's described with it's 8 (A-H) `Points`.

## NewCubeByPoints

It returns a new cube defined by it's points. The order of the points matter.

## NewCubeByVectorAndLength

It returns a new cube defined by it's bottom down left point and length of the edge.

### Draw

It draws the cube. The `setupVao` function creates the VAO. The `setupVaoWithColor` updates the colors (different color for each side) setups the coordinate and color pointers.
