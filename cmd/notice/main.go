package main

import (
	"github.com/jdxj/notice/scheduler"
)

func main() {
	s := scheduler.New()
	s.Register()

	s.Run()
}
