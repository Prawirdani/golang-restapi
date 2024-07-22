package common

type TokenType uint8

const (
	AccessToken TokenType = iota
	RefreshToken
)

// Label returns the string representation of the token type and used as token cookie name.
func (t TokenType) Label() string {
	return [...]string{"accessToken", "refreshToken"}[t]
}
