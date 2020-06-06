# BoundingObject

This package contains the necessary setup for initializing a bounding object. It holds a type name (`AABB` or `Sphere`) and the params, that are necessary to know. For the `AABB`, we need to have a `width` a `height` and a `length` key in the params. For the `Sphere` We only need to have a `radius` key. The package provides getter functions for the parameters. The **Type()** returns the typeName, and the **Params()** returns the params map.
