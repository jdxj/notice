package scheduler

import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
)

const ()

var (
	scheduler = cron.New()
)

func Start() {
	logs.Info("start scheduled tasks")
	defer scheduler.Start()

	addNeoTask()
	addRaphaelTask()
}

func Stop() {
	logs.Info("stop scheduled tasks")
	scheduler.Stop()
}
