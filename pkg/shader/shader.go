package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
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
	FragmentShaderBasicSource = `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
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
