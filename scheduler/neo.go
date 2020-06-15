package scheduler

import (
	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/app/neoproxy"
	"github.com/jdxj/notice/config"
)

func addNeoTask(neoCfg *config.Neo, emailCfg *config.Email) error {
	flow, err := neoproxy.NewFlow(neoCfg, emailCfg)
	if err != nil {
		return err
	}

	// ------------------------------------------------------------------------------
	// 每小时更新一次:
	//   - 用量
	//   - 新闻
	_, err = scheduler.AddFunc("0 * * * *", func() {
		logs.Info("execute neo 'update dosage' and 'crawl news'")

		flow.UpdateDosage()
		flow.CrawlLastNews()
	})
	if err != nil {
		return err
	}

	// ------------------------------------------------------------------------------
	// 每天 23:00 发送用量
	_, err = scheduler.AddFunc("0 23 * * *", func() {
		logs.Info("execute neo 'send dosage' task")

		flow.SendDosage()
	})
	if err != nil {
		return err
	}

	// ------------------------------------------------------------------------------
	// 如果发现 news 更新, 在5分钟内发送消息
	_, err = scheduler.AddFunc("*/5 * * * *", func() {
		logs.Info("execute neo 'send news' task")

		flow.SendLastNews()
	})
	if err != nil {
		return err
	}

	return nil
}
