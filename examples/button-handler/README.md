# Static button handler

This application drows a square and a triangle to the screen. The color of the shapes can be swapped with the `T` button. Until the button is pressed, the colors are swapped. The cube's position is changable. It goes forward with the `W`, backward with the `S`, left with the `A` and right with the `D` button. The triangle also changes it's position, it moves to the opposite directions.

## Functions

- `initGlfw`

Basic function for glfw initialization.

- `initOpenGL`

It is responsible for openGL initialization. It uses the `shader.FragmentShaderBasicSource` fragment shader and the `shader.VertexShaderBasicSource` vertex shader.

- `keyHandler`

This functions is responsible for the keyboard actions (`W`, `A`, `S`, `D`, `T`).
