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

type ContentFormat int

const (
	Text ContentFormat = iota
	Html
)

func NewSender(cfg *config.Email) *Sender {
	sender := &Sender{
		e:        email.NewEmail(),
		emailCfg: cfg,
	}
	return sender
}

type Sender struct {
	e        *email.Email
	emailCfg *config.Email
}

func (s *Sender) Send(format ContentFormat, subject string, data []byte, to ...string) error {
	emailCfg := s.emailCfg

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

func (s *Sender) SendTextSelf(subject, content string) error {
	emailCfg := s.emailCfg
	return s.Send(Text, subject, []byte(content), emailCfg.Addr)
}
func (s *Sender) SendTextSelfBytes(subject string, content []byte) error {
	emailCfg := s.emailCfg
	return s.Send(Text, subject, content, emailCfg.Addr)
}

func (s *Sender) SendHTMLSelf(subject, content string) error {
	emailCfg := s.emailCfg
	return s.Send(Html, subject, []byte(content), emailCfg.Addr)
}

func (s *Sender) SendHTMLSelfBytes(subject string, content []byte) error {
	emailCfg := s.emailCfg
	return s.Send(Html, subject, content, emailCfg.Addr)
}

func (s *Sender) SendTextOther(subject, content string, to ...string) error {
	return s.Send(Text, subject, []byte(content), to...)
}
