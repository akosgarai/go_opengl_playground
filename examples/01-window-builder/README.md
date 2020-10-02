# Window builder application

This application demonstrates a window builder tool. It can create full screen window (width, height is the current monitor resolution), or a given size window. The title bar can also be managed.
The following environment variables are used for the setup:
- `WIDTH` - The width of the window in pixels. In full screen mode, this value is overwritten with the current monitor width.
- `HEIGHT` - The height of the window in pixels. In full screen mode, this value iverwritten with the current monitor height.
- `DECORATED` - If this env is set to "0", then the decoration of the window will be turned off (eg: no title bar).
- `TITLE` - The value of this env (if not empty) overwrites the default window title.
- `FULL` - If this env is set to "1", then the application will start in full screen mode.

How to run the application (if you are in the main directory):

With default setup:

```
go run examples/01-window-builder/app.go
```

Fullscreen mode:

```
FULL=1 go run examples/01-window-builder/app.go
```

600 * 800 window without decoration

```
WIDTH=600 HEIGHT=800 DECORATED=0 go run examples/01-window-builder/app.go
```
