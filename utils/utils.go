package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing password: %w", err)
	}
	return string(hash), nil
}

func PasswordMatchesHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

type err struct {
	Error string `json:"error"`
}

func WriteError(errStr string, w http.ResponseWriter) {
	bytes, _ := json.Marshal(err{Error: errStr})
	w.Write(bytes)
}

func WriteData(data interface{}, w http.ResponseWriter) {
	bytes, _ := json.Marshal(data)
	w.Write(bytes)
}

func GenerateRandomDigit() string {
	return fmt.Sprint(time.Now().Nanosecond())[:6]
}
