package auth

import "errors"

var (
	ErrInvalidTokenType = errors.New("Invalid token type")
	ErrMissingTokenType = errors.New("Missing token type in payload")
)

type TokenType uint8

const (
	AccessToken TokenType = iota
	RefreshToken
)

// Label returns the string representation of the token type and used as token cookie name.
func (t TokenType) Label() string {
	return [...]string{"accessToken", "refreshToken"}[t]
}

func (t TokenType) Validate() error {
	if t != AccessToken && t != RefreshToken {
		return ErrInvalidTokenType
	}
	return nil
}

// GetTokenType returns the token type from the given payload.
func GetTokenType(payload map[string]interface{}) (*TokenType, error) {
	if typ, ok := payload["typ"].(float64); ok {
		t := TokenType(typ)
		if err := t.Validate(); err != nil {
			return nil, err
		}
		return &t, nil

	}
	return nil, ErrMissingTokenType
}
