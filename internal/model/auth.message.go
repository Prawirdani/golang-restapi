package model

import "time"

type ResetPasswordEmailMessage struct {
	To       string        `json:"to"`         // Recipient's email address
	Name     string        `json:"name"`       // Recipient's name
	ResetURL string        `json:"reset_url"`  // Link for resetting the password
	Expiry   time.Duration `json:"expiry_min"` // Expiration time of the reset token in minutes
}
