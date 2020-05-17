# Cuboid

This package represents a cuboid, that is described with the points of the rectangles of the sides. The cuboid has 6 sides, each side is describes with 4 points, so that we have 24 points. Each side point so different direction, so that the length of Normals is 6.

## New

This function returns a cuboid. The inputs are the width (length in the `X` axis) the length (length in the `Z` axis) and the height (length in the `Y` axis) in this order. The sides are scaled, the longest side becomes 1. The center point is the origo, where center point is the intersection point of the diagonals.

## NewCube

This function returns a unit cube with origo as center point, where center point is the intersection point of the diagonals.

## MeshInput

MeshInput method returns the verticies, indicies - inputs for the New Mesh function.

## ColoredMeshInput

ColoredMeshInput method returns the verticies, indicies - inputs for the New Mesh function.

## TexturedColoredMeshInput

This method returns the verticies, indicies - inputs for the NewTexturedColoredMesh function.
