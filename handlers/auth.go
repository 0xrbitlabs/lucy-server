package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/interfaces"
	"lucy/types"
	"net/http"
)

type AuthHandler struct {
	service interfaces.AuthService
	logger  interfaces.Logger
}

func NewAuthHandler(service interfaces.AuthService, logger interfaces.Logger) AuthHandler {
	return AuthHandler{
		service: service,
		logger:  logger,
	}
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
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, map[string]string{
		"token": *token,
	})
	return
}
