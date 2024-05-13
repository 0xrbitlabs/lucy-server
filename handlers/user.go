package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/models"
	"lucy/types"
	"net/http"
)

type UserService interface {
	CreateAdminAccount(dtos.CreateAdminDTO) error
	GetAllUsers() (*[]models.User, error)
	GetUserByID(id string) (*models.User, error)
	ChangePassword(dtos.ChangeUserPasswordDTO) error
}

type UserHandler struct {
	service UserService
	logger  Logger
}

func NewUserHandler(service UserService, logger Logger) UserHandler {
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
		WriteBadReqErr(w, errs)
		return
	}
	err = h.service.CreateAdminAccount(*payload)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusCreated, nil)
}

func (h UserHandler) HandleChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.ChangeUserPasswordDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		WriteBadReqErr(w, errs)
		return
	}
	err = h.service.ChangePassword(*payload)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, nil)
}
