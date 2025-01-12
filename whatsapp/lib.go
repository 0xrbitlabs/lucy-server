package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const whatsappAPIURL = "https://graph.facebook.com/v21.0"

type Client struct {
	AccessToken   string
	PhoneNumberID string
}

func NewClient(accessToken, phoneNumberID string) *Client {
	return &Client{AccessToken: accessToken, PhoneNumberID: phoneNumberID}
}

func (c *Client) SendBasicMessage(recipient, message string) error {
	payload := BasicMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               recipient,
		Type:             "text",
		Text: TextPayload{
			PreviewURL: false,
			Body:       message,
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error while sending basic message: Failed to marshal payload %w", err)
	}
	url := fmt.Sprintf("%s/%s/messages", whatsappAPIURL, c.PhoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error while sending basic message: Failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending basic message: Failed to send HTTP request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error while sending basic message: Failed to read response body: %w", err)
	}
	fmt.Println(string(responseBody))
	fmt.Println(response.Status)
	response.Body.Close()
	return nil
}

func (c *Client) SendVerificationCodeMessage(to, code string) error {
	payload := TemplateMessagePayload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template: Template{
			Name: "verification_code",
			Language: TemplateLanguage{
				Code: "fr",
			},
			Components: []TemplateComponent{
        {
          Type: "button",
          SubType: "url",
          Index: "0",
          Parameters: []ComponentParameter{
            {
              Type: "text",
              Text: code,
            },
          },
        },
				{
					Type: "body",
					Parameters: []ComponentParameter{
						{
							Type: "text",
							Text: code,
						},
					},
				},
			},
		},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error while sending template message: Failed to marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/%s/messages", whatsappAPIURL, c.PhoneNumberID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error while sending template message: Failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending template message: Failed to send HTTP request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error while sending template message: Failed to read response body: %w", err)
	}
	fmt.Println(string(responseBody))
	fmt.Println(response.Status)
	response.Body.Close()
	return nil
}
