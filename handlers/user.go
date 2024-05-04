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
}

type UserHandler struct {
	UserService
	types.ILogger
}

func (h UserHandler) handleCreateAdminAccount(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateAdminDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		WriteData(w, http.StatusBadRequest, errs)
		return
	}
	err = h.CreateAdminAccount(*payload)
	if err != nil {

	}
	WriteData(w, http.StatusCreated, nil)
	return
}
