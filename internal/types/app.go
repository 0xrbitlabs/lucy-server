package types

type ButtonSchema struct {
	Text    string `json:"text"`
	Payload string `json:"payload"`
}

type TextSchema struct {
	Body string `json:"body"`
}

type ContextSchema struct {
	From string `json:"from"`
	ID   string `json:"id"`
}

type MessageSchema struct {
	Context   ContextSchema `json:"context"`
	From      string        `json:"from"`
	ID        string        `json:"id"`
	TimeStamp string        `json:"timestamp"`
	Type      string        `json:"type"`
	Text      TextSchema    `json:"text"`
	Button    ButtonSchema  `json:"button"`
}

type ProfileSchema struct {
	Name string `json:"name"`
}

type ContactSchema struct {
	Profile ProfileSchema `json:"profile"`
	WaID    string        `json:"wa_id"`
}

type MetaDataSchema struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type ValueSchema struct {
	MessagingProduct string          `json:"messaging_product"`
	MetaData         MetaDataSchema  `json:"metadata"`
	Contacts         []ContactSchema `json:"contacts"`
	Messages         []MessageSchema `json:"messages"`
}

type ChangeSchema struct {
	Field string      `json:"field"`
	Value ValueSchema `json:"value"`
}

type EntrySchema struct {
	ID      string         `json:"id"`
	Changes []ChangeSchema `json:"changes"`
}

type WebhookMessage struct {
	Object string        `json:"object"`
	Entry  []EntrySchema `json:"entry"`
}
