package app

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/rest/api/v2010"
	"os"
	"server/internal/models"
)

func HandleFirstMessage(message models.InboundMessage, twilioClient *twilio.RestClient) {
	//Store user in database
	//Send welcome message
	params := &openapi.CreateMessageParams{}
	params.SetFrom(os.Getenv("PHONE_NUMBER"))
  params.SetTo(fmt.Sprintf("whatsapp:+%s", message.From))
	params.SetBody("Hello\n")

	resp, err := twilioClient.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error while sending message ", err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println(string(response))
	}
}
