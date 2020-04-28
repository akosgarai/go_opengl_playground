package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/shader"
	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Rectangle struct {
	precision int
	vao       *vao.VAO
	shader    *shader.Shader

	color        mgl32.Vec3
	Points       [4]mgl32.Vec3
	material     *Material
	invertNormal bool
}

func NewRectangle(points [4]mgl32.Vec3, mat *Material, prec int, shader *shader.Shader) *Rectangle {
	return &Rectangle{
		precision:    prec,
		material:     mat,
		Points:       points,
		shader:       shader,
		vao:          vao.NewVAO(),
		invertNormal: false,
	}
}

// Log returns the string representation of this object.
func (r *Rectangle) Log() string {
	logString := "Rectangle:\n"
	logString += " - A : Vector{" + Vec3ToString(r.Points[0]) + "}\n"
	logString += " - B : Vector{" + Vec3ToString(r.Points[1]) + "}\n"
	logString += " - C : Vector{" + Vec3ToString(r.Points[2]) + "}\n"
	logString += " - D : Vector{" + Vec3ToString(r.Points[3]) + "}\n"
	logString += " - Normal : Vector{" + Vec3ToString(r.GetNormal()) + "}\n"
	logString += " - precision : " + IntegerToString(r.precision) + "\n"
	if r.invertNormal {
		logString += " - inverted normal vector\n"
	}
	logString += " - " + r.material.Log() + "\n"
	return logString
}

// SetPrecision updates the precision of the rectangle
func (r *Rectangle) SetPrecision(p int) {
	r.precision = p
}

// SetShader updates the shader of the rectangle.
func (r *Rectangle) SetShader(s *shader.Shader) {
	r.shader = s
}

// SetMaterial updates the material of the rectangle
func (r *Rectangle) SetMaterial(m *Material) {
	r.material = m
}
func (r *Rectangle) GetNormal() mgl32.Vec3 {
	v1 := r.Points[1].Sub(r.Points[0])
	v2 := r.Points[3].Sub(r.Points[0])
	cp := v1.Cross(v2).Normalize()
	if r.invertNormal {
		cp = cp.Mul(-1)
	}
	return cp
}
func (r *Rectangle) IsNormalInverted() bool {
	return r.invertNormal
}
func (r *Rectangle) SetInvertNormal(i bool) {
	r.invertNormal = i
}
func (r *Rectangle) SetupExternalVao(v *vao.VAO) *vao.VAO {
	verticalStep := (r.Points[1].Sub(r.Points[0])).Mul(1.0 / float32(r.precision))
	horisontalStep := (r.Points[3].Sub(r.Points[0])).Mul(1.0 / float32(r.precision))

	normal := r.GetNormal()
	for horisontalLoopIndex := 0; horisontalLoopIndex < r.precision; horisontalLoopIndex++ {
		for verticalLoopIndex := 0; verticalLoopIndex < r.precision; verticalLoopIndex++ {
			a := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			b := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			c := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			d := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			v.AppendVectors(a, normal)
			v.AppendVectors(b, normal)
			v.AppendVectors(c, normal)
			v.AppendVectors(a, normal)
			v.AppendVectors(c, normal)
			v.AppendVectors(d, normal)
		}
	}
	return v
}

func (r *Rectangle) setupVao() {
	r.vao.Clear()
	verticalStep := (r.Points[1].Sub(r.Points[0])).Mul(1.0 / float32(r.precision))
	horisontalStep := (r.Points[3].Sub(r.Points[0])).Mul(1.0 / float32(r.precision))

	normal := r.GetNormal()
	for horisontalLoopIndex := 0; horisontalLoopIndex < r.precision; horisontalLoopIndex++ {
		for verticalLoopIndex := 0; verticalLoopIndex < r.precision; verticalLoopIndex++ {
			a := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			b := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			c := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex + 1)))
			d := r.Points[0].Add(
				verticalStep.Mul(float32(verticalLoopIndex + 1))).Add(
				horisontalStep.Mul(float32(horisontalLoopIndex)))
			r.vao.AppendVectors(a, normal)
			r.vao.AppendVectors(b, normal)
			r.vao.AppendVectors(c, normal)
			r.vao.AppendVectors(a, normal)
			r.vao.AppendVectors(c, normal)
			r.vao.AppendVectors(d, normal)
		}
	}
}
func (r *Rectangle) buildVao() {
	// Create the vao object
	r.setupVao()

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(r.vao.Get()), gl.Ptr(r.vao.Get()), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	// setup color
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
}

func (r *Rectangle) Update(dt float64) {
}

// DrawWithLight is for drawing the rectangle to the screen. but with lightsource.
func (r *Rectangle) DrawWithLight(view, projection mgl32.Mat4, lightPos mgl32.Vec3) {
	r.shader.Use()

	r.shader.SetUniformMat4("view", view)
	r.shader.SetUniformMat4("projection", projection)
	M := mgl32.Ident4()
	r.shader.SetUniformMat4("model", M)

	// diffuse color
	diffCol := r.material.GetDiffuse()
	r.shader.SetUniform3f("diffuseColor", diffCol.X(), diffCol.Y(), diffCol.Z())
	// specular color
	specCol := r.material.GetSpecular()
	r.shader.SetUniform3f("specularColor", specCol.X(), specCol.Y(), specCol.Z())
	// shininess
	r.shader.SetUniform1f("shininess", r.material.GetShininess())
	// light position
	r.shader.SetUniform3f("lightPosition", lightPos.X(), lightPos.Y(), lightPos.Z())
	// normal matrix
	N := mgl32.Mat4Normal(M.Mul4(view))
	r.shader.SetUniformMat3("normal", N)

	r.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(r.vao.Get())/6))
}
