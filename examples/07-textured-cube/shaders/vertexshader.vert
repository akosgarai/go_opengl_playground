# version 410
layout (location = 0) in vec3 vVertex;
layout (location = 1) in vec3 vColor;
layout (location = 2) in vec2 vTexCoord;

out vec3 vSmoothColor;
out vec2 vSmoothTexCoord;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main()
{
    gl_Position = projection * view * model * vec4(vVertex,1);
    vSmoothColor = vColor;
    vSmoothTexCoord = vec2(vTexCoord.x, vTexCoord.y);
}
