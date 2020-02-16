package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
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
	VertexShaderDirectOutputSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    void main()
    {
	gl_Position = vec4(vVertex,1);
    }
    ` + "\x00"
	VertexShaderDirectOutputWithColorSource = `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    void main()
    {
	vSmoothColor = vec4(vColor,1);
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
