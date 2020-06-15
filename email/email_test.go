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
}
