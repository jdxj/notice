package ruanyifeng

import (
	"testing"
)

func TestRuanYiFeng_UpdateEntry(t *testing.T) {
	ryf := NewRuanYiFeng()
	ryf.UpdateEntry()
	ryf.SendUpdate()
}
