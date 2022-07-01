package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jdxj/notice/api/router"
	_ "github.com/jdxj/notice/subscription/github"
	_ "github.com/jdxj/notice/subscription/rss"
)

func main() {
	router.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	router.Stop()
}
