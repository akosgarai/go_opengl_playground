# Mesh

It contains everything that we need for drawing a stuff. Now i have 3 kind of meshes above the base one.

## Base mesh

Those parameters that are always needed for managing an object.

- **Verticies** - The verticies of the mesh.
- **vbo** - vertex buffer object.
- **vao** - vertex array object.
- **position** - The center position of the mesh. The model transformation is calculated based on this.
- **direction** - The mesh is moving to this direction. If this value is null vector, then the mes is not moving.
- **velocity** - The mesh is moving to the direction with this speed. if this value is null, then the mash is not moving.
- **angle** - The mesh is rotated with this angle. This value has to be radian. The model transformation is calculated based on this.
- **axis** - The mesh is rotated on this axis. If this vector is null, then the mesh is not rotated. The model transformation is calculated based on this.
- **scale** - The mesh is scaled by this vector. The model transformation is calculated based on this.

It has setter functions for the parameters, and getters for the necessary ones, also one for the model transformation matrix calculation and one for updating the state of the mesh.

## Textured mesh

It is a mesh extension for textured objects. It's parameter list is extended with the followings:

- **Indicies** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indicies. These are the values that i can pass to the buffer.
- **Textures** - The textures that are used for covering the drew mesh.
- **ebo** - The element buffer object identifier. the indicies are stored here.

It's `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewTexturedMesh` function returns a textured mesh.

## Material mesh

It is a mesh extension for material objects. It's parameter list is extended with the followings:

- **Indicies** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indicies. These are the values that i can pass to the buffer.
- **Material** - The material that is used for calculating the color of the mesh.
- **ebo** - The element buffer object identifier. the indicies are stored here.

It's `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewMaterialMesh` function returns a material mesh.

## Point mesh

It is a mesh extension for point objects. It's parameter list isn't extended, but it has an `Add` function for extending the mesh. It's necessary, because the `NewPointMesh` function returns an empty mesh, so that we have to fill it with the verticies one by one. It's `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with points, and then cleans up.

## Color mesh

It is a mesh extension for colored opbjects. The colors are set in the verticies. Its parameter list is exended with the followings:

- **Indicies** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indicies. These are the values that i can pass to the buffer.
- **ebo** - The element buffer object identifier. the indicies are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles and then cleans up. The `NewColorMesh` returns a color mesh.


## Textured color mesh

It is a mesh extension for textured object where the object color also counts. Its parameter list is extended with the followings:

- **Indicies** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indicies. These are the values that i can pass to the buffer.
- **Textures** - The textures that are used for covering the drew mesh.
- **ebo** - The element buffer object identifier. the indicies are stored here.

It's `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewTexturedColoredMesh` function returns a textured colored mesh.
