#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vNormal;

smooth out vec4 vSmoothColor;


struct Light {
    vec3 position;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

struct Material {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float shininess;
};

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

uniform Light light;
uniform Material material;

uniform vec3 viewPosition;
void main()
{
    // ambient componenet
    vec3 ambientColor = light.ambient * material.ambient;

    // diffuse component
    vec3 normalizedNormal = normalize(vNormal);
    vec3 lightDirection = normalize(light.position - vVertex);
    float diff = max(dot(normalizedNormal, lightDirection), 0.0);
    vec3 diffuseColor = light.diffuse * (diff * material.diffuse);

    // specular component
    vec3 viewDirection = normalize(viewPosition - vVertex);
    vec3 reflectDir = reflect(-lightDirection, normalizedNormal);
    float spec = pow(max(dot(viewDirection, reflectDir), 0.0), material.shininess);
    vec3 specularColor = light.specular * (spec * material.specular);

    vec3 resultColor = (ambientColor + diffuseColor + specularColor);
    vSmoothColor = vec4(resultColor,1);
    
    gl_Position = projection * view * model * vec4(vVertex,1);
}
