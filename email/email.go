package email

import (
	"fmt"
	"net/smtp"

	"github.com/jdxj/notice/config"

	"github.com/jordan-wright/email"
)

const (
	host = "smtp.qq.com"
	addr = host + ":587"
)

var (
	emailCfg *config.Email
)

func init() {
	cfg, err := config.GetEmail()
	if err != nil {
		// todo
		//panic(err)
	}
	emailCfg = cfg
}

func Send(subject string, data []byte, to ...string) error {
	e := email.NewEmail()

	e.Subject = subject
	e.Text = data

	e.From = fmt.Sprintf("notice <%s>", emailCfg.Addr)
	e.To = to

	return e.Send(addr, smtp.PlainAuth("", emailCfg.Addr, emailCfg.Token, host))
}

func SendSelf(subject, content string) error {
	self := emailCfg.Addr

	return Send(subject, []byte(content), self)
}
