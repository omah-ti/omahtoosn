package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

var (
	ErrPasswordTooLong = errors.New("password exceeds maximum length of 72 characters")
)

// menerima plain text password dan mengembalikan bcrypt hash nya
// tolak password yang > 72 byte
func HashPassword(plain string) (string, error) {
	if len(plain) > 72 {
		return "", ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// membandingkan plain text password dengan bcrypt hash yang tersimpan di DB
func VerifyPassword(hash, plain string) bool {
	if len(plain) > 72 {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
