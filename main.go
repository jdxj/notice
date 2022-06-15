package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := NewRSS()
	r.Start()

	g := NewGithub()
	g.Start()

	logger.Infof("started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	r.Stop()
	g.Stop()

	logger.Infof("stopped")
}
