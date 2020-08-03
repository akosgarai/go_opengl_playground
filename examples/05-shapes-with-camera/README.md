# Shapes with camera

The purpose of this application is the demonstration of the currently existing shapes. The camera is controlable. We can change the position with the `q`, `w`, `e`, `a`, `s`, `d` keys. With the mouse, the camera direction can be changed.

The application could be started with a settings screen, where the color component, position, size of the items and the background color, and camera parameters could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/05-shapes-with-camera/app.go
```

![Sample gif app without settings](./sample/sample.gif)

- with settings:

```
SETTINGS=on go run examples/05-shapes-with-camera/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
