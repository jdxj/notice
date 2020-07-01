package scheduler

import (
	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/app/neoproxy"
)

func addNeoTask() error {
	flow := neoproxy.NewFlow()

	// ------------------------------------------------------------------------------
	_, err := scheduler.AddFunc("0 23 * * *", func() {
		logs.Info("execute neo 'update dosage' task")

		flow.NotifyDosage()
	})
	return err
}
