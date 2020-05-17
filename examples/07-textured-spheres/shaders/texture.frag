# version 410
out vec4 FragColor;

struct Material {
    sampler2D diffuse;
    sampler2D specular;
    float shininess;
};

struct PointLight {
    vec3 position;

    vec3 ambient;
    vec3 diffuse;
    vec3 specular;

    float constant;
    float linear;
    float quadratic;
};

in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoords;

#define MAX_POINT_LIGHTS 1

uniform PointLight pointLight[MAX_POINT_LIGHTS];
uniform Material material;

uniform vec3 viewPosition;

// function prototypes
vec3 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main()
{
    vec3 norm = normalize(Normal);
    vec3 viewDirection = normalize(viewPosition - FragPos);

    vec3 result = vec3(0);
    // calculate Point lighting
    for (int i = 0; i < MAX_POINT_LIGHTS; i++) {
        result += CalculatePointLight(pointLight[i], norm, FragPos, viewDirection);
    }
    FragColor = vec4(result, 1.0);
}

// calculates the color when using a point light.
vec3 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir)
{
    vec3 lightDir = normalize(light.position - fragPos);
    // diffuse shading
    float diff = max(dot(normal, lightDir), 0.0);
    // specular shading
    vec3 reflectDir = reflect(-lightDir, normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    // attenuation
    float distance = length(light.position - fragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));
    // combine results
    vec3 ambient = light.ambient * texture(material.diffuse, TexCoords).rbg;
    vec3 diffuse = light.diffuse * diff * texture(material.diffuse, TexCoords).rbg;
    vec3 specular = light.specular * spec * texture(material.specular, TexCoords).rbg;
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;
    return (ambient + diffuse + specular);
}
