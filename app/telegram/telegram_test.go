package telegram

import (
	"fmt"
	"os"
	"testing"

	"github.com/jdxj/notice/config"
)

func TestMain(t *testing.M) {
	config.Init("../../config")
	Init()
	os.Exit(t.Run())
}

func TestSendMessage(t *testing.T) {
	fmt.Printf("%+v\n", config.GetTelegramBot())
	err := SendMessage("aha")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
