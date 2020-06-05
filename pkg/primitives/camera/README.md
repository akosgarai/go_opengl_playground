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

It updates the pitch and yaw values.

## GetBoundingObject

It returns the bounding object of the camera. Now it is defined as a sphere. The position is the current position of the camera. The radius is hardcoded to 0.1.

## BouncingObjectAfterWalk

BouncingObjectAfterWalk returns the bouncing object of the new position. It is used for colision detection. The step is forbidden if it leads to collision.

## BouncingObjectAfterStrafe

BouncingObjectAfterStrafe returns the bouncing object of the new position. It is used for colision detection. The step is forbidden if it leads to collision.

## BouncingObjectAfterLift

BouncingObjectAfterLift returns the bouncing object of the new position. It is used for colision detection. The step is forbidden if it leads to collision.
