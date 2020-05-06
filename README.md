# opengl playground

Just for fun. How to implement 3d applications in golang.

## Useful links:

- [Godoc glfw](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw)
- [Godoc mgl32](https://godoc.org/github.com/go-gl/mathgl/mgl32)
- [Godoc gl](https://godoc.org/github.com/go-gl/gl/v4.1-core/gl)
- [Learnopengl](https://learnopengl.com/) - good explanations and cpp examples.
- [About glsl](https://www.khronos.org/opengl/wiki/OpenGL_Shading_Language)
- A tutorial [first part](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl) and [second part](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-2-drawing-the-game-board)
- [Other tutorial](https://medium.com/@drgomesp/opengl-and-golang-getting-started-abcd3d96f3db)
- [About transformations](http://www.codinglabs.net/article_world_view_projection_matrix.aspx)

## About the applications

- Dependencies are handled with gomod.
- The example readmes might be outdated.
- Light application is very wip.
- How to run the example apps?

In the main directory run the following command, after you replaced the directory name to a valid one.

```
go run examples/directory-name/app.go
```

## Possible issues ubuntu.

- Opengl version mismatch.

The applications are using the opengl 4.1 package. If your version is same or higher, the appliactions should run without issues.
To check your opengl version just run the following command in terminal (based on [this](https://askubuntu.com/questions/47062/what-is-terminal-command-that-can-show-opengl-version)):

```bash
glxinfo | grep "OpenGL version"
```

The output is something like: `OpenGL version string: 4.6.0 NVIDIA 440.82`.


## v3.3

This is a test branch. Check whether the current apps are working with the opengl 3.3 version or not.

The following steps needs to be done, to replace the versions:

```bash
sed -i 's%github.com/go-gl/gl/v4.1-core/gl%github.com/go-gl/gl/v3.3-core/gl%g' examples/*/*.go
sed -i 's%github.com/go-gl/gl/v4.1-core/gl%github.com/go-gl/gl/v3.3-core/gl%g' pkg/*/*.go
sed -i 's%github.com/go-gl/gl/v4.1-core/gl%github.com/go-gl/gl/v3.3-core/gl%g' pkg/primitives/*/*.go
sed -i 's%#version 410%#version 330%g' examples/*/*
```

Manually update the window pkg.

```
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
```
