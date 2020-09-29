#version 410
smooth in vec4 vSmoothColor;
layout(location=0) out vec4 vFragColor;
void main()
{
    vFragColor = vSmoothColor;
}
