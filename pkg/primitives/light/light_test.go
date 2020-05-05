package light

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultLightPosition     = mgl32.Vec3{0, 0, 0}
	DefaultAmbientComponent  = mgl32.Vec3{1, 1, 1}
	DefaultDiffuseComponent  = mgl32.Vec3{0.2, 0.2, 0.2}
	DefaultSpecularComponent = mgl32.Vec3{0.4, 0.4, 0.4}
)

func TestNew(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	if l.position != DefaultLightPosition {
		t.Error("Invalid light position")
	}
	if l.ambient != DefaultAmbientComponent {
		t.Error("Invalid ambient component")
	}
	if l.diffuse != DefaultDiffuseComponent {
		t.Error("Invalid diffuse component")
	}
	if l.specular != DefaultSpecularComponent {
		t.Error("Invalid specular component")
	}
}
func TestLog(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	log := l.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestGetAmbient(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	if l.GetAmbient() != DefaultAmbientComponent {
		t.Error("Invalid ambient color")
	}
}
func TestGetDiffuse(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	if l.GetDiffuse() != DefaultDiffuseComponent {
		t.Error("Invalid diffuse color")
	}
}
func TestGetSpecular(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	if l.GetSpecular() != DefaultSpecularComponent {
		t.Error("Invalid specular color")
	}
}
func TestGetPosition(t *testing.T) {
	l := New(DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent)
	if l.GetPosition() != DefaultLightPosition {
		t.Error("Invalid position vector")
	}
}
