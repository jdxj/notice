package ruanyifeng

import (
	"testing"

	"github.com/jdxj/notice/config"
)

func TestRuanYiFeng_UpdateEntry(t *testing.T) {
	emailCfg := &config.Email{
		Addr:  "985759262@qq.com",
		Token: "",
	}
	ryf := NewRuanYiFeng(emailCfg)
	ryf.UpdateEntry()
	ryf.SendUpdate()
}
