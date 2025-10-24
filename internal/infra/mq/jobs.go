package mq

const EmailResetPasswordJobKey = "email:reset-password"

type EmailResetPasswordJob struct {
	Type      string `json:"type"` // "password_reset", "verification", "welcome"
	To        string `json:"to"`
	Name      string `json:"name"`
	ResetURL  string `json:"reset_url"`  // For password reset
	ExpiryMin int    `json:"expiry_min"` // Token expiry in minutes
}
