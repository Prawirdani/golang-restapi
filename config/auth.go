package config

import (
	"os"
	"time"
)

type Auth struct {
	JwtSecret                 string
	JwtTTL                    time.Duration
	SessionTTL                time.Duration
	ResetPasswordTTL          time.Duration
	ResetPasswordFormEndpoint string
}

func (t *Auth) Parse() error {
	t.JwtSecret = os.Getenv("AUTH_JWT_SECRET")
	t.ResetPasswordFormEndpoint = os.Getenv("AUTH_RESET_PASSWORD_FORM_ENDPOINT")

	if val := os.Getenv("AUTH_JWT_TTL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.JwtTTL = d
		}
	}
	if val := os.Getenv("AUTH_SESSION_TTL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.SessionTTL = d
		}
	}
	if val := os.Getenv("AUTH_RESET_PASSWORD_TTL"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.ResetPasswordTTL = d
		}
	}
	return nil
}
