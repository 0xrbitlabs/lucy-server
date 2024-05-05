package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/types"
	"net/http"
)

type AuthService interface {
	Login(data dtos.LoginDTO) (*string, error)
}

type AuthHandler struct {
	service AuthService
	logger  types.ILogger
}

func (h AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.LoginDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		WriteData(w, http.StatusBadRequest, errs)
		return
	}
	token, err := h.service.Login(*payload)
	if err != nil {
		h.logger.Error(err.Error())
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, map[string]string{
		"token": *token,
	})
	return
}
