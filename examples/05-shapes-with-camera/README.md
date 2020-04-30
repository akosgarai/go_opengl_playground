# Shapes with camera

The purpose of this application is the demonstration of the currently existing shapes. The camera is controlable. We can change the position with the `q`, `w`, `e`, `a`, `s`, `d` keys. With the mouse, the camera direction can be changed.

## Application

This stucture represents our world. It contains the shapes and the options for the camera.

### NewApplication

It returns an application with the necessary setup.

### GenerateCube

It creates the cube.

### GenerateSphere

It creates the sphere.

### KeyHandler

This function is responsible for the movement handling (`q`, `w`, `e`, `a`, `s`, `d` keys) and for the debug log (`h`).

### MouseHandler

This function is responsible for the mouse movement handling. If the mouse is in the target zone, it updates the camera direction.
