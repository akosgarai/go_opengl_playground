# Multiple light source

This app aims to implement an application based on the [ligh casters](https://learnopengl.com/Lighting/Light-casters) and the [multiple light](https://learnopengl.com/Lighting/Multiple-lights) tutorials. I'm planning to draw a textured rectangle as a floor plane, a textured cube, a directional light, 2 spotlight and 2 point light.
There is a textured square (1000 * 1000) on the `x-z` plane. The center point is the origo The texture is grass from [here](https://pixabay.com/hu/photos/r%C3%A9t-f%C5%B1-strukt%C3%BAra-anyagminta-halme-253616/).
The box was copy-pasted from a previous application. Only the position & size were updated. The lamps has been made with the StreetLampBuilder tool, so that they are used as spot light sources.
The texture for the `Bug2` sphere was downloaded from [here](https://www.solarsystemscope.com/textures/).

The application could be started with a settings screen, where the position of the items, the background color, lightsource, and camera parameters could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/08-multiple-light/app.go
```

![Sample gif app without settings](./sample/sample.gif)

- with settings:

```
SETTINGS=on go run examples/08-multiple-light/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
