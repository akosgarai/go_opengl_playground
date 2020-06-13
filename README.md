# Opengl Playground

Just for fun. How to implement 3d applications in golang. The 3D engine used to be in this repo, but it was difficult to manage everything inside one repository, so i decided to move the engine to a [separate repo](https://github.com/akosgarai/playground_engine).

Now this repo contains only the example application that i have written with the engine.
The gifs under the examples directory were made with [peek](https://github.com/phw/peek) application.

## About the applications

- Dependencies are handled with gomod.
- How to run the example apps?

In the main directory run the following command, after you replaced the directory name with a valid one.

```
go run examples/directory-name/app.go
```

![Sample gif from outer space](./examples/07-textured-spheres/sample/sample.gif)

## Possible issues ubuntu.

- Opengl version mismatch.

The applications are using the opengl 4.1 package. If your version is same or higher, the appliactions should run without issues.
To check your opengl version just run the following command in terminal (based on [this](https://askubuntu.com/questions/47062/what-is-terminal-command-that-can-show-opengl-version)):

```bash
glxinfo | grep "OpenGL version"
```

The output is something like: `OpenGL version string: 4.6.0 NVIDIA 440.82`.
