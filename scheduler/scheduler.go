package scheduler

import (
	"github.com/jdxj/notice/neoproxy"
	"github.com/robfig/cron/v3"
)

var (
	scheduler = cron.New()
)

func Start() {
	// 定时通知 neo 用量
	scheduler.AddFunc("0 23 * * *", neoproxy.NotifyDosage)
}

func Stop() {
	scheduler.Stop()
}
