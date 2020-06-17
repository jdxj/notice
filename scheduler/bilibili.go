package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/bilibili"
	"github.com/jdxj/notice/config"
)

func addMultiBiliBiliTask(biliCfg *config.BiliBili, emailCfg *config.Email) error {
	if len(biliCfg.Cookies) <= 0 {
		logs.Warn("have no bilibili task")
	}

	for emailAddr, cookie := range biliCfg.Cookies {
		if err := addBiliBiliTask(emailAddr, cookie, emailCfg); err != nil {
			return err
		}
	}
	return nil
}

func addBiliBiliTask(emailAddr, cookie string, emailCfg *config.Email) error {
	bili := bilibili.NewBiliBili(emailAddr, cookie, emailCfg)

	_, err := scheduler.AddFunc("0 8 * * *", func() {
		logs.Info("execute bilibili 'collect coins' task")

		bili.CollectCoins()
	})
	if err != nil {
		return err
	}
	return nil
}
