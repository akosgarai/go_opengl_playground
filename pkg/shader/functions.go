package shader

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// InitOpenGL is for initializing the gl lib. It also prints out the gl version.
func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}
func textureMap(index int) uint32 {
	switch index {
	case 0:
		return gl.TEXTURE0
	case 1:
		return gl.TEXTURE1
	case 2:
		return gl.TEXTURE2
	}
	return 0
}

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
