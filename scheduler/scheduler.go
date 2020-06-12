package scheduler

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/app/neoproxy"
	"github.com/robfig/cron/v3"
)

const ()

var (
	scheduler = cron.New()
)

func Start() {
	var err error
	// 避免忘记 start
	defer scheduler.Start()

	err = addNeoTask()
	if err != nil {
		logs.Error("add neo task failed: %s", err)
	}
}

func Stop() {
	scheduler.Stop()
}

func StartTest() {
	var err error

	err = addHelloTask()
	if err != nil {
		logs.Error("add hello task failed: %s", err)
	}

	scheduler.Start()
}

func addHelloTask() (err error) {
	_, err = scheduler.AddFunc("* * * * *", func() {
		fmt.Printf("%s - hello world\n", time.Now())
	})
	return
}

func addNeoTask() (err error) {
	flow, err := neoproxy.NewFlow()
	if err != nil {
		return err
	}

	//                         "0 23 * * *" 每天 23:00 运行
	_, err = scheduler.AddFunc("0 23 * * *", func() {
		neoproxy.NotifyDosage(flow)
	})
	return
}
