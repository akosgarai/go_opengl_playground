# Textured spheres

This application aims to show how to draw textured spheres. The textures for the spheres were downloaded from [here](https://www.solarsystemscope.com/textures/).
The skybox textures were generated with [wwwtyro](https://wwwtyro.github.io/space-3d/#animationSpeed=1&fov=80&nebulae=true&pointStars=true&resolution=1024&seed=2hnqv2e7hhg0&stars=true&sun=true).

The application could be started with a settings screen, where the position, size of the items, the background color, lightsource, and camera parameters could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/07-textured-spheres/app.go
```

![Sample gif app without settings](./sample/sample.gif)

- with settings:

```
SETTINGS=on go run examples/07-textured-spheres/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
