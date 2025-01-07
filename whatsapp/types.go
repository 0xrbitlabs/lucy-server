package whatsapp

const (
	TextMessage MessageType = "text"
)

type MessageType string

type MetaData struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Text struct {
	Body string `json:"body"`
}

type Message struct {
	From      string      `json:"from"`
	ID        string      `json:"id"`
	Timestamp string      `json:"timestamp"`
	Text      Text        `json:"text"`
	Type      MessageType `json:"type"`
}

type Profile struct {
	Name string `json:"name"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	MetaData         MetaData  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Change struct {
	Value Value  `json:"value"`
	Field string `json:"field"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Payload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Envelope struct {
	Object string `json:"object"`
	Entry  Entry  `json:"entry"`
}
