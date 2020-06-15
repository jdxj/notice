package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/config"

	"github.com/jdxj/notice/app/sourceforge"
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
	rap := sourceforge.NewSourceforge(url, emailCfg)

	// ------------------------------------------------------------------------------
	// 每6个小时更新一次
	//   - ex rom
	//   - im kernel
	_, err := scheduler.AddFunc("0 */6 * * *", func() {
		logs.Info("execute sourceforge 'update item' task")

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
		logs.Info("execute sourceforge 'send update' task")

		rap.SendUpdate()
	})
	if err != nil {
		return err
	}

	return nil
}
