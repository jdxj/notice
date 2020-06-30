package scheduler

import (
	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/app/neoproxy"
)

func addNeoTask() error {
	flow := neoproxy.NewFlow()

	// ------------------------------------------------------------------------------
	_, err := scheduler.AddFunc("0 * * * *", func() {
		logs.Info("execute neo 'update dosage' and 'crawl news' task")

		flow.NotifyDosage()
	})
	return err
}
