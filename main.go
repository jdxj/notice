package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/cmd"
)

func init() {
	logs.SetLogger(logs.AdapterFile,
		`{"filename":"notice.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":3,"color":true}`)
}

func main() {
	//sigs := make(chan os.Signal, 2)
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	//select {
	//case <-sigs:
	//	logs.Info("receive stop signal")
	//}
	if err := cmd.Execute(); err != nil {
		logs.Error("cmd execute failed: %s", err)
	}
}
