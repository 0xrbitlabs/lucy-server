package models

type UserProfile struct {
	Name string `json:"name"`
}

type InboundMessage struct {
	To            string      `json:"to"`
	From          string      `json:"from"`
	Channel       string      `json:"channel"`
	MessageUUID   string      `json:"message_uuid"`
	TimeStamp     string      `json:"time_stamp"`
	MessageType   string      `json:"message_type"`
	Text          string      `json:"text"`
	ContextStatus string      `json:"context_status"`
	Profile       UserProfile `json:"profile"`
}
