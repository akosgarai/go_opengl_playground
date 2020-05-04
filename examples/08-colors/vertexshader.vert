#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform vec3 lightColor;
void main()
{
    vSmoothColor = vec4(lightColor*vColor,1);
    gl_Position = projection * view * model * vec4(vVertex,1);
}
