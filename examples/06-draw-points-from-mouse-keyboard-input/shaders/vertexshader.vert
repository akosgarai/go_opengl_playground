#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
layout(location = 2) in float vSize;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
void main()
{
    vSmoothColor = vec4(vColor,1);
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = vSize;
}
