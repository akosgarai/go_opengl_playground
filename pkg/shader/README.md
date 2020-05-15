# Shader


This readme is outdated. the package tests are broken. Needs to be fixed and updated.

This package is a kind of wrapper for the gl.\* commands.

## Functions

The stuff that we can do with the shaders.

### LoadShaderFromFile

LoadShaderFromFile takes a filepath string arguments. It loads the file and returns it as a `\x00` terminated string. It returns an error also.

### CompileShader

It compiles the given shader program string as the given type of shader. In case of success the shader program id and nil returns. On case of error, it's returned instead of the nil.

### NewShader

NewShader returns a Shader. It's inputs are the filenames of the shaders. It reads the files and compiles them. The shaders are attached to the shader program.

### Use

Use is a wrapper for gl.UseProgram

### SetUniformMat4

SetUniformMat4 gets an uniform name string and the value matrix as input and calls the gl.UniformMatrix4fv function

### SetUniformMat3

SetUniformMat3 gets an uniform name string and the value matrix as input and calls the gl.UniformMatrix3fv function

### SetUniform3f

SetUniform3f gets an uniform name string and 3 float values as input and calls the gl.Uniform3f function

### SetUniform1f

SetUniform1f gets an uniform name string and a float value as input and calls the gl.Uniform1f function

### BindBufferData

BindBufferData gets a float array as an input, generates a buffer binds it as array buffer, and sets the input as buffer data.

### BindVertexArray

BindVertexArray generates a vertex array and binds it.

### VertexAttribPointer

VertexAttribPointer sets the pointer.

### Close

Close disables the vertexarraypointers and the vertex array.

### DrawPoints

DrawPoints is the draw functions for points

### DrawTriangles

DrawTriangles is the draw function for triangles
