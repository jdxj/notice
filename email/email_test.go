package email

import (
	"testing"

	"github.com/jdxj/notice/config"
)

func TestSetEmail(t *testing.T) {
	email := &config.Email{
		Addr:  "985759262@qq.com",
		Token: "",
	}
	if err := config.SetEmail(email); err != nil {
		t.Fatalf("%s\n", err)
	}
	config.Close()
}

func TestSend(t *testing.T) {
	err := Send("test", []byte("lll"), "985759262@qq.com")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}

func TestSendSelf(t *testing.T) {
	if err := SendSelf("abc", "def"); err != nil {
		t.Fatalf("%s\n", err)
	}
}
