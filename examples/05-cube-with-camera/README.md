# Cube with camera

This application demonstrates the camera functionalities. The camera position is controlable with the `W` (forward), `A` (left), `S` (backward), `D` (right), `Q` (down), `E` (up) keys. The look direction is updatable with the mouse (boarder of the window, 8 direction). On the screen, there is a cube with different side colors.

The application could be started with a settings screen, where the color component of the items and the background color, and camera parameters could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/05-cube-with-camera/app.go
```

![Sample image app without settings](./sample/sample.png)

- with settings:

```
SETTINGS=on go run examples/05-cube-with-camera/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
