package material

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DefaultDiffuse   = mgl32.Vec3{0.6, 0.6, 0.6}
	DefaultSpecular  = mgl32.Vec3{0.7, 0.7, 0.7}
	DefaultShininess = float32(0.9)
)

func TestNew(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	if material.ambient != DefaultAmbient {
		t.Error("Invalid ambient color")
	}
	if material.diffuse != DefaultDiffuse {
		t.Error("Invalid diffuse color")
	}
	if material.specular != DefaultSpecular {
		t.Error("Invalid specular color")
	}
	if material.shininess != DefaultShininess {
		t.Errorf("Invalid shininess. Instead of '%f', we have '%f'.", DefaultShininess, material.shininess)
	}
}
func TestLog(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	log := material.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestGetAmbient(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	if material.GetAmbient() != DefaultAmbient {
		t.Error("Invalid ambient color")
	}
}
func TestGetDiffuse(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	if material.GetDiffuse() != DefaultDiffuse {
		t.Error("Invalid diffuse color")
	}
}
func TestGetSpecular(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	if material.GetSpecular() != DefaultSpecular {
		t.Error("Invalid specular color")
	}
}
func TestGetShininess(t *testing.T) {
	material := New(DefaultAmbient, DefaultDiffuse, DefaultSpecular, DefaultShininess)
	if material.GetShininess() != DefaultShininess {
		t.Errorf("Invalid shininess. Instead of '%f', we have '%f'.", DefaultShininess, material.shininess)
	}
}
