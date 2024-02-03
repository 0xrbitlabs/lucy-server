package handlers

import (
	"github.com/oklog/ulid/v2"
	"net/http"
	"server/internal/types"
	"server/internal/utils"
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

func (h *WebhookHandler) HandleButtonReplyEvent() {

}
