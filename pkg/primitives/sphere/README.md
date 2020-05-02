# Sphere

It represents a sphere, that is described with it's center, radius and color, a direction vector + speed (for moving objects).
It has a VAO and a Shader also. It implements the Drawable interface.

## Functions

The stuff that we can do with a sphere.

## New

It creates a new Sphere. The inputs of this functions are the center, the color the radius and the shader. The direction & speed is initialized as null vector & 0 speed.

### Log

The string representation of the current state of the object.

## SetCenter

It updates the center of the sphere.

## GetCenter

it returns the center of the sphere.

## SetColor

It updates the color of the sphere.

## GetColor

It returns the color of the sphere.

## SetRadius

It updates the radius of the sphere.

## GetRadius

It returns the radius of the sphere.

### SetDirection

It updates the direction of the sphere with the given new one.

### SetIndexDirection

It updates the indexed direction of the sphere to the given value.

### SetSpeed

It updates the speed of the sphere.

### Draw

It draws the sphere. Transformations are not applied in this case.

### DrawWithUniforms

It draws the sphere. It gets the V & P matrices as inputs. It sets the model, view, projection uniforms for the shader program.

### Update

It updates the state of the sphere. It gets the delta time as input and it calculates the movement of the sphere.
