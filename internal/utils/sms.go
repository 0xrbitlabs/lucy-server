package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func SendMessage(receiver, message string) error {
	apiURL := "https://rest.nexmo.com/sms/json"
	formData := url.Values{}
	formData.Set("from", "Ask Lucy")
	formData.Set("text", message)
	formData.Set("to", receiver)
	formData.Set("api_key", os.Getenv("VONAGE_API_KEY"))
	formData.Set("api_secret", os.Getenv("VONAGE_API_SECRET"))
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("Error while creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending HTTP request: %w", err)
	}
  fmt.Println(resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error while sending SMS: Wanted HTTP 200 but got %s", resp.Status)
	}
	return nil
}
