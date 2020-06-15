package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/config"

	"github.com/jdxj/notice/app/raphael"
)

func addRaphaelTask(emailCfg *config.Email) error {
	rap := raphael.NewRaphael(emailCfg)

	// ------------------------------------------------------------------------------
	// 每6个小时更新一次
	//   - ex rom
	//   - im kernel
	_, err := scheduler.AddFunc("0 */6 * * *", func() {
		logs.Info("execute raphael 'update item' task")

		rap.UpdateItem()
		rap.UpdateItemIm()
	})
	if err != nil {
		return err
	}

	// ------------------------------------------------------------------------------
	// 如果发现 item 更新, 5分钟内发送
	//   - ex rom
	//   - im kernel
	_, err = scheduler.AddFunc("*/5 * * * *", func() {
		logs.Info("execute raphael 'send update' task")

		rap.SendUpdate()
		rap.SendUpdateIm()
	})
	if err != nil {
		return err
	}

	return nil
}
