package server

import (
	"encoding/json"
	"fmt"
	"lucy/app_errors"
	"net/http"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	data := &payload{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		s.logger.Error(fmt.Sprintf(
			"Error while decoding request body: %s",
			err,
		))
		s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {

}
