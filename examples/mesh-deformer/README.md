# Mesh deformer

This application demonstrates the mesh deformer vertex shader. The movement is based on a periodic (sin) function. 2 applications has been written one with camera (app-with-camera.go) and one without camera (app.go)

## App

This is the "static" version of the mesh deformer. Static means, that the camera is static, so we can't change it's position or direction from the application.

### Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderDeformVertexPositionSource` vertex shader.

## App with camera

This one is the same app, but with moving camera. We can change the position with the `q`, `w`, `e`, `a`, `s`, `d` keys. With the mouse, the camera direction can be changed.

### Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderDeformVertexPositionModelViewProjectionSource` vertex shader.
