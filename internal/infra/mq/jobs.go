package mq

const EmailResetPasswordJobKey = "email:reset-password"

type EmailResetPasswordJob struct {
	To        string `json:"to"`
	Name      string `json:"name"`
	ResetURL  string `json:"reset_url"`  // For password reset
	ExpiryMin int    `json:"expiry_min"` // Token expiry in minutes
}
