package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/raphael"
)

func addRaphaelTask() {
	rap := raphael.NewRaphael()

	// ------------------------------------------------------------------------------
	// 每6个小时更新一次
	_, err := scheduler.AddFunc("0 */6 * * *", func() {
		logs.Info("execute raphael 'update item' task")

		rap.UpdateItem()
	})
	if err != nil {
		logs.Error("add raphael 'update item' task failed: %s", err)
		return
	}

	// ------------------------------------------------------------------------------
	// 如果发现 item 更新, 5分钟内发送
	_, err = scheduler.AddFunc("*/5 * * * *", func() {
		logs.Info("execute raphael 'send update' task")

		rap.SendUpdate()
	})
	if err != nil {
		logs.Error("add raphael 'send update' task failed: %s", err)
		return
	}
}
