package service

import (
	"github.com/prawirdani/golang-restapi/config"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	"gopkg.in/gomail.v2"
)

type SendEmailParams struct {
	To      []string
	Cc      []string
	Subject string
}

type MailService struct {
	dialer *gomail.Dialer
	cfg    *config.SMTPConfig
	logger logging.Logger
}

func NewMailService(cfg *config.Config, logger logging.Logger) *MailService {
	dialer := gomail.NewDialer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.AuthEmail,
		cfg.SMTP.AuthPassword,
	)

	return &MailService{
		dialer: dialer,
		cfg:    &cfg.SMTP,
		logger: logger,
	}
}

// Send implements MailService.
func (m *MailService) Send(mailParams SendEmailParams, body string) error {
	mail := m.createHeader(mailParams)
	mail.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(mail); err != nil {
		m.logger.Error(logging.Service, "MailService.Send", err.Error())
		return err
	}

	return nil
}

func (m *MailService) createHeader(params SendEmailParams) *gomail.Message {
	mail := gomail.NewMessage()

	mail.SetHeader("From", m.cfg.SenderName)
	mail.SetHeader("To", params.To...)
	mail.SetHeader("Cc", params.Cc...)
	mail.SetHeader("Subject", params.Subject)

	return mail
}
