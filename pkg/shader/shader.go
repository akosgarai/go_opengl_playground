package shader

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

type texture struct {
	textureId   uint32
	targetId    uint32
	texUnitId   uint32
	uniformName string
}

func (t *texture) Bind(id uint32) {
	gl.ActiveTexture(id)
	gl.BindTexture(t.targetId, t.textureId)
	t.texUnitId = id
}
func (t *texture) IsBinded() bool {
	if t.texUnitId == 0 {
		return false
	}
	return true
}
func (t *texture) UnBind() {
	t.texUnitId = 0
	gl.BindTexture(t.targetId, t.texUnitId)
}

type Shader struct {
	shaderProgramId       uint32
	textures              []texture
	lightColor            mgl32.Vec3
	lightColorUniformName string
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
		shaderProgramId:       program,
		textures:              []texture{},
		lightColorUniformName: "",
		lightColor:            mgl32.Vec3{1, 1, 1},
	}
}
func (s *Shader) AddTexture(filePath string, wrapR, wrapS, minificationFilter, magnificationFilter int32, uniformName string) {
	img, err := loadImageFromFile(filePath)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("not 32 bit color")
	}

	tex := texture{
		textureId:   s.genTexture(),
		targetId:    gl.TEXTURE_2D,
		texUnitId:   0,
		uniformName: uniformName,
	}

	tex.Bind(gl.TEXTURE0)
	defer tex.UnBind()

	s.TexParameteri(gl.TEXTURE_WRAP_R, wrapR)
	s.TexParameteri(gl.TEXTURE_WRAP_S, wrapS)
	s.TexParameteri(gl.TEXTURE_MIN_FILTER, minificationFilter)
	s.TexParameteri(gl.TEXTURE_MAG_FILTER, magnificationFilter)

	gl.TexImage2D(tex.targetId, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, uint32(gl.UNSIGNED_BYTE), gl.Ptr(rgba.Pix))

	gl.GenerateMipmap(tex.textureId)

	s.textures = append(s.textures, tex)
}

func (s *Shader) genTexture() uint32 {
	var id uint32
	gl.GenTextures(1, &id)
	return id
}

func (s *Shader) UseLightColor(color mgl32.Vec3, uniformName string) {
	s.lightColor = color
	s.lightColorUniformName = uniformName
}
func (s *Shader) HasTexture() bool {
	if len(s.textures) > 0 {
		return true
	}
	return false
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

// BindBufferData gets a float array as an input, generates a buffer
// binds it as array buffer, and sets the input as buffer data.
func (s *Shader) BindBufferData(bufferData []float32) {
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(bufferData), gl.Ptr(bufferData), gl.STATIC_DRAW)
}

// BindVertexArray generates a vertex array and binds it.
func (s *Shader) BindVertexArray() {
	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
}

// VertexAttribPointer sets the pointer.
func (s *Shader) VertexAttribPointer(index uint32, size, stride int32, offset int) {
	gl.EnableVertexAttribArray(index)
	gl.VertexAttribPointer(index, size, gl.FLOAT, false, stride, gl.PtrOffset(offset))
}

// Close disables the vertexarraypointers and the vertex array.
func (s *Shader) Close(numOfVertexAttributes int) {
	for i := 0; i < numOfVertexAttributes; i++ {
		index := uint32(i)
		gl.DisableVertexAttribArray(index)
	}
	for index, _ := range s.textures {
		s.textures[index].UnBind()
	}
	gl.BindVertexArray(0)
}

// DrawPoints is the draw functions for points
func (s *Shader) DrawPoints(numberOfPoints int32) {
	if s.lightColorUniformName != "" {
		s.SetUniform3f(s.lightColorUniformName, s.lightColor.X(), s.lightColor.Y(), s.lightColor.Z())
	}
	gl.DrawArrays(gl.POINTS, 0, numberOfPoints)
}

// DrawTriangles is the draw function for triangles
func (s *Shader) DrawTriangles(numberOfPoints int32) {
	for index, _ := range s.textures {
		s.textures[index].Bind(textureMap(index))
		gl.Uniform1i(s.getUniformLocation(s.textures[index].uniformName), int32(s.textures[index].texUnitId-gl.TEXTURE0))
	}
	if s.lightColorUniformName != "" {
		s.SetUniform3f(s.lightColorUniformName, s.lightColor.X(), s.lightColor.Y(), s.lightColor.Z())
	}
	gl.DrawArrays(gl.TRIANGLES, 0, numberOfPoints)
}

// TexParameteri is a wrapper function for gl.TexParameteri
func (s *Shader) TexParameteri(pName uint32, param int32) {
	gl.TexParameteri(gl.TEXTURE_2D, pName, param)
}

// TextureBorderColor is a wrapper function for gl.glTexParameterfv with TEXTURE_BORDER_COLOR as pname.
func (s *Shader) TextureBorderColor(color [4]float32) {
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &color[0])
}
