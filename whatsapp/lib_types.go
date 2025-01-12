package whatsapp

type TextPayload struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type BasicMessagePayload struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             TextPayload `json:"text"`
}

type TemplateLanguage struct {
	Code string `json:"code"`
}

type ComponentParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type TemplateComponent struct {
	Type       string               `json:"type"`
	SubType    string               `json:"sub_type,omitempty"`
	Index      string               `json:"index,omitempty"`
	Parameters []ComponentParameter `json:"parameters"`
}

type Template struct {
	Name       string              `json:"name"`
	Language   TemplateLanguage    `json:"language"`
	Components []TemplateComponent `json:"components"`
}

type TemplateMessagePayload struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}
