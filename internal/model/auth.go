package model

type LoginInput struct {
	Email     string `json:"email"    validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	Token             string `json:"token"               validate:"required"`
	NewPassword       string `json:"new_password"        validate:"required,min=8"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type ChangePasswordInput struct {
	Password          string
	NewPassword       string `json:"new_password"        validate:"required,min=8"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,eqfield=NewPassword"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
