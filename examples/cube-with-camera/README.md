# Cube with camera

This application demonstrates the camera functionalities. The camera position is controlable with the `W` (forward), `A` (left), `S` (backward), `D` (right), `Q` (down), `E` (up) keys. The look direction is updatable with the mouse (boarder of the window, 8 direction). On the screen, there is a cube with different side colors.

## Application

This stucture represents our world. It contains the cube and the necessary options for drawing it. It also contains the camera and it's options.

### NewApplication

It creates an application with the minimal setup.

### GenerateCube

It sets the cube and it's position.

### KeyHandler

This functions is responsible for the keyboard actions (`W`, `A`, `S`, `D`, `T`, `Q`, `E`, `H`). Button `H` prints debug info to the console.

### MouseHandler

This function is responsible for the mouse movement handling. If the mouse is in the target zone, it updates the camera direction.

## Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderModelViewProjectionSource` vertex shader.
