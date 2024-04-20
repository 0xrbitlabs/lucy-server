package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"lucy/app_errors"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
)

func (s *Server) handleCategoryCreate(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Label       string `json:"label"`
		Description string `json:"description"`
		Enabled     bool   `json:"enabled"`
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
	if data.Label == "" || data.Description == "" {
		s.writeError(w, app_errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	newCategory := &types.Category{
		ID:          ulid.Make().String(),
		Label:       data.Label,
		Description: data.Description,
		Enabled:     data.Enabled,
	}
	err = s.Store.InsertCategory(*newCategory)
	if err != nil {
		if errors.Is(err, app_errors.ErrDuplicateResource) {
			s.writeError(w, app_errors.ErrConflict, http.StatusConflict)
			return
		}
		s.logger.Error(err.Error())
		s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
		return
	}
	s.writeData(w, http.StatusCreated, nil)
	return
}

func (s *Server) handleCategoryGetAll(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(*types.User)
	data, err := s.Store.GetCategories(ok)
	if err != nil {
		s.logger.Error(err.Error())
		s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
		return
	}
	s.writeData(w, http.StatusOK, map[string]interface{}{
		"categories": data,
	})
	return
}

func (s *Server) handleToggleCategory(w http.ResponseWriter, r *http.Request) {
}
