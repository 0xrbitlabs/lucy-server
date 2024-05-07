package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/interfaces"
	"lucy/types"
	"net/http"
)

type UserHandler struct {
	service interfaces.UserService
	logger  interfaces.Logger
}

func NewUserHandler(service interfaces.UserService, logger interfaces.Logger) UserHandler {
	return UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h UserHandler) HandleCreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateAdminDTO{}
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
	err = h.service.CreateAdminAccount(*payload)
	if err != nil {
		h.logger.Error(err.Error())
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusCreated, nil)
	return
}
