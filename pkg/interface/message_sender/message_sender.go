package message_sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	gomail "github.com/wneessen/go-mail"
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
	switch input.MessageType {
	case config.SMSMessageType:
		return r.sendSMS(ctx, input)
	case config.EmailMessageType:
		switch r.params.App().MailTransport() {
		case config.SESTransport:
			return r.sendEmailWithSES(ctx, input)
		case config.SQSTransport:
			return r.sendEmailWithSQS(ctx, input)
		default:
			return r.sendEmailWithMailTrap(ctx, input)
		}
	default:
		return contract.MessageSendResponse{}, fmt.Errorf("unsupported message type: %s", input.MessageType)
	}
}

func (r *MessageSender) sendEmailWithSES(ctx context.Context, input contract.MessageSendInput) (contract.MessageSendResponse, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return contract.MessageSendResponse{}, err
	}
	client := ses.NewFromConfig(cfg)
	output, err := client.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			BccAddresses: nil,
			CcAddresses:  nil,
			ToAddresses:  []string{input.To},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(input.Message),
					Charset: aws.String("UTF-8"),
				},
				Text: nil,
			},
			Subject: &types.Content{
				Data:    aws.String(input.Subject),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(input.From),
	})
	if err != nil {
		return contract.MessageSendResponse{}, err
	}
	messageId := ""
	if output.MessageId != nil {
		messageId = *output.MessageId
	}
	return contract.MessageSendResponse{
		MessageID: messageId,
	}, nil

}

type sendEmailWithSQSInput struct {
	Topic string         `json:"topic"`
	Email emailInput     `json:"email"`
	Data  map[string]any `json:"data"`
}

type emailInput struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
}

func (r *MessageSender) sendEmailWithSQS(ctx context.Context, input contract.MessageSendInput) (contract.MessageSendResponse, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return contract.MessageSendResponse{}, err
	}
	client := sqs.NewFromConfig(cfg)

	inputStr, err := json.Marshal(&sendEmailWithSQSInput{
		Topic: input.Topic,
		Email: emailInput{
			To:      input.To,
			Subject: input.Subject,
		},
		Data: input.Data,
	})
	if err != nil {
		return contract.MessageSendResponse{}, err
	}

	output, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(inputStr)),
		QueueUrl:    aws.String(r.params.Aws().SqsQueueURL()),
	})
	if err != nil {
		return contract.MessageSendResponse{}, err
	}

	messageId := ""
	if output.MessageId != nil {
		messageId = *output.MessageId
	}
	return contract.MessageSendResponse{
		MessageID: messageId,
	}, nil

}

func (r *MessageSender) sendEmailWithMailTrap(ctx context.Context, input contract.MessageSendInput) (contract.MessageSendResponse, error) {
	message := gomail.NewMsg()

	message.SetGenHeader("From", input.From)
	message.SetGenHeader("To", input.To)
	message.SetGenHeader("Subject", input.Subject)
	message.SetBodyString("text/html", input.Message)

	port := r.params.App().MailTrapPort()
	if port == 0 {
		port = gomail.DefaultPort
	}

	client, err := gomail.NewClient(r.params.App().MailTrapHost(),
		gomail.WithSMTPAuth(gomail.SMTPAuthLogin),
		gomail.WithUsername(r.params.App().MailTrapUser()),
		gomail.WithPassword(r.params.App().MailTrapPassword()),
		gomail.WithPort(port),
	)
	if err != nil {
		return contract.MessageSendResponse{}, err
	}
	err = client.Send(message)
	if err != nil {
		return contract.MessageSendResponse{}, err
	}

	return contract.MessageSendResponse{}, nil
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
