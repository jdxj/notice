package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"

	"github.com/jdxj/notice/config"
)

var (
	scheduler = cron.New()
)

func Start() error {
	cache, err := config.NewCache(config.CachePath)
	if err != nil {
		return err
	}
	defer cache.Close()

	// 获取配置
	neoCfg, err := cache.GetNeo()
	if err != nil {
		return err
	}

	emailCfg, err := cache.GetEmail()
	if err != nil {
		return err
	}

	// 注册任务
	if err := addNeoTask(neoCfg, emailCfg); err != nil {
		return err
	}
	if err := addRaphaelTask(emailCfg); err != nil {
		return err
	}

	logs.Info("start scheduled tasks")
	scheduler.Start()
	return nil
}

func Stop() {
	logs.Info("stop scheduled tasks")
	scheduler.Stop()
}
