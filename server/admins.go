package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"lucy/app_errors"
	"lucy/store"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
)

func (s *Server) handleAdminLogin() http.HandlerFunc {
	type loginPayload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &loginPayload{}
		err := json.NewDecoder(r.Body).Decode(payload)
		if err != nil {
			s.logger.Error(fmt.Sprintf(
				"Error while decoding request body: %s",
				err,
			))
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		admin, err := s.Store.GetAdmin(store.GetAdminFilter{Column: "username", Value: payload.UserName})
		if err != nil {
			if errors.Is(err, app_errors.ErrResourceNotFound) {
				s.writeError(w, app_errors.ErrBadRequest, http.StatusBadRequest)
				return
			}
			s.logger.Error(err.Error())
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		if !s.passwordIsCorrect(payload.Password, admin.Password) {
			s.writeError(w, app_errors.ErrBadRequest, http.StatusBadRequest)
			return
		}
		_, token, err := s.jwt.Encode(map[string]interface{}{
			"id":      admin.ID,
			"isAdmin": true,
		})
		if err != nil {
			s.writeError(w, app_errors.ErrTokenEncoding, http.StatusInternalServerError)
			return
		}
		s.writeData(w, http.StatusOK, map[string]interface{}{
			"token": token,
		})
		return
	}
}

func (s *Server) handleCreateAdmin() http.HandlerFunc {
	type payload struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		IsSuper  bool   `json:"is_super"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		admin, ok := r.Context().Value("user").(*types.Admin)
		if !ok || !admin.IsSuper {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		data := &payload{}
		err := json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			s.logger.Error(
				"Error while decoding request body: %s",
				err,
			)
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		duplicateCount, err := s.Store.CountAdminByUsername(data.UserName)
		if err != nil {
			s.logger.Error(err.Error())
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		if duplicateCount != 0 {
			s.writeError(w, app_errors.ErrConflict, http.StatusConflict)
			return
		}
		hash, err := s.hash(data.Password)
		if err != nil {
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		newAdmin := &types.Admin{
			ID:       ulid.Make().String(),
			Username: data.UserName,
			Password: hash,
		}
		err = s.Store.InsertAdmin(*newAdmin)
		if err != nil {
			s.logger.Error(err.Error())
			s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
			return
		}
		s.writeData(w, http.StatusCreated, nil)
		return
	}
}

func (s *Server) handleChangePassword() http.HandlerFunc {
	type payload struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		currentUser := r.Context().Value("user")
		switch currentUser.(type) {
		case *types.Admin:
			currUser := currentUser.(*types.Admin)
			data := &payload{}
			err := json.NewDecoder(r.Body).Decode(data)
			if err != nil {
				s.logger.Error(
					"Error while decoding request body: %s",
					err,
				)
				s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
				return
			}
			if !s.passwordIsCorrect(data.Old, currUser.Password) {
				s.writeError(w, app_errors.ErrBadRequest, http.StatusBadRequest)
				return
			}
			hash, err := s.hash(data.New)
			if err != nil {
				s.logger.Error(err.Error())
				s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
				return
			}
			err = s.Store.UpdateAdminInfo(currUser.Username, hash)
			if err != nil {
				s.logger.Error(err.Error())
				s.writeError(w, app_errors.ErrInternal, http.StatusInternalServerError)
				return
			}
			s.writeData(w, http.StatusOK, nil)
			return
		default:
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
