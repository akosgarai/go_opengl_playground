# Vector

It represents a 3 dimensional vector. It stores `X`, `Y`, `Z` coordinates. The coordinates are stored as `float64` values.

## Vector functions.

We can do different functions with the vectors.

### Length

It returns the `length` of the vector. The length is defined with the following formula: `Square root (X*X + Y*Y + Z*Z)`

### Dot aka SquaredLength

It returns the `length^2` of the vector. It's defined as the following formula: `X*X + Y*Y + Z*Z`
[Wikipedia](https://en.wikipedia.org/wiki/Dot_product)

### Add

The `add` function is defined as the following formula:
Let `V1`, `V2` vectors. The `V = V1 + V2` operation defines `V` as: `{ X = V1.X + V2.X, Y = V1.Y + V2.Y, Z = V1.Z + V2.Z }`

### Subtract

The substraction operand on vectors is defined with the following formula:
Let `V1`, `V2` vectors. The `V = V1 - V2` operation defines `V` as: `{ X = V1.X - V2.X, Y = V1.Y - V2.Y, Z = V1.Z - V2.Z }`

### Multiply
 
The multiply opearand returns a vector defined as the following formula:
Let `V1`, `V2` vectors. The `V = V1 * V2` operation defines `V` as: `{ X = V1.X * V2.X, Y = V1.Y * V2.Y, Z = V1.Z * V2.Z }`

### Divide

The divide opearand returns a vector defined as the following formula:
Let `V1`, `V2` vectors. The `V = V1 / V2` operation defines `V` as: `{ X = V1.X / V2.X, Y = V1.Y / V2.Y, Z = V1.Z / V2.Z }`

### AddScalar

This operand increments each coordinate of the vector with the given scalar value.
Let `V1` vector and `t` scalar. The `V = V1 + t` operation defines `V` as: `{ X = V1.X + t, Y = V1.Y + t, Z = V1.Z + t }`.
As we can see, `v.AddScalar(t) = v.Add(Vector{t,t,t})`.

### SubtractScalar

This operand decrements each coordinate of the vector with the given scalar value.
Let `V1` vector and `t` scalar. The `V = V1 - t` operation defines `V` as: `{ X = V1.X - t, Y = V1.Y - t, Z = V1.Z - t }`.
As we can see, `v.SubtractScalar(t) = v.Subtract(Vector{t,t,t})`.

### MultiplyScalar

This operand multiplies each coordinate of the vector with the given scalar value.
Let `V1` vector and `t` scalar. The `V = V1 * t` operation defines `V` as: `{ X = V1.X * t, Y = V1.Y * t, Z = V1.Z * t }`.
As we can see, `v.MultiplyScalar(t) = v.Multiply(Vector{t,t,t})`.

### DivideScalar

This operand divides each coordinate of the vector with the given scalar value.
Let `V1` vector and `t` scalar. The `V = V1 / t` operation defines `V` as: `{ X = V1.X / t, Y = V1.Y / t, Z = V1.Z / t }`.
As we can see, `v.DivideScalar(t) = v.Divide(Vector{t,t,t})`.

### Normalize

The normalize operand makes the vector unit. It divides each coordinate of the vector with it's length.

### Cross

The cross product returns the cross pruduct of the given  vectors.
[wikipedia](https://en.wikipedia.org/wiki/Cross_product)

### ToString

It returns the string representation of the vector. It is used for debugging, error logging.
