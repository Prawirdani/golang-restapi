package mail

import (
	"bytes"

	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	"gopkg.in/gomail.v2"
)

type HeaderParams struct {
	To      []string
	Cc      []string
	Subject string
}

type Mailer struct {
	dialer *gomail.Dialer
	cfg    *config.SMTPConfig
	logger logging.Logger
}

func NewMailer(cfg *config.Config, logger logging.Logger) *Mailer {
	dialer := gomail.NewDialer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.AuthEmail,
		cfg.SMTP.AuthPassword,
	)

	return &Mailer{
		dialer: dialer,
		cfg:    &cfg.SMTP,
		logger: logger,
	}
}

func (m *Mailer) Send(headerParams HeaderParams, body bytes.Buffer) error {
	mail := m.createHeader(headerParams)
	mail.SetBody("text/html", body.String())

	if err := m.dialer.DialAndSend(mail); err != nil {
		m.logger.Error(logging.Service, "MailService.Send", err.Error())
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
