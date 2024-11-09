package domain

import whatsapptypes "github.com/joseph0x45/lucy/whatsapp_types"

type Context struct {
	// is true if the user is messaging the bot for the first time
	FirstMessage bool
	// contains the payload sent by Meta
	Envelope *whatsapptypes.Envelope
}

func NewContext(envelope *whatsapptypes.Envelope) *Context {
	return &Context{
		Envelope: envelope,
	}
}
