package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/bilibili"
)

func addBiliBiliTask() error {
	bili := bilibili.NewBiliBili()
	_, err := scheduler.AddFunc("0 8 * * *", func() {
		logs.Info("execute bilibili 'collect coins' task")

		bili.KeepCollect()
	})
	return err
}
