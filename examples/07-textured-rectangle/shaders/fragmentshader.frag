# version 410
out vec4 FragColor;
  
in vec3 vSmoothColor;
in vec2 vSmoothTexCoord;

uniform sampler2D textureOne;

void main()
{
    FragColor = texture(textureOne, vSmoothTexCoord) * vec4(vSmoothColor, 1.0);
}
