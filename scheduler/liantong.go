package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/liantong"
)

func addLianTongTask() error {
	lt := liantong.NewLianTong()
	_, err := scheduler.AddFunc("0 13 2 */2 *", func() {
		logs.Info("execute liantong 'notify payment' task")

		lt.NotifyPayment()
	})
	return err
}
