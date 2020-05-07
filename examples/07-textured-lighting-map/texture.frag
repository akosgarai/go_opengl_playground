# version 410
out vec4 FragColor;

in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;

struct Light {
    vec3 position;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
};

struct Material {
    sampler2D diffuse;
    sampler2D specular;
    float shininess;
};

uniform Light light;
uniform Material material;

uniform vec3 viewPosition;

void main()
{

    // ambient componenet
    vec3 ambientColor = light.ambient * texture(material.diffuse, TexCoords).rgb;

    // diffuse component
    vec3 normalizedNormal = normalize(Normal);
    vec3 lightDirection = normalize(light.position - FragPos);
    float diff = max(dot(normalizedNormal, lightDirection), 0.0);
    vec3 diffuseColor = light.diffuse * diff * texture(material.diffuse, TexCoords).rgb;

    // specular component
    vec3 viewDirection = normalize(viewPosition - FragPos);
    vec3 reflectDir = reflect(-lightDirection, normalizedNormal);
    float spec = pow(max(dot(viewDirection, reflectDir), 0.0), material.shininess);
    vec3 specularColor = light.specular * spec * texture(material.specular, TexCoords).rgb;

    vec3 resultColor = ambientColor + diffuseColor + specularColor;
    FragColor = vec4(resultColor, 1.0);
}
