package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := NewRSS()
	r.Start()

	logger.Infof("started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	r.Stop()

	logger.Infof("stopped")
}
