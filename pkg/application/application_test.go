package application

import (
	"testing"
)

func TestNew(t *testing.T) {
	app := New()
	if len(app.items) != 0 {
		t.Error("Invalid application - items should be empty")
	}
	if app.cameraSet {
		t.Error("Camera shouldn't be set")
	}
}
func TestLog(t *testing.T) {
	app := New()
	log := app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
}
