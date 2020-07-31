# Static Mesh deformer

This application demonstrates the mesh deformer vertex shader. The movement is based on a periodic (sin) function in the [vertex shader](./shaders/vertexshader.vert). The camera movement is possible with the `W`, `Q`, `A`, `S`, `D`, `E` buttons and with the mouse.

The application could be started with a settings screen, where the item colors, the number of the items, the size, and the camera options could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/04-mesh-deformer-with-camera/app.go
```

![Sample image app without settings](./sample/sample.png)

- with settings:

```
SETTINGS=on go run examples/04-mesh-deformer-with-camera/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
