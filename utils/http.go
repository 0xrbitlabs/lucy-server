package utils

import (
	"encoding/json"
	"net/http"
)

func write(w http.ResponseWriter, status int, data interface{}, key string) {
	w.WriteHeader(status)
  w.Header().Set("Content-Type", "application/json")
	if data == nil {
		data = map[string]interface{}{}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		key: data,
	})
}

func WriteData(w http.ResponseWriter, status int, data interface{}) {
	write(w, status, data, "data")
}

func WriteError(w http.ResponseWriter, status int, data interface{}) {
	write(w, status, data, "error")
}
