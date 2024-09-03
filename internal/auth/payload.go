package auth

// These are the payload structures for the access token and refresh token.
// Should align with the structure of user entity.

type AccessTokenPayload struct {
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	Type TokenType `json:"type"`
}

type RefreshTokenPayload struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Type TokenType `json:"type"`
}

type TokenPayload interface {
	AccessTokenPayload | RefreshTokenPayload
}
