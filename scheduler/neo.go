package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/neoproxy"
)

func addNeoTask() {
	flow, err := neoproxy.NewFlow()
	if err != nil {
		logs.Error("new flow failed: %s", err)
		return
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
		logs.Error("add neo 'update dosage' and 'crawl news' task failed: %s", err)
		return
	}

	// ------------------------------------------------------------------------------
	// 每天 23:00 发送用量
	_, err = scheduler.AddFunc("0 23 * * *", func() {
		logs.Info("execute neo 'send dosage' task")

		flow.SendDosage()
	})
	if err != nil {
		logs.Error("add neo 'send dosage' task failed: %s", err)
		return
	}

	// ------------------------------------------------------------------------------
	// 如果发现 news 更新, 在5分钟内发送消息
	_, err = scheduler.AddFunc("*/5 * * * *", func() {
		logs.Info("execute neo 'send news' task")

		flow.SendLastNews()
	})
	if err != nil {
		logs.Error("add neo 'send news' task failed: %s", err)
		return
	}
}
