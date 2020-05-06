#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;

struct Light {
    vec3 position;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform Light light;

void main()
{
    vSmoothColor = vec4(light.ambient*vColor,1);
    gl_Position = projection * view * model * vec4(vVertex,1);
}
