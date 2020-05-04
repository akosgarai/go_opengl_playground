# Cuboid

It represents a cuboid, so it contains 6 `rectangles` as it's sides. It has a VAO and a Shader also. It implements the Drawable interface.

## Functions

The stuff that we can do with a cuboid.

### New

It creates a new cuboid. The inputs of this functions are a side of the cuboid (rectangle) and the lengt of the perpendicular sides and the shader. The direction & speed is initialized as null vector & 0 speed. The rest of the sides are calculated based on the input side and length.

### Log

The string representation of the current state of the object.

### SetColor

It updates the color of the sides to the given one.

### SetIndexColor

It updates the color of a given pont of each sides.

### SetSideColor

It updates every color of a given side.

### SetDirection

It updates the direction vector of the sides.

### SetIndexDirection

It updates the given direction component to the given value of each sides.

### SetSpeed

It updates the speed of the sides to the given one.

### SetPrecision

It updates the precision of the sides to the given one.

### SetAngle

It updates the rotation angle (radian) of the cube.

### SetAxis

It updates the rotation axis of the cube.

### Draw

It draws the cuboid (rectangles). Transformations are not applied.

### DrawWithUniforms

It draws the cuboid (rectangles). It gets the V & P matrices as inputs. It sets the model, view, projection uniforms for the shader program.

### Update

It updates the state of the cuboid (rectangles). It gets the delta time as input and it calculates the movement of the cuboid (rectangles).
