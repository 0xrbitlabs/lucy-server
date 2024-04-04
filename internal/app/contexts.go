package app

import (
	"server/internal/models"
	"server/internal/types"
)

func Context(message *models.InboundMessage) (types.MessageContext, map[string]interface{}) {
	return 0, map[string]interface{}{}
}
