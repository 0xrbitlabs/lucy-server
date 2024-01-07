package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/store"
	"server/internal/types"
	"github.com/oklog/ulid/v2"
)

type UserHandler struct {
	users *store.Users
}

func NewUserHandler(users *store.Users) *UserHandler {
	return &UserHandler{
		users: users,
	}
}

func (h *UserHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	reqPayload := new(struct {
		PhoneNumber    string `json:"phone_number"`
		ProfilePicture string `json:"profile_picture"`
		FullName       string `json:"full_name"`
	})
	err := json.NewDecoder(r.Body).Decode(reqPayload)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser := &types.User{
		Id:             ulid.Make().String(),
		PhoneNumber:    reqPayload.PhoneNumber,
		ProfilePicture: reqPayload.ProfilePicture,
		FullName:       reqPayload.FullName,
	}
	err = h.users.Insert(newUser)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}
