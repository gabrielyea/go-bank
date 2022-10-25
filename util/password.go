package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}

	return string(hashedPass), nil
}

func IsValidPassword(input, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
}
