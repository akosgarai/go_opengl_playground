#version 410

layout(location=0) in vec3 vVertex;		//per-vertex position
layout(location=1) in vec3 vNormal;		//per-vertex normal

//uniforms
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat3 normal;		//normal matrix
uniform vec3 lightPosition;		//light position in object space
uniform vec3 diffuseColor;		//diffuse colour of object
uniform vec3 specularColor;		//specular colour of object
uniform float shininess;		//specular shininess

//shader outputs to the fragment shader
smooth out vec4 vSmoothColor;    //final diffuse colour to the fragment shader

//shader constant
const vec3 vEyeSpaceCameraPosition = vec3(0,0,0); //eye is at vec3(0,0,0) in eye space

void main()
{
    //multiply the object space light position with the modelview matrix
    //to get the eye space light position
    vec4 vEyeSpaceLightPosition = view * model * vec4(lightPosition,1);

    //multiply the object space vertex position with the modelview matrix
    //to get the eye space vertex position
    vec4 vEyeSpacePosition = view * model * vec4(vVertex,1);

    //multiply the object space normal with the normal matrix
    //to get the eye space normal
    vec3 vEyeSpaceNormal   = normalize(normal * vNormal);

    //get the light vector
    vec3 L = normalize(vEyeSpaceLightPosition.xyz-vEyeSpacePosition.xyz);
    //get the view vector
    vec3 V = normalize(vEyeSpaceCameraPosition.xyz-vEyeSpacePosition.xyz);
    //get the half way vector between light and view vectors
    vec3 H = normalize(L+V);

    //calculate the diffuse and specular components
    float diffuse = max(0, dot(vEyeSpaceNormal, L));
    float specular = max(0, pow(dot(vEyeSpaceNormal, H), shininess));

    //calculate the final colour by adding the diffuse and specular components
    vSmoothColor = diffuse*vec4(diffuseColor,1) + specular*vec4(specularColor, 1);

    //multiply the combiend modelview projection matrix with the object space vertex
    //position to get the clip space position
    gl_Position = projection * view * model * vec4(vVertex,1);
}

