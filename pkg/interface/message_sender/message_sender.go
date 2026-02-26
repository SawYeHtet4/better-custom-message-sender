package message_sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/yomafleet/better-custom-message-sender/config"
	"github.com/yomafleet/better-custom-message-sender/pkg/domain/contract"
)

type MessageSender struct {
	params *config.EnvProvider
}

func NewMessageSender(params *config.EnvProvider) contract.MessageSenderContract {

	return &MessageSender{
		params: params,
	}
}

func (r *MessageSender) SendMessage(ctx context.Context, input contract.MessageSendInput) (contract.MessageSendResponse, error) {
	if input.MessageType == config.SMSMessageType {
		return r.sendSMS(ctx, input)
	}
	return contract.MessageSendResponse{}, fmt.Errorf("unsupported message type: %s", input.MessageType)
}

type smsAPIRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

type smsAPIResponse struct {
	Status bool `json:"status"`
}

func (r *MessageSender) sendSMS(ctx context.Context, input contract.MessageSendInput) (contract.MessageSendResponse, error) {
	// Prepare SMS API request payload
	smsPayload := smsAPIRequest{
		To:      input.To,
		Message: input.Message,
		Sender:  input.From,
	}

	jsonData, err := json.Marshal(smsPayload)
	if err != nil {
		return contract.MessageSendResponse{}, fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	// Create HTTP request to SMS API
	smsHost := r.params.App().SmsHost()
	if smsHost == "" {
		return contract.MessageSendResponse{}, fmt.Errorf("SMS_HOST not configured")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", smsHost+"/api/v2/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return contract.MessageSendResponse{}, fmt.Errorf("failed to create SMS request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	authKey := r.params.App().SmsAuthKey()
	if authKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authKey))
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return contract.MessageSendResponse{}, fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return contract.MessageSendResponse{}, fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return contract.MessageSendResponse{}, fmt.Errorf("SMS API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var smsResp smsAPIResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		return contract.MessageSendResponse{}, fmt.Errorf("failed to parse SMS response: %w", err)
	}

	if !smsResp.Status {
		return contract.MessageSendResponse{}, fmt.Errorf("SMS sending failed: status is false")
	}

	return contract.MessageSendResponse{
		MessageID: "sms-sent", // SMS API might not return a message ID, use a placeholder
	}, nil
}
