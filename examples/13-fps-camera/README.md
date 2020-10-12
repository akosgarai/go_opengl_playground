# FPS camera application

This application demonstrates my FPS camera solution. It has to provide a terrain generator, a room editor and a streetlamp editor.

## Features

- Window setup with environment variables.
- Ground builder settings
- Camera builder settings
- Room builder settings
- Lamp builder settings

How to run the application (if you are in the main directory):

```
go run examples/13-fps-camera/app.go
```

The app starts the menu screen, where you can start the world screen with the current settings, activate the settings screen to update the settings, exit the application. If the world has been started, the menu screen changes, the continue activates the world screen, with the latest state, the restart option activates the world screen with the latest settings.

![Sample gif](./sample/sample.gif)
