#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vNormal;

smooth out vec4 vSmoothColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

uniform vec3 lightColor;
uniform vec3 lightPosition;

uniform float ambientStrength;
uniform float specularStrength;
uniform vec3 objectColor;
uniform vec3 viewPosition;
void main()
{
    vec3 ambientColor = ambientStrength * lightColor;

    vec3 viewDirection = normalize(viewPosition - vVertex);
    vec3 lightDirection = normalize(lightPosition - vVertex);
    vec3 normalizedNormal = normalize(vNormal);
    vec3 reflectDir = reflect(-lightDirection, normalizedNormal);

    float spec = pow(max(dot(viewDirection, reflectDir), 0.0), 32);
    vec3 specularColor = specularStrength * spec * lightColor;

    float diff = max(dot(normalizedNormal, lightDirection), 0.0);
    vec3 diffuseColor = diff * lightColor;

    vec3 resultColor = (ambientColor + diffuseColor + specularColor) * objectColor;
    vSmoothColor = vec4(resultColor,1);
    
    gl_Position = projection * view * model * vec4(vVertex,1);
}
