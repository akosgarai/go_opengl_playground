# Transformations

It contains functions that represents transformations. 

## DegToRad

This function converts a `float64` degree value to `float64` radian.

```
    rad = deg * PI / 180
```

## RadToDeg

This function converts a `float64` radian value to `float64` degree.

```
    deg = rad * 180 / PI
```

## Perspective

This function returns the perspective matrix.

## LookAt

This function returns the transformation matrix from the world space to the camera space.

## MouseCoordinates

This function transforms the mouse coordinates to window coordinate. (windowwidth\*windowheight -> [-1, 1])

## ScaleMatrix

This function returns a scale matrix. The inputs are the scale ratios.

## TranslationMatrix

This function returns a translation matrix. The inputs are the translation values.

## RotationXMatrix

This function returns a rotation matrix, where the rotation is based on the `X` axis.

## RotationYMatrix

This function returns a rotation matrix, where the rotation is based on the `Y` axis.

## RotationZMatrix

This function returns a rotation matrix, where the rotation is based on the `Z` axis.

## ProjectionMatrix

This functions returns a perspective matrix, but with 1 as ratio.
