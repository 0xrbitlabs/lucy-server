package server

import (
	"encoding/json"
	"fmt"
	"lucy/app_errors"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
)

func (s *Server) CreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Type     string `json:"type"`
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
	count, err := s.Store.CheckForDuplicate(data.Phone)
	if err != nil {
		s.logger.Error(err.Error())
		s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
		return
	}
	if count != 0 {
		s.writeError(w, app_errors.ErrConflict, http.StatusConflict)
		return
	}
	hash, err := s.hash(data.Password)
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
	user := &types.User{
		ID:          ulid.Make().String(),
		Type:        data.Type,
		Password:    hash,
		PhoneNumber: data.Phone,
	}
	err = s.Store.InsertUser(user)
	if err != nil {
		s.logger.Error(err.Error())
		s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
		return
	}
	s.writeData(w, http.StatusCreated, nil)
	return
}
