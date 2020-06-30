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
	ds := config.DataStorage

	// 获取配置 ----------------------------------
	neoCfg, err := ds.GetNeo()
	if err != nil {
		return err
	}
	emailCfg, err := ds.GetEmail()
	if err != nil {
		return err
	}
	sfCfg, err := ds.GetSourceforge()
	if err != nil {
		return err
	}
	biliCfg, err := ds.GetBiliBili()
	if err != nil {
		return err
	}
	// -----------------------------------------

	// 注册任务 ------------------------------------------------------
	if err := addNeoTask(neoCfg, emailCfg); err != nil {
		return err
	}
	if err := addMultiSourceforgeTask(sfCfg, emailCfg); err != nil {
		return err
	}
	if err := addRuanYiFengTask(emailCfg); err != nil {
		return err
	}
	if err := addMultiBiliBiliTask(biliCfg, emailCfg); err != nil {
		return err
	}
	// -------------------------------------------------------------

	logs.Info("start scheduled tasks")
	scheduler.Start()
	return nil
}

func Stop() {
	logs.Info("stop scheduled tasks")
	scheduler.Stop()
}
