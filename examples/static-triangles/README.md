# Static triangles app

This application drows triangles to the screen. It hasn't any mouse or keyboard control, just the image on the screen.

## Application

This stucture represents our world. It contains the triangles and the necessary options for generating them.

### NewApplication

This function returns a configured application. It's input the number of rows. The columns and the triangle lengths are calculated based on that value.

### GenerateTriangles

This function is responsible for the triangle generation. It calculates the coordinates based on the number of rows

### Draw

Draws the triangles.

## Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderBasicSource` vertex shader.
