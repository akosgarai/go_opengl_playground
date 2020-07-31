#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float alpha;
void main()
{
    vSmoothColor = vec4(vColor,alpha);
    gl_Position = projection * view * model * vec4(vVertex,1);
}
