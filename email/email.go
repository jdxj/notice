package email

import (
	"fmt"
	"net/smtp"
	"sync"

	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/config"

	"github.com/jordan-wright/email"
)

const (
	host = "smtp.qq.com"
	addr = host + ":587"
)

var (
	// mutex 保护 emailCfg 的初始化
	mutex    = sync.Mutex{}
	emailCfg *config.Email
)

// 应用单例模式
func getEmailCfg() (*config.Email, error) {
	if emailCfg != nil {
		return emailCfg, nil
	}

	var err error
	mutex.Lock()
	defer mutex.Unlock()

	if emailCfg != nil {
		return emailCfg, nil
	}

	emailCfg, err = config.GetEmail()
	return emailCfg, err
}

func Send(subject string, data []byte, to ...string) error {
	cfg, err := getEmailCfg()
	if err != nil {
		logs.Warn("send by email failed:\n\tsubject: %s\n\tdata: %s",
			subject, data)
		return err
	}

	e := email.NewEmail()
	e.Subject = subject
	e.Text = data
	e.From = fmt.Sprintf("notice <%s>", cfg.Addr)
	e.To = to

	return e.Send(addr, smtp.PlainAuth("", cfg.Addr, cfg.Token, host))
}

func SendSelf(subject, content string) error {
	self := "985759262@qq.com"

	return Send(subject, []byte(content), self)
}
