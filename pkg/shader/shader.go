package shader

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	shaderProgramId         uint32
	textures                []texture
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
	vertexShader, err := CompileShader(vertexShaderSource, wrapper.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := LoadShaderFromFile(fragmentShaderPath)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(fragmentShaderSource, wrapper.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program := wrapper.CreateProgram()
	wrapper.AttachShader(program, vertexShader)
	wrapper.AttachShader(program, fragmentShader)
	wrapper.LinkProgram(program)

	return &Shader{
		shaderProgramId:         program,
		textures:                []texture{},
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
		targetId:    wrapper.TEXTURE_2D,
		texUnitId:   0,
		uniformName: uniformName,
	}

	tex.Bind(wrapper.TEXTURE0)
	defer tex.UnBind()

	s.TexParameteri(wrapper.TEXTURE_WRAP_R, wrapR)
	s.TexParameteri(wrapper.TEXTURE_WRAP_S, wrapS)
	s.TexParameteri(wrapper.TEXTURE_MIN_FILTER, minificationFilter)
	s.TexParameteri(wrapper.TEXTURE_MAG_FILTER, magnificationFilter)

	wrapper.TexImage2D(tex.targetId, 0, wrapper.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, wrapper.RGBA, uint32(wrapper.UNSIGNED_BYTE), wrapper.Ptr(rgba.Pix))

	wrapper.GenerateMipmap(tex.textureId)

	s.textures = append(s.textures, tex)
}

func (s *Shader) genTexture() uint32 {
	var id uint32
	wrapper.GenTextures(1, &id)
	return id
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
// that contains the model related info, and it also contains the uniform names in [10]string format.
// The order has to be the following: 'PositionUniformName', 'DirectionUniformName', 'AmbientUniformName',
// 'DiffuseUniformName', 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName',
// 'QuadraticTermUniformName', 'CutoffUniformName'.
func (s *Shader) AddSpotLightSource(lightSource SpotLight, uniformNames [10]string) {
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
	sSource.CutoffUniformName = uniformNames[8]
	sSource.OuterCutoffUniformName = uniformNames[8]

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
	wrapper.UseProgram(s.shaderProgramId)
}

// SetUniformMat4 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix4fv function
func (s *Shader) SetUniformMat4(uniformName string, mat mgl32.Mat4) {
	location := wrapper.GetUniformLocation(s.shaderProgramId, uniformName)
	wrapper.UniformMatrix4fv(location, 1, false, &mat[0])
}

// SetUniformMat3 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix3fv function
func (s *Shader) SetUniformMat3(uniformName string, mat mgl32.Mat3) {
	location := wrapper.GetUniformLocation(s.shaderProgramId, uniformName)
	wrapper.UniformMatrix3fv(location, 1, false, &mat[0])
}

// SetUniform3f gets an uniform name string and 3 float values as input and
// calls the gl.Uniform3f function
func (s *Shader) SetUniform3f(uniformName string, v1, v2, v3 float32) {
	location := wrapper.GetUniformLocation(s.shaderProgramId, uniformName)
	wrapper.Uniform3f(location, v1, v2, v3)
}

// SetUniform1f gets an uniform name string and a float value as input and
// calls the gl.Uniform1f function
func (s *Shader) SetUniform1f(uniformName string, v1 float32) {
	location := wrapper.GetUniformLocation(s.shaderProgramId, uniformName)
	wrapper.Uniform1f(location, v1)
}

// BindBufferData gets a float array as an input, generates a buffer
// binds it as array buffer, and sets the input as buffer data.
func (s *Shader) BindBufferData(bufferData []float32) {
	vertexBufferObject := wrapper.GenBuffers()
	wrapper.BindBuffer(wrapper.ARRAY_BUFFER, vertexBufferObject)
	wrapper.ArrayBufferData(bufferData)
}

// BindVertexArray generates a vertex array and binds it.
func (s *Shader) BindVertexArray() {
	vertexArrayObject := wrapper.GenVertexArrays()
	wrapper.BindVertexArray(vertexArrayObject)
}

// VertexAttribPointer sets the pointer.
func (s *Shader) VertexAttribPointer(index uint32, size, stride int32, offset int) {
	wrapper.VertexAttribPointer(index, size, wrapper.FLOAT, false, stride, wrapper.PtrOffset(offset))
}

// Close disables the vertexarraypointers and the vertex array.
func (s *Shader) Close(numOfVertexAttributes int) {
	for i := 0; i < numOfVertexAttributes; i++ {
		index := uint32(i)
		wrapper.DisableVertexAttribArray(index)
	}
	for index, _ := range s.textures {
		s.textures[index].UnBind()
	}
	wrapper.BindVertexArray(0)
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
		if source.OuterCutoffUniformName != "" {
			s.SetUniform1f(source.OuterCutoffUniformName, source.LightSource.GetOuterCutoff())
		}
	}
}

// DrawPoints is the draw functions for points
func (s *Shader) DrawPoints(numberOfPoints int32) {
	s.lightHandler()
	wrapper.DrawArrays(wrapper.POINTS, 0, numberOfPoints)
}

// DrawTriangles is the draw function for triangles
func (s *Shader) DrawTriangles(numberOfPoints int32) {
	for index, _ := range s.textures {
		s.textures[index].Bind(textureMap(index))
		wrapper.Uniform1i(wrapper.GetUniformLocation(s.shaderProgramId, s.textures[index].uniformName), int32(s.textures[index].texUnitId-wrapper.TEXTURE0))
	}
	s.lightHandler()
	wrapper.DrawArrays(wrapper.TRIANGLES, 0, numberOfPoints)
}

// TexParameteri is a wrapper function for gl.TexParameteri
func (s *Shader) TexParameteri(pName uint32, param int32) {
	wrapper.TexParameteri(wrapper.TEXTURE_2D, pName, param)
}

// TextureBorderColor is a wrapper function for gl.glTexParameterfv with TEXTURE_BORDER_COLOR as pname.
func (s *Shader) TextureBorderColor(color [4]float32) {
	wrapper.TexParameterfv(wrapper.TEXTURE_2D, wrapper.TEXTURE_BORDER_COLOR, &color[0])
}
