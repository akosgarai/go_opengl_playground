# Menu screen

The purpose of this application is to demonstrate a menu screen, where we can manage the game state (start it, exit, setup stuff). It also should demonstrate the fonts (ttf) as texts. When the app starts, it sets the Menu screen and displays it. On case of hover, the button is highlighted. On case of click on the bottom below, the application exits. On case of clicking the button above, the clear color changes and the menu button disappeares. On case of pushing the `q` button, the menu appeares again, and the clear color changes back to blue.
The font solution is based on the following packages:

- [`github.com/nullboundary/glfont`](https://github.com/nullboundary/glfont) package.
- [`freetype`](https://godoc.org/github.com/golang/freetype) package.
- [`truetype`](https://godoc.org/github.com/golang/freetype/truetype) package.
