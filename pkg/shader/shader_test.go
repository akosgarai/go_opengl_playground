package shader

import (
	"os"
	"runtime"
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	ValidFragmentShaderString = `
#version 410
smooth in vec4 vSmoothColor;
layout(location=0) out vec4 vFragColor;
void main()
{
    vFragColor = vSmoothColor;
}
    `
	ValidVertexShaderWithUniformsString = `
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
    `
	ValidVertexShaderWithUniformsStringWithTrailingChars = `
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
	ValidVertexShaderWithMat3String = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat3 normal;		//normal matrix
void main()
{
    vec3 vNormal = normalize(normal * vVertex);
    gl_Position = projection * view * model * vec4(vNormal,1);
    vSmoothColor = vec4(vColor,1);
}
    `
	ValidVertexShaderWithFloatUniformString = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float pointSize;
void main()
{
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = pointSize;
    vSmoothColor = vec4(vColor,1);
}
    `
	InvalidShaderString = `
    This string is not valid as
    a shader progrma.
    `
	InvalidShaderStringWithTrailingChars = `
    This string is not valid as
    a shader progrma.
    ` + "\x00"
	ValidTextureVertexShader = `
# version 410
layout (location = 0) in vec3 vVertex;
layout (location = 1) in vec3 vColor;
layout (location = 2) in vec2 vTexCoord;

out vec3 vSmoothColor;
out vec2 vSmoothTexCoord;

void main()
{
    gl_Position = vec4(vVertex, 1.0);
    vSmoothColor = vColor;
    vSmoothTexCoord = vec2(vTexCoord.x, vTexCoord.y);
}
    `
	ValidTextureFragmentShader = `
# version 410
out vec4 FragColor;
  
in vec3 vSmoothColor;
in vec2 vSmoothTexCoord;

uniform sampler2D textureOne;

void main()
{
    FragColor = texture(textureOne, vSmoothTexCoord) * vec4(vSmoothColor, 1.0);
}
    `
	EmptyString            = ""
	FragmentShaderFileName = "fragmentShader.frag"
	VertexShaderFileName   = "vertexShader.vert"
)

func CreateFileWithContent(name, content string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
func DeleteFile(name string) error {
	return os.Remove(name)
}
func InitGlfw() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(600, 600, "Test-window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
}
func NewTestShader(t *testing.T, validFragmentShaderContent, validVertexShaderContent string) *Shader {
	CreateFileWithContent(FragmentShaderFileName, validFragmentShaderContent)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, validVertexShaderContent)
	defer DeleteFile(VertexShaderFileName)
	InitGlfw()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
		t.Fail()
	}
	if len(shader.textures) != 0 {
		t.Error("Invalid shader length.")
		t.Fail()
	}
	return shader
}
func TestInitOpenGL(t *testing.T) {
	func() {
		defer func() {
			defer glfw.Terminate()
			if r := recover(); r != nil {
				t.Errorf("InitOpenGL shouldn't panicked!")
			}
		}()
		InitOpenGL()
	}()
}
func TestLoadShaderFromFile(t *testing.T) {
	// Create tmp file with a known content.
	// call function with
	// - bad filename, that doesn't exist, so that we should have an error.
	// - good filename, that exists and we know it's content
	wrongFileName := "badfile.name"
	content, err := LoadShaderFromFile(wrongFileName)
	if err == nil {
		t.Error("Wrong filename should return error")
	}
	if content != EmptyString {
		t.Errorf("Wrong filename should return empty content. We got: '%s'", content)
	}
	goodFileName := "goodfile.name"
	CreateFileWithContent(goodFileName, InvalidShaderString)
	defer DeleteFile(goodFileName)
	content, err = LoadShaderFromFile(goodFileName)
	if err != nil {
		t.Error("Good file shouldn't return error")
	}
	if content == InvalidShaderString {
		t.Error("Good file content should have the trailing '\\x00'")
	}
	if content != InvalidShaderString+"\x00" {
		t.Error("Good file content should be the same")
	}
}
func TestCompileShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	_, err := CompileShader(InvalidShaderStringWithTrailingChars, gl.VERTEX_SHADER)
	if err == nil {
		t.Error("Compile should fail with wrong content.")
	}
	prog, err := CompileShader(ValidVertexShaderWithUniformsStringWithTrailingChars, gl.VERTEX_SHADER)
	if err != nil {
		t.Error(err)
	}
	if prog == 0 {
		t.Error("Invalid shader program id")
	}
}
func TestNewShaderPanicOnVertexContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the invalid content!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
		defer DeleteFile(FragmentShaderFileName)
		CreateFileWithContent(VertexShaderFileName, InvalidShaderString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName)
	}()
}
func TestNewShaderPanicOnFragmentContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the invalid content!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, InvalidShaderString)
		defer DeleteFile(FragmentShaderFileName)
		CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName)
	}()
}
func TestNewShaderPanicOnFragmentFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the missing file!")
			}
		}()
		CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName)
	}()
}
func TestNewShaderPanicOnVertexFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the missing file!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
		defer DeleteFile(FragmentShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName)
	}()
}
func TestNewShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
}
func TestUse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
}
func TestSetUniformMat4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
	shader.SetUniformMat4("model", mgl32.Ident4())
}
func TestSetUniformMat3(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithMat3String)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
	shader.SetUniformMat3("model", mgl32.Ident3())
}
func TestSetUniform3f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithMat3String)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
	shader.SetUniform3f("ambientColor", 1, 1, 1)
}
func TestSetUniform1f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithFloatUniformString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	var valueToSet float32
	valueToSet = 20
	shader.Use()
	shader.SetUniform1f("pointSize", valueToSet)
}
func TestGetUniformLocation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName)
	shader.Use()
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	testData := []struct {
		Name     string
		Location int32
	}{
		{"model", 0},
		{"view", 2},
		{"projection", 1},
		{"notValidUniformName", -1},
	}
	for _, tt := range testData {
		location := shader.getUniformLocation(tt.Name)
		if location != tt.Location {
			t.Error("Invalid location identifier")
		}
	}
}
func TestBindBufferData(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1}
	shader.BindBufferData(bufferData)
}
func TestBindVertexArray(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1}
	shader.BindBufferData(bufferData)
	shader.BindVertexArray()
}
func TestVertexAttribPointer(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1}
	shader.BindBufferData(bufferData)
	shader.BindVertexArray()
	shader.VertexAttribPointer(uint32(0), int32(3), int32(6*4), 0)
}
func TestClose(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1}
	shader.BindBufferData(bufferData)
	shader.BindVertexArray()
	shader.VertexAttribPointer(uint32(0), int32(3), int32(6*4), 0)
	shader.VertexAttribPointer(uint32(1), int32(3), int32(6*4), 3*4)
	shader.Close(1)
}
func TestDrawPoints(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1, 5}
	shader.BindBufferData(bufferData)
	shader.BindVertexArray()
	shader.VertexAttribPointer(uint32(0), int32(3), int32(7*4), 0)
	shader.VertexAttribPointer(uint32(1), int32(3), int32(7*4), 3*4)
	shader.VertexAttribPointer(uint32(2), int32(1), int32(7*4), 6*4)
	shader.DrawPoints(1)
	shader.Close(2)
}
func TestDrawTriangles(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	bufferData := []float32{0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1}
	shader.BindBufferData(bufferData)
	shader.BindVertexArray()
	shader.VertexAttribPointer(uint32(0), int32(3), int32(7*4), 0)
	shader.VertexAttribPointer(uint32(1), int32(3), int32(7*4), 3*4)
	shader.DrawTriangles(1)
	shader.Close(1)
}
func TestTexParameteri(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestTextureBorderColor(t *testing.T) {
	t.Skip("Unimplemented")
}

func TestLoadImageFromFile(t *testing.T) {
	_, err := loadImageFromFile("this-image-does-not-exist.jpg")
	if err == nil {
		t.Error("Image load should be failed.")
	}
	_, err = loadImageFromFile("transparent-image-for-texture-testing.jpg")
	if err != nil {
		t.Error("Issue during load.")
	}
}
func TestAddTexture(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	shader.AddTexture("transparent-image-for-texture-testing.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	if len(shader.textures) != 1 {
		t.Error("Invalid shader length.")
	}
	shader.Close(2)
}
func TestAddTextureInvalidFilename(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("AddTexture should have panicked due to the missing file!")
			}
		}()
		runtime.LockOSThread()
		shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
		defer glfw.Terminate()
		shader.AddTexture("this-file-does-not-exist.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	}()
}
func TestTextureBind(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer glfw.Terminate()
	shader.AddTexture("transparent-image-for-texture-testing.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	shader.textures[0].Bind(gl.TEXTURE0)
	defer shader.textures[0].UnBind()

	if shader.textures[0].texUnitId != gl.TEXTURE0 {
		t.Error("Invalid texUnitId")
	}
	shader.Close(2)
}
func TestTextureIsBinded(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer shader.Close(2)
	defer glfw.Terminate()
	shader.AddTexture("transparent-image-for-texture-testing.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	shader.textures[0].Bind(gl.TEXTURE0)
	if !shader.textures[0].IsBinded() {
		t.Error("Texture should be binded")
	}
	shader.textures[0].UnBind()
	if shader.textures[0].IsBinded() {
		t.Log(shader.textures)
		t.Errorf("Texture shouldn't be binded. id: '%d'", shader.textures[0].texUnitId)
	}
}
func TestTextureUnbind(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer shader.Close(2)
	defer glfw.Terminate()
	shader.AddTexture("transparent-image-for-texture-testing.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	shader.textures[0].Bind(gl.TEXTURE0)
	if !shader.textures[0].IsBinded() {
		t.Error("Texture should be binded")
	}
	shader.textures[0].UnBind()
	if shader.textures[0].IsBinded() {
		t.Error("Texture shouldn't be binded")
	}
}
func TestHasTexture(t *testing.T) {
	runtime.LockOSThread()
	shader := NewTestShader(t, ValidTextureFragmentShader, ValidTextureVertexShader)
	defer shader.Close(2)
	defer glfw.Terminate()
	if shader.HasTexture() {
		t.Error("Shouldn't have texture")
	}
	shader.AddTexture("transparent-image-for-texture-testing.jpg", gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE, gl.LINEAR, gl.LINEAR, "textureOne")
	if !shader.HasTexture() {
		t.Error("it has texture")
	}
}
