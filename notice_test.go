package main

import (
	"fmt"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestParseRSS(t *testing.T) {
	p := gofeed.NewParser()
	feed, err := p.ParseURL("https://coolshell.cn/feed")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", feed)
}
