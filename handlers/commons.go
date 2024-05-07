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

func WriteError(w http.ResponseWriter, error types.ServiceError) {
	bytes, _ := json.Marshal(map[string]interface{}{
		"error": map[string]string{
			"code": error.ErrorCode,
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.StatusCode)
	w.Write(bytes)
}

func WriteBadReqErr(w http.ResponseWriter, errors map[string]string) {
	bytes, _ := json.Marshal(map[string]interface{}{
		"errors": errors,
	})
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusBadRequest)
  w.Write(bytes)
}
