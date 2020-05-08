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

type Light interface {
	GetSpecular() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetPosition() mgl32.Vec3
}
type DirectionalLight interface {
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
}
type PointLight interface {
	GetPosition() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
}
type SpotLight interface {
	GetPosition() mgl32.Vec3
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
	GetCutoff() float32
}

type LightSource struct {
	LightSource         Light
	PositionUniformName string
	AmbientUniformName  string
	DiffuseUniformName  string
	SpecularUniformName string
}
type DirectionalLightSource struct {
	LightSource          DirectionalLight
	DirectionUniformName string
	AmbientUniformName   string
	DiffuseUniformName   string
	SpecularUniformName  string
}
type PointLightSource struct {
	LightSource              PointLight
	PositionUniformName      string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
}
type SpotLightSource struct {
	LightSource              SpotLight
	PositionUniformName      string
	DirectionUniformName     string
	AmbientUniformName       string
	DiffuseUniformName       string
	SpecularUniformName      string
	ConstantTermUniformName  string
	LinearTermUniformName    string
	QuadraticTermUniformName string
	CutoffUniformName        string
}

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
	shaderProgramId         uint32
	textures                []texture
	lightSource             LightSource
	directionalLightSources []DirectionalLightSource
	pointLightSources       []PointLightSource
	spotLightSources        []SpotLightSource
	viewPosition            mgl32.Vec3
	viewPositionUniformName string
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
		textures:        []texture{},
		lightSource: LightSource{
			PositionUniformName: "",
			AmbientUniformName:  "",
			DiffuseUniformName:  "",
			SpecularUniformName: "",
		},
		directionalLightSources: []DirectionalLightSource{},
		pointLightSources:       []PointLightSource{},
		spotLightSources:        []SpotLightSource{},

		viewPosition:            mgl32.Vec3{0, 0, 0},
		viewPositionUniformName: "",
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

// SetLightSource setups a light source.
// It takes a Light input that contains the model related info,
// and it also takes 4 strings, the uniform names that are used in the shader applications
// the 'PositionUniformName', 'AmbientUniformName', 'DiffuseUniformName', 'SpecularUniformName'
func (s *Shader) SetLightSource(lightSource Light, position, ambient, diffuse, specular string) {
	s.lightSource.LightSource = lightSource
	s.lightSource.PositionUniformName = position
	s.lightSource.AmbientUniformName = ambient
	s.lightSource.DiffuseUniformName = diffuse
	s.lightSource.SpecularUniformName = specular
}

// AddDirectionalLightSource sets up a directional light source.
// It takes a DirectionalLight input that contains the model related info,
// and it also takes a [4]string, with the uniform names that are used in the shader applications
// the 'DirectionUniformName', 'AmbientUniformName', 'DiffuseUniformName', 'SpecularUniformName'.
// They has to be in this order.
func (s *Shader) AddDirectionalLightSource(lightSource DirectionalLight, uniformNames [4]string) {
	var dSource DirectionalLightSource
	dSource.LightSource = lightSource
	dSource.DirectionUniformName = uniformNames[0]
	dSource.AmbientUniformName = uniformNames[1]
	dSource.DiffuseUniformName = uniformNames[2]
	dSource.SpecularUniformName = uniformNames[3]

	s.directionalLightSources = append(s.directionalLightSources, dSource)
}

// AddPointLightSource sets up a point light source. It takes a PointLight
// input that contains the model related info, and it also containt the uniform names in [7]string format.
// The order has to be the following: 'PositionUniformName', 'AmbientUniformName', 'DiffuseUniformName',
// 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName', 'QuadraticTermUniformName'.
func (s *Shader) AddPointLightSource(lightSource PointLight, uniformNames [7]string) {
	var pSource PointLightSource
	pSource.LightSource = lightSource
	pSource.PositionUniformName = uniformNames[0]
	pSource.AmbientUniformName = uniformNames[1]
	pSource.DiffuseUniformName = uniformNames[2]
	pSource.SpecularUniformName = uniformNames[3]
	pSource.ConstantTermUniformName = uniformNames[4]
	pSource.LinearTermUniformName = uniformNames[5]
	pSource.QuadraticTermUniformName = uniformNames[6]

	s.pointLightSources = append(s.pointLightSources, pSource)
}

// AddSpotLightSource sets up a spot light source. It takes a SpotLight input
// that contains the model related info, and it also contains the uniform names in [8]string format.
// The order has to be the following: 'PositionUniformName', 'DirectionUniformName', 'AmbientUniformName',
// 'DiffuseUniformName', 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName',
// 'QuadraticTermUniformName'.
func (s *Shader) AddSpotLightSource(lightSource SpotLight, uniformNames [8]string) {
	var sSource SpotLightSource
	sSource.LightSource = lightSource
	sSource.PositionUniformName = uniformNames[0]
	sSource.DirectionUniformName = uniformNames[1]
	sSource.AmbientUniformName = uniformNames[2]
	sSource.DiffuseUniformName = uniformNames[3]
	sSource.SpecularUniformName = uniformNames[4]
	sSource.ConstantTermUniformName = uniformNames[5]
	sSource.LinearTermUniformName = uniformNames[6]
	sSource.QuadraticTermUniformName = uniformNames[7]

	s.spotLightSources = append(s.spotLightSources, sSource)
}

func (s *Shader) SetViewPosition(position mgl32.Vec3, uniformName string) {
	s.viewPosition = position
	s.viewPositionUniformName = uniformName
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

// Setup light related uniforms.
func (s *Shader) lightHandler() {
	s.directionalLightHandler()
	s.pointLightHandler()
	s.spotLightHandler()
}

// Setup directional light related uniforms. It iterates over the directional sources
// and setups each uniform, where the name is not empty.
func (s *Shader) directionalLightHandler() {
	for _, source := range s.directionalLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
	}

}

// Setup point light relates uniforms. It iterates over the point light sources and sets
// up every uniform, where the name is not empty.
func (s *Shader) pointLightHandler() {
	for _, source := range s.pointLightSources {
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
	}
}

// Setup spot light related uniforms. It iterates over the spot light sources and sets up
// every uniform, where the name is not empty.
func (s *Shader) spotLightHandler() {
	for _, source := range s.spotLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
		if source.CutoffUniformName != "" {
			s.SetUniform1f(source.CutoffUniformName, source.LightSource.GetCutoff())
		}
	}
}

// DrawPoints is the draw functions for points
func (s *Shader) DrawPoints(numberOfPoints int32) {
	s.lightHandler()
	gl.DrawArrays(gl.POINTS, 0, numberOfPoints)
}

// DrawTriangles is the draw function for triangles
func (s *Shader) DrawTriangles(numberOfPoints int32) {
	for index, _ := range s.textures {
		s.textures[index].Bind(textureMap(index))
		gl.Uniform1i(s.getUniformLocation(s.textures[index].uniformName), int32(s.textures[index].texUnitId-gl.TEXTURE0))
	}
	s.lightHandler()
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
