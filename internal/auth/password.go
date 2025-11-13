package auth

import (
	"github.com/prawirdani/golang-restapi/pkg/errorsx"
	"golang.org/x/crypto/bcrypt"
)

var ErrWrongCredentials = errorsx.Unauthorized("Check your credentials")

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
