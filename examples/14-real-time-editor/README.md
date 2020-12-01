# Real time editor

This application is a kind of material editor. It displays a material sphere and a form, where i can set the parameters of the material. The UI items will be implemented in the engine later.

How to run the application (if you are in the main directory):

```
go run examples/14-real-time-editor/app.go
```

## UI items

- **Label**

The label is not an actual UI item. It holds the parameters for rendering a label. The surface mesh can be set as the surface of the label. It can be any surface of the ui items.

- **Button**

The button item is a model that contains 2 rectangles. One for the background, and one for the foreground. The foreground is the surface of label. The background changes the color on mouse hover. It handles the click events. On case of left mouse button release above the button it fires its clickCallback function.

- **TextInput**

It is the representation of the text input ui item. The idea about the item: Based on rectangles. The background rectangle is responsible for the hover event The foreground rectangle is split (horizontal) two half. The top half contains the label of the item. The bottom half contains the current value of the input. Currently its implementation is still WIP.

- **SliderInput**

It is the representation of the slider input ui item. The idea about the item: Based on rectangles. The background rectangle is responsible for the hover event The foreground rectangle is split (horizontal) two half. The top half contains the label of the item. The bottom half contains a slip bar and also a label for the value of the slip bar. The ratio between she slider and the label is 3/1. It handles the click events. If the left mouse button is clicked above the slider, and the mouse moves on the vertical axis, the slider follows the mouse movement, and the value of the slider input also changes.

![Sample gif material editor](./sample/sample.gif)
