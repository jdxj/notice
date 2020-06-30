package email

import (
	"fmt"
	"net/smtp"
	"os"

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
	cfg, err := config.GetEmail()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] can not get email config: %s\n\n", err)
		return
	}
	emailCfg = cfg
}

func send(format ContentFormat, subject string, data []byte, to ...string) error {
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

func SendText(subject, content string, to ...string) error {
	return send(Text, subject, []byte(content), to...)
}

func SendTextBytes(subject string, content []byte, to ...string) error {
	return send(Text, subject, content, to...)
}

func SendHTML(subject, content string, to ...string) error {
	return send(Html, subject, []byte(content), to...)
}

func SendHTMLBytes(subject string, content []byte, to ...string) error {
	return send(Html, subject, content, to...)
}
