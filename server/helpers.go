package server

import (
	"encoding/json"
	"fmt"
	"lucy/app_errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) hash(plainText string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error while hashing: %w", err)
	}
	return string(hash), nil
}

func (s *Server) passwordIsCorrect(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Server) writeError(w http.ResponseWriter, errorCode app_errors.ErrorCode, statusCode int) {
	bytes, _ := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"code": errorCode,
		},
	})
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func (s *Server) writeData(w http.ResponseWriter, statusCode int, data interface{}) {
	bytes := []byte("")
	if data == nil {
		bytes, _ = json.Marshal(map[string]interface{}{
			"data": map[string]interface{}{},
		})
	} else {
		bytes, _ = json.Marshal(map[string]interface{}{
			"data": data,
		})
	}
	w.WriteHeader(statusCode)
	w.Write(bytes)
}
