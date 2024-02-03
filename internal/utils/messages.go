package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendMessageSingle(to, content string) (int, error) {
	reqBody, err := json.Marshal(struct {
		MessagingProduct string `json:"messaging_product"`
		RecipientType    string `json:"recipient_type"`
		To               string `json:"to"`
		Text             struct {
			PreviewURL bool   `json:"preview_url"`
			Body       string `json:"body"`
		} `json:"text"`
	}{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Text: struct {
			PreviewURL bool   "json:\"preview_url\""
			Body       string "json:\"body\""
		}{
			PreviewURL: false,
			Body:       content,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("Error while marshalling request body: %w", err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", os.Getenv("PHONE_NUMBER_ID"))
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return 0, fmt.Errorf("Error while constructing request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("META_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error while sending request: %w", err)
	}
	status := resp.StatusCode
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error while reading request body: %w", err)
	}
	fmt.Printf("%s\n", string(responseData))
	return status, nil
}

func SendTemplateMessage(template, to, language string) {}
