package shader

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

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
