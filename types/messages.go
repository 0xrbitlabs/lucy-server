package types

import (
	"net/http"
)

type InboundMessage struct {
	SMSMessageSID    string
	NumMedia         string
	ProfileName      string
	MessageType      string
	SMSSID           string
	WaID             string
	SMSStatus        string
	Body             string
	To               string
	NumSegments      string
	ReferralNumMedia string
	MessageSid       string
	AccountSid       string
	From             string
	ApiVersion       string
}

func NewInboundMessage(r *http.Request) *InboundMessage {
	return &InboundMessage{
		SMSMessageSID:    r.FormValue("SmsMessageSid"),
		NumMedia:         r.FormValue("NumMedia"),
		ProfileName:      r.FormValue("ProfileName"),
		MessageType:      r.FormValue("MessageType"),
		SMSSID:           r.FormValue("SmsSid"),
		WaID:             r.FormValue("WaId"),
		SMSStatus:        r.FormValue("SmsStatus"),
		Body:             r.FormValue("Body"),
		To:               r.FormValue("To"),
		NumSegments:      r.FormValue("NumSegments"),
		ReferralNumMedia: r.FormValue("ReferralNumMedia"),
		MessageSid:       r.FormValue("MessageSid"),
		AccountSid:       r.FormValue("AccountSid"),
		From:             r.FormValue("From"),
		ApiVersion:       r.FormValue("ApiVersion"),
	}
}
