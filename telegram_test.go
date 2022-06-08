package main

import (
	"testing"
)

func TestSendMessage(t *testing.T) {
	err := SendMessage("https://coolshell.cn/articles/22242.html")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
