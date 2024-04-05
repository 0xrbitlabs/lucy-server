package app

import (
	"fmt"
	"server/internal/models"
	"server/internal/pkg"
)

func HandleFirstMessage(message models.InboundMessage) {
	//Store user in database
	//Send welcome message
	err := pkg.SendTextMessage(message.From, "Hello, seems like this is your first message to me :)")
	if err != nil {
		fmt.Println(err.Error())
	}
}
