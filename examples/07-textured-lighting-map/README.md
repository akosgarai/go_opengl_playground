# Textured lighting map

The purpose of this application is to demonstrate the lighting maps and it's usage. Based on the [learnopengl](https://learnopengl.com/Lighting/Lighting-maps) tutorial. It draws a textured cube and a sphere as lighting source to the screen. The light source is rotating around the cube. The application contains a camera, so that we can move around with the `W`, `Q`, `A`, `S`, `D`, `E` keys and with the mouse.

The application could be started with a settings screen, where the position, size of the items, the background color, lightsource, and camera parameters could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/07-textured-lighting-map/app.go
```

![Sample gif app without settings](./sample/sample.gif)

- with settings:

```
SETTINGS=on go run examples/07-textured-lighting-map/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
