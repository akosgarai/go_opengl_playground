# version 410
out vec4 FragColor;

struct Tex {
    sampler2D diffuse;
    sampler2D specular;
};
struct Material {
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
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
uniform Tex tex;

uniform vec3 viewPosition;

// function prototypes
vec4 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir);

void main()
{
    vec3 norm = normalize(Normal);
    vec3 viewDirection = normalize(viewPosition - FragPos);

    vec4 result = vec4(0);
    // calculate Point lighting
    for (int i = 0; i < MAX_POINT_LIGHTS; i++) {
        result += CalculatePointLight(pointLight[i], norm, FragPos, viewDirection);
    }
    FragColor = result;
}

// calculates the color when using a point light.
vec4 CalculatePointLight(PointLight light, vec3 normal, vec3 fragPos, vec3 viewDir)
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
    vec3 ambient = light.ambient * material.ambient;
    vec3 diffuse = light.diffuse * material.diffuse * diff;
    vec3 specular = light.specular * material.specular * spec;
    //vec3 ambient = light.ambient * texture(tex.diffuse, TexCoords).rbg;
    //vec3 diffuse = light.diffuse * diff * texture(tex.diffuse, TexCoords).rbg;
    //vec3 specular = light.specular * spec * texture(tex.specular, TexCoords).rbg;
    ambient *= attenuation;
    diffuse *= attenuation;
    specular *= attenuation;
    vec4 ambientComponent = texture(tex.diffuse, TexCoords) * vec4(ambient, 1.0);
    vec4 diffuseComponent = texture(tex.diffuse, TexCoords) * vec4(diffuse, 1.0);
    vec4 specularComponent = texture(tex.specular, TexCoords) * vec4(specular, 1.0);
    return (ambientComponent + diffuseComponent + specularComponent);
}
