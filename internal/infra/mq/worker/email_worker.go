package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"

	"github.com/prawirdani/golang-restapi/internal/infra/mq"
	"github.com/prawirdani/golang-restapi/internal/mail"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type EmailWorker struct {
	mailer    *mail.Mailer
	templates *EmailTemplates
}

type EmailTemplates struct {
	ResetPassword *template.Template
}

func NewEmailWorker(mailer *mail.Mailer) (*EmailWorker, error) {
	// Parse templates once at worker startup
	templates := &EmailTemplates{
		ResetPassword: template.Must(
			template.ParseFiles("./templates/email/reset-password-mail.html"),
		),
	}

	return &EmailWorker{
		mailer:    mailer,
		templates: templates,
	}, nil
}

// HandlePasswordReset processes password reset email jobs
func (w *EmailWorker) HandlePasswordReset(ctx context.Context, payload json.RawMessage) error {
	var job mq.EmailResetPasswordJob
	if err := json.Unmarshal(payload, &job); err != nil {
		log.Error("Failed to unmarshal email payload", "err", err.Error())
		return err
	}

	// Execute template
	var buf bytes.Buffer
	if err := w.templates.ResetPassword.Execute(&buf, map[string]any{
		"Name":    job.Name,
		"Minutes": job.ExpiryMin,
		"URL":     job.ResetURL,
	}); err != nil {
		log.Error("Failed to execute email template", "err", err.Error())
		return err
	}

	// Send email
	err := w.mailer.Send(
		mail.HeaderParams{
			To:      []string{job.To},
			Subject: "Password Reset",
		},
		buf,
	)
	if err != nil {
		return err
	}

	log.Info("Password reset mail sent", "to", job.To)
	return nil
}
