package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	VertexShaderPointSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    const float pointSize = 20.0;
    void main()
    {
	gl_Position = vec4(vVertex,1);
	gl_PointSize = pointSize;
    }
    ` + "\x00"
	VertexShaderPointWithColorSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    const float pointSize = 20.0;
    smooth out vec4 vSmoothColor;
    void main()
    {
	gl_Position = vec4(vVertex,1);
	gl_PointSize = pointSize;
	vSmoothColor = vec4(vColor,1);
    }
    ` + "\x00"
	VertexShaderPointWithColorModelViewProjectionSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    const float pointSize = 20.0;
    smooth out vec4 vSmoothColor;
    uniform mat4 model;
    uniform mat4 view;
    uniform mat4 projection;
    void main()
    {
	gl_Position = projection * view * model * vec4(vVertex,1);
	gl_PointSize = pointSize;
	vSmoothColor = vec4(vColor,1);
    }
    ` + "\x00"
	VertexShaderDirectOutputSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    void main()
    {
	gl_Position = vec4(vVertex,1);
    }
    ` + "\x00"
	VertexShaderBasicSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    uniform mat4 MVP;
    void main()
    {
	vSmoothColor = vec4(vColor,1);
	gl_Position = MVP*vec4(vVertex,1);
    }
    ` + "\x00"
	VertexShaderModelViewProjectionSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    uniform mat4 model;
    uniform mat4 view;
    uniform mat4 projection;
    void main()
    {
	vSmoothColor = vec4(vColor,1);
	gl_Position = projection * view * model * vec4(vVertex,1);
    }
    ` + "\x00"
	VertexShaderDeformVertexPositionModelViewProjectionSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    uniform mat4 model;
    uniform mat4 view;
    uniform mat4 projection;
    uniform float time;
    const float amplitude = 0.125;
    const float frequency = 4;
    const float PI = 3.14159;
    void main()
    {
	vSmoothColor = vec4(vColor,1);
	float distance = length(vVertex);
	float z = amplitude*sin(-PI*distance*frequency+time);
	gl_Position = projection * view * model * vec4(vVertex.x, vVertex.y, z,1);
    }
    ` + "\x00"
	VertexShaderDeformVertexPositionSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    uniform mat4 MVP;
    uniform float time;
    const float amplitude = 0.125;
    const float frequency = 4;
    const float PI = 3.14159;
    void main()
    {
	vSmoothColor = vec4(vColor,1);
	float distance = length(vVertex);
	float z = amplitude*sin(-PI*distance*frequency+time);
	gl_Position = MVP*vec4(vVertex.x, vVertex.y, z,1);
    }
    ` + "\x00"
	VertexShaderWithLightSource = `
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
    ` + "\x00"
	FragmentShaderConstantSource = `
    #version 410
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vec4(1,1,1,1);
    }
    ` + "\x00"
	FragmentShaderBasicSource = `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    ` + "\x00"
	GeometryShaderQuadSubdivisionSource = `
    #version 410
    layout (triangles) in;
    layout (triangle_strip, max_vertices=256) out;
    uniform int sub_divisions;
    uniform mat4 MVP;
    void main() {
	vec4 v0 = gl_in[0].gl_Position;
	vec4 v1 = gl_in[1].gl_Position;
	vec4 v2 = gl_in[2].gl_Position;
	float dx = abs(v0.x-v2.x)/sub_divisions;
	float dz = abs(v0.z-v1.z)/sub_divisions;
	float x=v0.x;
	float z=v0.z;
	for(int j=0;j<sub_divisions*sub_divisions;j++) {
	    gl_Position = MVP * vec4(x,0,z,1);
	    EmitVertex();
	    gl_Position = MVP * vec4(x,0,z+dz,1);
	    EmitVertex();
	    gl_Position = MVP * vec4(x+dx,0,z,1);
	    EmitVertex();
	    gl_Position = MVP * vec4(x+dx,0,z+dz,1);
	    EmitVertex();
	    EndPrimitive();
	    x+=dx;
	    if((j+1) %sub_divisions == 0) {
		x=v0.x;
		z+=dz;
	    }
	}
    }
    ` + "\x00"
)

// LoadShaderFromFile takes a filepath string arguments.
// It loads the file and returns it as a '\x00' terminated string.
// It returns an error also.
func LoadShaderFromFile(path string) (string, error) {
	shaderCode, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	result := string(shaderCode) + "\x00"
	return result, nil
}
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

type Shader struct {
	shaderProgramId uint32
}

// NewShader returns a Shader. It's inputs are the filenames of the shaders.
// It reads the files and compiles them. The shaders are attached to the shader program.
func NewShader(vertexShaderPath, fragmentShaderPath string) *Shader {
	vertexShaderSource, err := LoadShaderFromFile(vertexShaderPath)
	if err != nil {
		panic(err)
	}
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := LoadShaderFromFile(fragmentShaderPath)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return &Shader{
		shaderProgramId: program,
	}
}

// Use is a wrapper for gl.UseProgram
func (s *Shader) Use() {
	gl.UseProgram(s.shaderProgramId)
}

// SetUniformMat4 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix4fv function
func (s *Shader) SetUniformMat4(uniformName string, mat mgl32.Mat4) {
	location := s.getUniformLocation(uniformName)
	gl.UniformMatrix4fv(location, 1, false, &mat[0])
}

// SetUniformMat3 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix3fv function
func (s *Shader) SetUniformMat3(uniformName string, mat mgl32.Mat3) {
	location := s.getUniformLocation(uniformName)
	gl.UniformMatrix3fv(location, 1, false, &mat[0])
}

// SetUniform3f gets an uniform name string and 3 float values as input and
// calls the gl.Uniform3f function
func (s *Shader) SetUniform3f(uniformName string, v1, v2, v3 float32) {
	location := s.getUniformLocation(uniformName)
	gl.Uniform3f(location, v1, v2, v3)
}

// SetUniform1f gets an uniform name string and a float value as input and
// calls the gl.Uniform1f function
func (s *Shader) SetUniform1f(uniformName string, v1 float32) {
	location := s.getUniformLocation(uniformName)
	gl.Uniform1f(location, v1)
}
func (s *Shader) getUniformLocation(uniformName string) int32 {
	return gl.GetUniformLocation(s.shaderProgramId, gl.Str(uniformName+"\x00"))
}
