package handlers

import (
	"encoding/json"
	"lucy/types"
	"net/http"
)

func WriteData(w http.ResponseWriter, statusCode int, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	bytes, _ := json.Marshal(map[string]interface{}{
		"data": data,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func WriteError(w http.ResponseWriter, statusCode int, error types.ServiceError) {

}
