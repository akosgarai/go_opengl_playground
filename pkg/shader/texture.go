package shader

import (
	wrapper "github.com/akosgarai/opengl_playground/pkg/glwrapper"
)

type texture struct {
	textureId   uint32
	targetId    uint32
	texUnitId   uint32
	uniformName string
}

func (t *texture) Bind(id uint32) {
	wrapper.ActiveTexture(id)
	wrapper.BindTexture(t.targetId, t.textureId)
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
	wrapper.BindTexture(t.targetId, t.texUnitId)
}
