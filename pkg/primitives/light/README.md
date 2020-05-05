# Light

This package aims to gather the lightsource related parameter into a package.

A light source is described with it's `position` vector, `ambient`, `diffuse`, `specular` color vector.

## Extend the material math.

For using the light source parameters, we have to extend the math mentioned in the material descriptions.

```
ambientColorComponent = light.ambient * material.ambient
diffuseColorComponent = light.diffuse * (diffMult * material.diffuse)
specularColorComponent = light.specular * (specMult * material.specular)
```
