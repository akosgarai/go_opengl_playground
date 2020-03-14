# Camera

It represents the camera or eye. We see our model world from the camera's point of view.

## NewCamera

Returns a new camera with the given setup. `position` - the camera or eye position, `worldUp` - the up direction in the world coordinate system, `yaw - the rotation in z` `pitch - the rotation in y`

## Walk

It updates the position (forward, back directions).

## Strafe

It updates the position (left, right directions).

## Lift

It updates the position (up, down directions).

## SetupProjection

It sets the projection related variables. `fov` - field of view, `aspectRation` - windowWidth/windowHeight, `near` - near clip plane, `far` - far clip plane

## GetProjectionMatrix

It returns the projectionMatrix of the camera. It setups a perspective transformation.

## GetViewMatrix

It gets the matrix to transform from world coordinates to this camera's coordinates. It returns the viewMatrix of the camera.

## UpdateDirection

Itupdates the pitch and yaw values.
