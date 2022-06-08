package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := NewRSS()
	r.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	logger.Infof("stop")
	r.Stop()
}
