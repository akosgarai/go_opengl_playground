# version 410
out vec4 FragColor;
  
in vec2 TexCoord;

uniform sampler2D paper;

void main()
{
    FragColor = vec4(vec3(texture(paper, TexCoord)), 1.0);
}
