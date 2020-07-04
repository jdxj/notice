package liantong

import (
	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"
)

func NewLianTong() *LianTong {
	return &LianTong{}
}

type LianTong struct {
}

func (lt *LianTong) NotifyPayment() {
	subject := "给联通手机号缴费"
	content := "缴费"

	// 暂时使用 neo 的配置
	neoCfg, err := config.GetNeo()
	if err != nil {
		logs.Error("get liantong config (neo) config failed: %s", err)
		return
	}

	if err := email.SendText(subject, content, neoCfg.User); err != nil {
		logs.Error("send email failed, subject: %s, content: %s",
			subject, content)
	}
}
