package mailer

import (
	"embed"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

type Templates struct {
	ResetPassword *template.Template
}

func parseTemplates() *Templates {
	return &Templates{
		ResetPassword: template.Must(
			template.ParseFS(templatesFS, "templates/reset-password-mail.html"),
		),
	}
}
