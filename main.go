package main

import (
	"notice/module"
	"notice/neoproxy"
	"os"
	"os/signal"
	"syscall"

	"github.com/astaxie/beego/logs"
)

func main() {
	logs.SetLogger(logs.AdapterFile,
		`{"filename":"notice.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	config, err := module.ReadConfig()
	if err != nil {
		logs.Error("%s", err)
		return
	}

	// 启动 flow
	flow, err := neoproxy.NewFlow(config.NeoProxy, config.Email)
	if err != nil {
		logs.Error("%s", err)
		return
	}
	flow.Start()
	defer flow.Stop()

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		logs.Info("receive stop signal")
	}
}
