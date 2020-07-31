# Static button handler

This application draws a square and a triangle to the screen. The cube's position is changeable. It goes forward with the `W`, backward with the `S`, left with the `A` and right with the `D` button. The triangle also changes it's position, it moves to the opposite directions.

## Controls

- **W** Move the square forward, the triangle backward direction.
- **S** Move the square backward, the triangle forward direction.
- **A** Move the square left, the triangle right direction.
- **D** Move the square right, the triangle left direction.

The application could be started with a settings screen, where the color, position, scale of the items and the background color and the alpha value (for blending) could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/03-button-handler/app.go
```

![Sample image app without settings](./sample/sample.png)

- with settings:

```
SETTINGS=on go run examples/03-button-handler/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
