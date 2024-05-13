package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/models"
	"lucy/types"
	"net/http"
)

type CategoryService interface {
	CreateCategory(dtos.CreateCategoryDTO) error
	GetAllCategories(*models.User) (*[]models.Category, error)
	ToggleEnabled(*dtos.ToggleEnabledDTO) error
}

type CategoryHandler struct {
	categoryService CategoryService
	logger          Logger
}

func NewCategoryHandler(categoryService CategoryService, logger Logger) CategoryHandler {
	return CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

func (h CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateCategoryDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := payload.Validate()
	if errs != nil {
		WriteBadReqErr(w, errs)
		return
	}
	err = h.categoryService.CreateCategory(*payload)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusCreated, nil)
}

func (h CategoryHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	currUser, ok := r.Context().Value("user").(*models.User)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	categories, err := h.categoryService.GetAllCategories(currUser)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, map[string]interface{}{
		"categories": *categories,
	})
}

func (h CategoryHandler) HandleToggleEnabled(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.ToggleEnabledDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.categoryService.ToggleEnabled(payload)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
    return
	}
	WriteData(w, http.StatusOK, nil)
}
