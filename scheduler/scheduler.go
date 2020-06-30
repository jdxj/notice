package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"

	"github.com/jdxj/notice/config"
)

var (
	scheduler = cron.New()
)

type Selected struct {
	Neo  bool
	Bili bool
	Ryf  bool
	Sf   bool
}

func Start(sel *Selected) error {
	// 获取配置 ----------------------------------
	sfCfg, err := config.GetSourceforge()
	if err != nil {
		return err
	}
	// -----------------------------------------

	// 注册任务 ------------------------------------------------------
	if sel.Neo {
		if err := addNeoTask(); err != nil {
			return err
		}
	}
	if sel.Sf {
		if err := addMultiSourceforgeTask(sfCfg); err != nil {
			return err
		}
	}
	if sel.Ryf {
		if err := addRuanYiFengTask(); err != nil {
			return err
		}
	}
	if sel.Bili {
		if err := addBiliBiliTask(); err != nil {
			return err
		}
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
