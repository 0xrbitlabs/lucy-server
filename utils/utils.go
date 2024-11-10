package utils

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomDigit() string {
	return fmt.Sprint(time.Now().Nanosecond())[:6]
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return string(hashed), nil
}

func HashMatches(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
