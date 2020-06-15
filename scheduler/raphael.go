package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/config"

	"github.com/jdxj/notice/app/raphael"
)

func addMultiSourceforgeTask(sfCfg *config.Sourceforge, emailCfg *config.Email) error {
	if len(sfCfg.SubsAddr) <= 0 {
		logs.Warn("have no sourceforge task")
	}

	for _, v := range sfCfg.SubsAddr {
		if err := addSourceforgeTask(v, emailCfg); err != nil {
			return err
		}
	}
	return nil
}

func addSourceforgeTask(url string, emailCfg *config.Email) error {
	rap := raphael.NewRaphael(url, emailCfg)

	// ------------------------------------------------------------------------------
	// 每6个小时更新一次
	//   - ex rom
	//   - im kernel
	_, err := scheduler.AddFunc("0 */6 * * *", func() {
		logs.Info("execute raphael 'update item' task")

		rap.UpdateItem()
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
	})
	if err != nil {
		return err
	}

	return nil
}
