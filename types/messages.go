package types

type Text struct {
	Body string
}

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	TimeStamp string `json:"timestamp"`
	Text      Text   `json:"text"`
	Type      string `json:"type"`
}

type Profile struct {
	Name string `json:"name"`
}

type Contact struct {
	Profile Profile `json:"profil"`
	WaID    string  `json:"wa_id"`
}

type MetaData struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	MetaData         MetaData  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type InboundMessage struct {
	Object string `json:"object"`
	Entry  []Entry
}
