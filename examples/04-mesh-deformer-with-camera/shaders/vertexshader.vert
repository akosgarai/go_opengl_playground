#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float time;
const float amplitude = 0.125;
const float frequency = 4;
const float PI = 3.14159;
void main()
{
    vSmoothColor = vec4(vColor,1);
    float distance = length(vec3(vVertex.x, 0, vVertex.z));
    float y = amplitude*sin(-PI*distance*frequency+time);
    gl_Position = projection * view * model * vec4(vVertex.x, y, vVertex.z,1);
}
