package crypt

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(text string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash : %w", err)
	}
	return string(hashed), nil
}

func VerifyPassword(normal string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}
