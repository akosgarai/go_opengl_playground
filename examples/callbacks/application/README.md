# Application package

The purpose of this package to implement the aplication related functions in this package, and the plan is to simplify the main function in the example applications. The common functions are placed in this package, so that i don't need to reimplement them in every application. The main goal is to refacto the whole example directory with this new solution.

## Functions

In the first iteration the most obvious getter and setter functions were implemented (window, program, camera, keymap).
Now i'm thinking about the drawable items. To make it general, i will define a drawable interface, and every drawable object will implement it.
