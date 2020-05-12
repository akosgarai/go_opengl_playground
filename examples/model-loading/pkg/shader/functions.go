package shader

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	wrapper "github.com/akosgarai/opengl_playground/examples/model-loading/pkg/glwrapper"
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

// LoadImageFromFile takes a filepath string argument.
// It loads the file, decodes it as PNG or jpg, and returns the image and error
func loadImageFromFile(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	return img, err

}
func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := wrapper.CreateShader(shaderType)

	csources, free := wrapper.Strs(source)
	wrapper.ShaderSource(shader, 1, csources, nil)
	free()
	wrapper.CompileShader(shader)

	var status int32
	wrapper.GetShaderiv(shader, wrapper.COMPILE_STATUS, &status)
	if status == wrapper.FALSE {
		var logLength int32
		wrapper.GetShaderiv(shader, wrapper.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		wrapper.GetShaderInfoLog(shader, logLength, nil, wrapper.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
