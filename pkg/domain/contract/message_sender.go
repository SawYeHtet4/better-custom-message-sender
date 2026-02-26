package contract

import "context"

type MessageSenderContract interface {
	SendMessage(ctx context.Context, input MessageSendInput) (MessageSendResponse, error)
}

type MessageSendInput struct {
	MessageType string
	Topic       string
	Message     string
	From        string
	To          string
	Subject     string
	Data        map[string]any
}

type MessageSendResponse struct {
	MessageID string
}
