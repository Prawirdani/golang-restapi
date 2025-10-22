package auth

import (
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrWrongCredentials = errors.Unauthorized("Check your credentials")

func HashPassword(plain string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// Verify / Decrypt user password
func VerifyPassword(plain, hashed string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return ErrWrongCredentials
	}
	return nil
}
