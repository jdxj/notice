package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jdxj/notice/logger"
	"github.com/jdxj/notice/model/github"
	"github.com/jdxj/notice/model/rss"
)

func main() {
	r := rss.NewRSS()
	r.Start()

	g := github.NewGithub()
	g.Start()

	logger.Infof("started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	r.Stop()
	g.Stop()

	logger.Infof("stopped")
}
