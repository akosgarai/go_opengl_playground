package shader

import (
	"os"
	"testing"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
func TestInitOpenGl(t *testing.T) {
	func() {
		defer func() {
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
	wrongFileContent := ""
	content, err := LoadShaderFromFile(wrongFileName)
	if err == nil {
		t.Error("Wrong filename should return error")
	}
	if content != wrongFileContent {
		t.Error("Wrong filename should return empty content")
	}
	goodFileName := "goodfile.name"
	goodFileContent := `This is My Content
	in my favorite file.
	`
	CreateFileWithContent(goodFileName, goodFileContent)
	defer DeleteFile(goodFileName)
	content, err = LoadShaderFromFile(goodFileName)
	if err != nil {
		t.Error("Good file shouldn't return error")
	}
	if content == goodFileContent {
		t.Error("Good file content should have the trailing '\\x00'")
	}
	if content != goodFileContent+"\x00" {
		t.Error("Good file content should be the same")
	}
}
func TestCompileShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	wrongFileContent := `This is My Content
	in my favorite file.
	` + "\x00"
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	_, err := CompileShader(wrongFileContent, gl.VERTEX_SHADER)
	if err == nil {
		t.Error("Compile should fail with wrong content.")
	}
	goodFileContent := `
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
	prog, err := CompileShader(goodFileContent, gl.VERTEX_SHADER)
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
		fragmentShaderFileName := "fragmentShader.frag"
		fragmentShaderRaw := `
	    #version 410
	    smooth in vec4 vSmoothColor;
	    layout(location=0) out vec4 vFragColor;
	    void main()
	    {
		vFragColor = vSmoothColor;
	    }
	    `
		vertexShaderFileName := "vertexShader.vert"
		vertexShaderRaw := `
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
		gl_PointSize = pointSize
		vSmoothColor = vec4(vColor,1);
	    }
	    `
		CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
		defer DeleteFile(fragmentShaderFileName)
		CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
		defer DeleteFile(vertexShaderFileName)
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(vertexShaderFileName, fragmentShaderFileName)
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
		fragmentShaderFileName := "fragmentShader.frag"
		fragmentShaderRaw := `
	    #version 410
	    smooth in vec4 vSmoothColor;
	    layout(location=0) out vec4 vFragColor;
	    void main()
	    {
		vFragColor = vSmoothColor
	    }
	    `
		vertexShaderFileName := "vertexShader.vert"
		vertexShaderRaw := `
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
		CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
		defer DeleteFile(fragmentShaderFileName)
		CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
		defer DeleteFile(vertexShaderFileName)
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(vertexShaderFileName, fragmentShaderFileName)
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
		fragmentShaderFileName := "fragmentShader.frag"
		vertexShaderFileName := "vertexShader.vert"
		vertexShaderRaw := `
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
		CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
		defer DeleteFile(vertexShaderFileName)
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(vertexShaderFileName, fragmentShaderFileName)
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
		fragmentShaderFileName := "fragmentShader.frag"
		fragmentShaderRaw := `
	    #version 410
	    smooth in vec4 vSmoothColor;
	    layout(location=0) out vec4 vFragColor;
	    void main()
	    {
		vFragColor = vSmoothColor;
	    }
	    `
		vertexShaderFileName := "vertexShader.vert"
		CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
		defer DeleteFile(fragmentShaderFileName)
		InitGlfw()
		defer glfw.Terminate()
		InitOpenGL()
		NewShader(vertexShaderFileName, fragmentShaderFileName)
	}()
}
func TestNewShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
}
func TestUse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
}
func TestSetUniformMat4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.SetUniformMat4("model", mgl32.Ident4())
}
func TestSetUniformMat3(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.SetUniformMat3("model", mgl32.Ident3())
}
func TestSetUniform3f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
    #version 410
    layout(location = 0) in vec3 vVertex;
    layout(location = 1) in vec3 vColor;
    smooth out vec4 vSmoothColor;
    uniform mat4 model;
    uniform mat4 view;
    uniform mat4 projection;
    uniform vec3 ambientColor;
    void main()
    {
	gl_Position = projection * view * model * vec4(vVertex,1);
	vSmoothColor = vec4(ambientColor*vColor,1);
    }
    `
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	shader.SetUniform3f("ambientColor", 1, 1, 1)
}
func TestSetUniform1f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
	if shader.shaderProgramId == 0 {
		t.Error("Invalid shader program id")
	}
	var valueToSet float32
	valueToSet = 20
	shader.SetUniform1f("pointSize", valueToSet)
}
func TestGetUniformLocation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	fragmentShaderFileName := "fragmentShader.frag"
	fragmentShaderRaw := `
    #version 410
    smooth in vec4 vSmoothColor;
    layout(location=0) out vec4 vFragColor;
    void main()
    {
	vFragColor = vSmoothColor;
    }
    `
	vertexShaderFileName := "vertexShader.vert"
	vertexShaderRaw := `
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
	CreateFileWithContent(fragmentShaderFileName, fragmentShaderRaw)
	defer DeleteFile(fragmentShaderFileName)
	CreateFileWithContent(vertexShaderFileName, vertexShaderRaw)
	defer DeleteFile(vertexShaderFileName)
	InitGlfw()
	defer glfw.Terminate()
	InitOpenGL()
	shader := NewShader(vertexShaderFileName, fragmentShaderFileName)
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
	t.Skip("Unimplemented")
}
func TestBindVertexArray(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestVertexAttribPointer(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestClose(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestDrawPoints(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestDrawTriangles(t *testing.T) {
	t.Skip("Unimplemented")
}
