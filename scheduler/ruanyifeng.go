package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/ruanyifeng"
)

func addRuanYiFengTask() error {
	ryf := ruanyifeng.NewRuanYiFeng()

	_, err := scheduler.AddFunc("0 * * * *", func() {
		logs.Info("execute ruanyifeng 'update entry' task")

		ryf.UpdateEntry()
	})
	if err != nil {
		return err
	}

	_, err = scheduler.AddFunc("0 13 * * *", func() {
		logs.Info("execute ruanyifeng 'send update' task")

		ryf.SendUpdate()
	})
	return err
}
