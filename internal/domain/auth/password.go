// Package auth provides authentication and authorization functionality.
// This package handles user authentication through sessions, access tokens, and
// password management including secure hashing and password reset flows. It manages
// the complete authentication lifecycle from login through logout, including token
// generation, validation, and session management.
package auth

import (
	"github.com/prawirdani/golang-restapi/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrWrongCredentials = domain.ErrUnauthorized("Check your credentials")

func HashPassword(plain string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// VerifyPassword Verify / Decrypt user password
func VerifyPassword(plain, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return ErrWrongCredentials
	}
	return nil
}
