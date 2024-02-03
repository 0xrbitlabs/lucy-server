package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"server/internal/types"
	"server/internal/utils"

	"github.com/oklog/ulid/v2"
)

func (h *WebhookHandler) HandleTextEvent(w http.ResponseWriter, userContactInfo types.ContactSchema, message types.MessageSchema) {
	userPhone := message.From
	count, err := h.users.CountByPhoneNumber(userPhone)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	if count == 0 {
		newUser := &types.User{
			Id:          ulid.Make().String(),
			UserType:    "regular",
			PhoneNumber: userPhone,
			Password:    "",
			Name:        userContactInfo.Profile.Name,
			Description: "",
			Country:     "",
			Town:        "",
		}
		err := h.users.Insert(newUser)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		err = utils.SendTemplateMessage("welcome", userPhone, "fr_FR")
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	err = utils.SendMessageSingle(userPhone, message.Text.Body)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) HandleRegistrationRequest(w http.ResponseWriter, userContactInfo types.ContactSchema, userPhone string) {
	dbUser, err := h.users.GetByPhoneNumber(userPhone)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			_ = utils.SendErrorMessage(userPhone, h.logger)
			w.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Error(err.Error())
		_ = utils.SendErrorMessage(userPhone, h.logger)
		w.WriteHeader(http.StatusOK)
		return
	}
	if dbUser.UserType == "seller" {
		err := utils.SendMessageSingle(userPhone, "Vous etes deja enregistre en tant que vendeur sur notre plateforme :)")
		if err != nil {
			h.logger.Error(err.Error())
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	code := utils.GenerateVerificationCode()
	err = h.codes.Create(userPhone, code)
	if err != nil {
		h.logger.Error(err.Error())
		_ = utils.SendErrorMessage(userPhone, h.logger)
		w.WriteHeader(http.StatusOK)
		return
	}
	userRegistrationURL := fmt.Sprintf("%s/register?auth_code=%s", os.Getenv("WEB_APP_URL"), code)
	message := fmt.Sprintf(
    "Veuillez vous connecter sur notre plateforme pour terminer votre enregistrement en tant que vendeur en suivant ce lien [%s](%s)", 
    userRegistrationURL, 
    userRegistrationURL,
  )
	err = utils.SendMessageSingle(userPhone, message)
	if err != nil {
		_ = utils.SendErrorMessage(userPhone, h.logger)
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
