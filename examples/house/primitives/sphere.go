package primitives

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/opengl_playground/pkg/vao"
)

type Sphere struct {
	center mgl32.Vec3
	radius float32
	color  mgl32.Vec3

	numOfRows       int
	numOfItemsInRow int
	vao             *vao.VAO
	shaderProgram   uint32
}

func NewSphere() *Sphere {
	return &Sphere{
		center:          mgl32.Vec3{0, 0, 0},
		radius:          1,
		color:           mgl32.Vec3{1, 1, 1},
		numOfRows:       20,
		numOfItemsInRow: 20,
		vao:             vao.NewVAO(),
	}
}

// Log returns the string representation of this object.
func (s *Sphere) Log() string {
	logString := "Sphere:\n"
	logString += " - center : Vector{" + Vec3ToString(s.center) + "}\n"
	logString += " - radius : " + Float32ToString(s.radius) + "\n"
	logString += " - color : Vector{" + Vec3ToString(s.color) + "}\n"
	logString += " - numOfRows : " + IntegerToString(s.numOfRows) + "\n"
	logString += " - numOfItemsInRow : " + IntegerToString(s.numOfItemsInRow) + "\n"
	return logString
}

// SetCenter updates the center of the sphere
func (s *Sphere) SetCenter(c mgl32.Vec3) {
	s.center = c
}

// GetCenter returns the center of the sphere
func (s *Sphere) GetCenter() mgl32.Vec3 {
	return s.center
}

// SetColor updates the color of the sphere
func (s *Sphere) SetColor(c mgl32.Vec3) {
	s.color = c
}

// GetColor returns the color of the sphere
func (s *Sphere) GetColor() mgl32.Vec3 {
	return s.color
}

// SetRadius updates the radius of the sphere
func (s *Sphere) SetRadius(r float32) {
	s.radius = r
}

// GetRadius returns the radius of the sphere
func (s *Sphere) GetRadius() float32 {
	return s.radius
}

// SetShaderProgram updates the shaderProgram of the sphere.
func (s *Sphere) SetShaderProgram(p uint32) {
	s.shaderProgram = p
}
func (s *Sphere) trianglePointsToVao(pa, pb, pc Point) {
	s.vao.AppendVectors(pa.Coordinate, pa.Color)
	s.vao.AppendVectors(pb.Coordinate, pb.Color)
	s.vao.AppendVectors(pc.Coordinate, pc.Color)
}
func (s *Sphere) setupVao() {
	s.vao.Clear()
	RefPoint := &mgl32.Vec3{0, 1, 0}
	step_Z := -mgl32.DegToRad(float32(360.0 / float32(s.numOfItemsInRow)))
	step_Y := -mgl32.DegToRad(float32(360.0 / float32(s.numOfRows)))
	for i := 0; i < s.numOfRows; i++ {
		i_Rotation := mgl32.HomogRotate3DZ(float32(i) * step_Z).Transpose()
		i1_Rotation := mgl32.HomogRotate3DZ(float32(i+1) * step_Z).Transpose()
		for j := 0; j < s.numOfItemsInRow; j++ {
			j1_Rotation := mgl32.HomogRotate3DY(float32(j+1) * step_Y).Transpose()
			j_Rotation := mgl32.HomogRotate3DY(float32(j) * step_Y).Transpose()
			if i == 0 {
				p1 := Point{*RefPoint, s.color}
				p2 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i1_Rotation)), s.color}
				p3 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i1_Rotation)), s.color}
				s.trianglePointsToVao(p1, p2, p3)
			} else {
				p1 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i_Rotation)), s.color}
				p2 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i_Rotation)), s.color}
				p3 := Point{mgl32.TransformCoordinate(*RefPoint, j1_Rotation.Mul4(i1_Rotation)), s.color}
				p4 := Point{mgl32.TransformCoordinate(*RefPoint, j_Rotation.Mul4(i1_Rotation)), s.color}
				s.trianglePointsToVao(p1, p2, p3)
				s.trianglePointsToVao(p1, p3, p4)
			}
		}
	}
}
func (s *Sphere) buildVao() {
	s.setupVao()
	gl.UseProgram(s.shaderProgram)

	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	// a 32-bit float has 4 bytes, so we are saying the size of the buffer,
	// in bytes, is 4 times the number of points
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(s.vao.Get()), gl.Ptr(s.vao.Get()), gl.STATIC_DRAW)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	// setup points
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// setup color
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 4*6, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
}
func (s *Sphere) Draw() {
	s.buildVao()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(s.vao.Get())/6))
}
func (s *Sphere) DrawWithUniforms(view, projection mgl32.Mat4) {
	gl.UseProgram(s.shaderProgram)

	viewLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLocation, 1, false, &view[0])
	projectionLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionLocation, 1, false, &projection[0])

	modelLocation := gl.GetUniformLocation(s.shaderProgram, gl.Str("model\x00"))

	M := mgl32.Translate3D(
		s.center.X(),
		s.center.Y(),
		s.center.Z()).Mul4(mgl32.Scale3D(
		s.radius,
		s.radius,
		s.radius))
	gl.UniformMatrix4fv(modelLocation, 1, false, &M[0])

	s.Draw()
}
func (s *Sphere) Update(delta float64) {
}
