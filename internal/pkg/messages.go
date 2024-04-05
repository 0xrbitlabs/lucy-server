package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendTextMessage(to, message string) error {
	payload := struct {
		From        string `json:"from"`
		To          string `json:"to"`
		MessageType string `json:"message_type"`
		Text        string `json:"text"`
		Channel     string `json:"channel"`
	}{
		From:        os.Getenv("PHONE_NUMBER"),
		To:          to,
		MessageType: "text",
		Text:        message,
		Channel:     "whatsapp",
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error while marshalling data: %w", err)
	}
	url := fmt.Sprintf("%s/messages", os.Getenv("VONAGE_URL"))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Error while constructing HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(os.Getenv("VONAGE_KEY"), os.Getenv("VONAGE_SECRET"))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending HTTP request: %w", err)
	}
	fmt.Println(response.Status)
	return nil
}
