package mailer

import (
	"bytes"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/log"
	"gopkg.in/gomail.v2"
)

type HeaderParams struct {
	To      []string
	Cc      []string
	Subject string
}

type Mailer struct {
	Templates *Templates
	dialer    *gomail.Dialer
	cfg       *config.SMTPConfig
}

func New(cfg *config.Config) *Mailer {
	dialer := gomail.NewDialer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.AuthEmail,
		cfg.SMTP.AuthPassword,
	)

	templates := parseTemplates()

	return &Mailer{
		dialer:    dialer,
		cfg:       &cfg.SMTP,
		Templates: templates,
	}
}

func (m *Mailer) Send(headerParams HeaderParams, body bytes.Buffer) error {
	mail := m.createHeader(headerParams)
	mail.SetBody("text/html", body.String())

	if err := m.dialer.DialAndSend(mail); err != nil {
		log.Error("Failed to send mail", err)
		return err
	}

	return nil
}

func (m *Mailer) createHeader(params HeaderParams) *gomail.Message {
	mail := gomail.NewMessage()

	mail.SetHeader("From", m.cfg.SenderName)
	mail.SetHeader("To", params.To...)
	mail.SetHeader("Cc", params.Cc...)
	mail.SetHeader("Subject", params.Subject)

	return mail
}
