package email

import (
	"fmt"
	"net/smtp"

	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/config"

	"github.com/jordan-wright/email"
)

const (
	host = "smtp.qq.com"
	addr = host + ":587"
)

type ContentFormat int

const (
	Text ContentFormat = iota
	Html
)

var (
	emailCfg *config.Email
)

func init() {
	logs.Info("check email cfg")

	cfg, err := config.GetEmail()
	if err != nil {
		logs.Error("get email cfg failed: %s", err)
		return
	}

	emailCfg = cfg
	logs.Info("check email cfg success")
}

func Send(format ContentFormat, subject string, data []byte, to ...string) error {
	if emailCfg == nil {
		return fmt.Errorf("email cfg not found:\n\tsubject: %s\n\tdata: %s",
			subject, data)
	}

	e := email.NewEmail()
	e.Subject = subject
	e.From = fmt.Sprintf("notice <%s>", emailCfg.Addr)
	e.To = to

	switch format {
	case Text:
		e.Text = data
	case Html:
		e.HTML = data
	}

	return e.Send(addr, smtp.PlainAuth("", emailCfg.Addr, emailCfg.Token, host))
}

func SendTextSelf(subject, content string) error {
	return Send(Text, subject, []byte(content), emailCfg.Addr)
}
func SendTextSelfBytes(subject string, content []byte) error {
	return Send(Text, subject, content, emailCfg.Addr)
}

func SendHTMLSelf(subject, content string) error {
	return Send(Html, subject, []byte(content), emailCfg.Addr)
}
