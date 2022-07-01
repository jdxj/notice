package main

import (
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jdxj/notice/subscription/github"
	_ "github.com/jdxj/notice/subscription/rss"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
