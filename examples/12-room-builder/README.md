# Room builder application

The purpose of this application is to demonstrate the RoomBuilder tool. The settings screen contains the parameters that you can set with the builder. The world screen contains a flat surface as ground on the xz plane. The camera position is updatable with the `q`, `w`, `e`, `a`, `s`, `d` keys. The door of the room could be opened and closed with key `c`. The `esc` key activates the menu screen.

How to run the application (if you are in the main directory):

```
go run examples/12-room-builder/app.go
```

The app starts the menu screen, where you can start the world screen with the current settings, activate the settings screen to update the settings, exit the application. If the world has been started, the menu screen changes, the continue activates the world screen, with the latest state, the restart option activates the world screen with the latest settings.

![Sample gif app without settings](./sample/sample.gif)
