package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func SendMessageSingle(to, content string) error {
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
		return fmt.Errorf("Error while marshalling request body: %w", err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", os.Getenv("PHONE_NUMBER_ID"))
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return fmt.Errorf("Error while constructing request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("META_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending request: %w", err)
	}
	status := resp.StatusCode
	if status != 200 {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while reading request body: %w", err)
		}
		return fmt.Errorf("Error while sending message, wanted HTTP 200 but got HTTP %d with error message: %s", status, string(responseData))
	}
	return nil
}

func SendTemplateMessage(template, to, language string) error {
	reqBody, err := json.Marshal(struct {
		MessagingProduct string `json:"messaging_product"`
		RecipientType    string `json:"recipient_type"`
		To               string `json:"to"`
		Type             string `json:"type"`
		Template         struct {
			Name     string `json:"name"`
			Language struct {
				Code string `json:"code"`
			} `json:"language"`
		} `json:"template"`
	}{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "template",
		Template: struct {
			Name     string "json:\"name\""
			Language struct {
				Code string "json:\"code\""
			} "json:\"language\""
		}{
			Name: template,
			Language: struct {
				Code string "json:\"code\""
			}{
				Code: "fr",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Error while marshalling request body: %w", err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", os.Getenv("PHONE_NUMBER_ID"))
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return fmt.Errorf("Error while constructing request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("META_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending request: %w", err)
	}
	status := resp.StatusCode
	if status != 200 {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while reading request body: %w", err)
		}
		return fmt.Errorf("Error while sending template message, wanted HTTP 200 but got HTTP %d with error message: %s", status, string(responseData))
	}
	return nil
}

func SendErrorMessage(to string, logger *slog.Logger) error {
	content := `Une erreur est survenue ;) J'en ai notifie mes createurs et ils travaillent deja dessus. Veuillez reessayer ou les contacter pour obtenir de l'aide`
	err := SendMessageSingle(to, content)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

func SendRegistrationConfirmationMessage(to, targetURL string) error {
	reqBody, err := json.Marshal(struct {
		MessagingProduct string `json:"messaging_product"`
		RecipientType    string `json:"recipient_type"`
		To               string `json:"to"`
		Type             string `json:"type"`
		Interactive      struct {
			Type string `json:"type"`
			Body struct {
				Text string `json:"text"`
			} `json:"body"`
			Action struct {
				Name       string `json:"name"`
				Parameters struct {
					DisplayText string `json:"display_text"`
					URL         string `json:"url"`
				} `json:"parameters"`
			} `json:"action"`
		} `json:"interactive"`
	}{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "interactive",
		Interactive: struct {
			Type string "json:\"type\""
			Body struct {
				Text string "json:\"text\""
			} "json:\"body\""
			Action struct {
				Name       string "json:\"name\""
				Parameters struct {
					DisplayText string "json:\"display_text\""
					URL         string "json:\"url\""
				} "json:\"parameters\""
			} "json:\"action\""
		}{
			Type: "cta_url",
			Body: struct {
				Text string "json:\"text\""
			}{
				Text: "Votre demande d'enregistrement en tant que vendeur a ete approuvee. Veuillez vous connecter sur notre plateforme pour terminer la creation de votre compte",
			},
			Action: struct {
				Name       string "json:\"name\""
				Parameters struct {
					DisplayText string "json:\"display_text\""
					URL         string "json:\"url\""
				} "json:\"parameters\""
			}{
				Name: "cta_url",
				Parameters: struct {
					DisplayText string "json:\"display_text\""
					URL         string "json:\"url\""
				}{
					DisplayText: "Confirmer mon compte",
					URL:         targetURL,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Error while marshalling request body: %w", err)
	}
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", os.Getenv("PHONE_NUMBER_ID"))
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return fmt.Errorf("Error while constructing request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("META_TOKEN")))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while sending request: %w", err)
	}
	status := resp.StatusCode
	if status != 200 {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while reading request body: %w", err)
		}
		return fmt.Errorf("Error while sending message, wanted HTTP 200 but got HTTP %d with error message: %s", status, string(responseData))
	}
	return nil
}
