package types

type MessageContext int

const (
	ContextFirstMessage MessageContext = iota
	ContextRegistrationRequest
	ContextQuery
)
