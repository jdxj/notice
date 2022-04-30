package telegram

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	err := SendMessage("aha")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
