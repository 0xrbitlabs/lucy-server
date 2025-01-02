package utils

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
