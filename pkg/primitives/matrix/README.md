# Matrix

It represents a 4 dimensional matrix (4\*4). It stores 16 single `float32` values.

## Matrix functions

The most necessary functions were implemented.

### GetMatrix

This function returns the points of the matrix in the given order. Let's define the matrix as

```
M = [
    11,12,13,14,
    21,22,23,24,
    31,32,33,34,
    41,42,43,44,
]
```

The function returns the points in the following order: `11,12,13,14,21,22,23,24,31,32,33,34,41,42,43,44`

### GetTransposeMatrix

This function returns the points of the transposed matrix of this matrix. Let's define the matrix as

```
M = [
    11,12,13,14,
    21,22,23,24,
    31,32,33,34,
    41,42,43,44,
]
```

The function returns the points in the following order: `11,21,31,41,12,22,32,42,13,23,33,43,14,24,34,44`

### TransposeMatrix

This function returns the transposed matrix. Let's define the matrix as

```
M = [
    11,12,13,14,
    21,22,23,24,
    31,32,33,34,
    41,42,43,44,
]
```

This function will return `MT` matrix as

```
MT = [
    11,21,31,41,
    12,22,32,42,
    13,23,33,43,
    14,24,34,44,
]
```

### Dot

The dot returns the matrix that is defined as the dot product of the given matrices. Let's define `M1` and `M2` matrices as

```
M1 = [
    11,12,13,14,
    21,22,23,24,
    31,32,33,34,
    41,42,43,44,
]
M2 = [
    11,12,13,14,
    21,22,23,24,
    31,32,33,34,
    41,42,43,44,
]
```

The `dot` product of the matrices `M = M1.Dot(M2)` is

```
M = [
    M1[0]*M2[0] + M1[4]*M2[1] + M1[8]*M2[2] + M1[12]*M2[3],
    M1[1]*M2[0] + M1[5]*M2[1] + M1[9]*M2[2] + M1[13]*M2[3],
    M1[2]*M2[0] + M1[6]*M2[1] + M1[10]*M2[2] + M1[14]*M2[3],
    M1[3]*M2[0] + M1[7]*M2[1] + M1[11]*M2[2] + M1[15]*M2[3],
    M1[0]*M2[4] + M1[4]*M2[5] + M1[8]*M2[6] + M1[12]*M2[7],
    M1[1]*M2[4] + M1[5]*M2[5] + M1[9]*M2[6] + M1[13]*M2[7],
    M1[2]*M2[4] + M1[6]*M2[5] + M1[10]*M2[6] + M1[14]*M2[7],
    M1[3]*M2[4] + M1[7]*M2[5] + M1[11]*M2[6] + M1[15]*M2[7],
    M1[0]*M2[8] + M1[4]*M2[9] + M1[8]*M2[10] + M1[12]*M2[11],
    M1[1]*M2[8] + M1[5]*M2[9] + M1[9]*M2[10] + M1[13]*M2[11],
    M1[2]*M2[8] + M1[6]*M2[9] + M1[10]*M2[10] + M1[14]*M2[11],
    M1[3]*M2[8] + M1[7]*M2[9] + M1[11]*M2[10] + M1[15]*M2[11],
    M1[0]*M2[12] + M1[4]*M2[13] + M1[8]*M2[14] + M1[12]*M2[15],
    M1[1]*M2[12] + M1[5]*M2[13] + M1[9]*M2[14] + M1[13]*M2[15],
    M1[2]*M2[12] + M1[6]*M2[13] + M1[10]*M2[14] + M1[14]*M2[15],
    M1[3]*M2[12] + M1[7]*M2[13] + M1[11]*M2[14] + M1[15]*M2[15],
]
```

### MultiVector

Let M1 matrix a transformation matrix. So it represents a transformation (translation, scale, rotate, ...). This function returns the transformed vector.
