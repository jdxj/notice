package module

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func NewEmailSender(config *EmailConfig) (*EmailSender, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid email config")
	}

	es := &EmailSender{
		config: config,
		email:  email.NewEmail(),
	}
	return es, nil
}

type EmailSender struct {
	config *EmailConfig
	email  *email.Email
}

func (es *EmailSender) SendMsg(subject, content string) error {
	e := es.email
	addr := es.config.Address
	token := es.config.Token

	e.From = fmt.Sprintf("notice <%s>", addr)
	e.To = []string{addr}

	e.Subject = subject
	e.Text = []byte(content)

	return e.Send("smtp.qq.com:587", smtp.PlainAuth("", addr, token, "smtp.qq.com"))
}
