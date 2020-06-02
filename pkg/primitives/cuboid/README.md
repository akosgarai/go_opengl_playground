# Cuboid

This package represents a cuboid, that is described with the points of the rectangles of the sides. The cuboid has 6 sides, each side is describes with 4 points, so that we have 24 points. Each side point so different direction, so that the length of Normals is 6.

## New

This function returns a cuboid. The inputs are the width (length in the `X` axis) the length (length in the `Z` axis) and the height (length in the `Y` axis) in this order. The center point is the origo, where center point is the intersection point of the diagonals.

## NewCube

This function returns a unit cube with origo as center point, where center point is the intersection point of the diagonals.

## TexturedMeshInput

TexturedMeshInput method returns the vertices, indices - inputs for the NewTexturedMesh function.

## MaterialMeshInput

MaterialMeshInput method returns the vertices, indices - inputs for the NewMaterialMesh function.

## ColoredMeshInput

ColoredMeshInput method returns the vertices, indices - inputs for the NewColorMesh function.

## TexturedColoredMeshInput

This method returns the vertices, indices - inputs for the NewTexturedColoredMesh function.

---

## Texture orientations

Currently it supports the `default` orientation, that means every texture is positioned in the same order, and the `same` ordet that means that the textures on the opposite sides are oriented opposite.
