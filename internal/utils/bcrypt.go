package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return string(hash), nil
}

func PasswordMatches(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
