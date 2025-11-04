package mailer

import "text/template"

type EmailTemplates struct {
	ResetPassword *template.Template
}
