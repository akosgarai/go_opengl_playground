# Static square app

This application draws a green square to the screen. It hasn't any mouse or keyboard control, just the image on the screen.

The application could be started with a settings screen, where the square and the background colors, the width could be set.

How to run the application (if you are in the main directory):

- without settings:

```
go run examples/02-static-square/app.go
```

![Sample image app without settings](./sample/sample.png)

- with settings:

```
SETTINGS=on go run examples/02-static-square/app.go
```

In settings mode, the `escape` key displays the menu screen, where the main screen could be started / continued / restarted with the latest settings. The settings page and exit function also available from the menu screen.
