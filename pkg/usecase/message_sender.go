package usecase

import (
	"context"
	"fmt"
	"github.com/yomafleet/better-custom-message-sender/config"
	"github.com/yomafleet/better-custom-message-sender/pkg/domain/contract"
	usecaseerr "github.com/yomafleet/better-custom-message-sender/pkg/usecase/e"
	"github.com/yomafleet/better-custom-message-sender/pkg/usecase/model"
	"strings"
)

type MessageSender struct {
	sender    contract.MessageSenderContract
	encrypter contract.EncrypterContract
	logger    contract.Logger
	params    *config.EnvProvider
}

func NewMessageSender(params *config.EnvProvider, sender contract.MessageSenderContract,
	encrypter contract.EncrypterContract, logger contract.Logger) *MessageSender {

	return &MessageSender{
		sender:    sender,
		encrypter: encrypter,
		logger:    logger,
		params:    params,
	}
}

const (
	SMSSignUp              = "CustomSMSSender_SignUp"
	SMSResendCode          = "CustomSMSSender_ResendCode"
	SMSForgotPassword      = "CustomSMSSender_ForgotPassword"
	SMSUpdateUserAttribute = "CustomSMSSender_UpdateUserAttribute"
	SMSVerifyUserAttribute = "CustomSMSSender_VerifyUserAttribute"
	SMSAdminCreateUser     = "CustomSMSSender_AdminCreateUser"

	DefaultUserName  = "User"
	DefaultSMSSender = "Yoma Fleet"
)

func (r *MessageSender) Handle(ctx context.Context, input model.CognitoEventUserPoolsCustomMessage) error {
	r.logger.Info("event", input)

	// Decrypt the OTP code from Cognito
	output, err := r.encrypter.Decrypt(ctx, contract.DecryptInput{
		EncryptedData: input.Request.CodeParameter,
	})
	if err != nil {
		r.logger.Error("decryption failed", err)
		return err
	}

	r.logger.Info("decrypted data", output)

	// Determine if this is an SMS trigger
	sources := strings.Split(input.TriggerSource, "_")
	r.logger.Info("trigger source", sources)

	if len(sources) > 1 && sources[0] == "CustomSMSSender" {
		return r.sendSMS(ctx, output.DecryptedData, input)
	}

	return fmt.Errorf("unsupported trigger source: %s (only SMS triggers are supported)", input.TriggerSource)
}

func (r *MessageSender) sendSMS(ctx context.Context, code string, input model.CognitoEventUserPoolsCustomMessage) error {
	// Get phone number from user attributes
	phoneNumber, ok := input.Request.UserAttributes["phone_number"].(string)
	if !ok {
		err := usecaseerr.NewApplicationError(config.StatusCode500ExternalIntegrationError, "phone number not found in user attributes")
		r.logger.Error(err.Error(), nil)
		return err
	}

	// Get user name (optional, default to "User")
	userName, ok := input.Request.UserAttributes["name"].(string)
	if !ok {
		userName = DefaultUserName
	}

	// Format the SMS message based on trigger type
	message := r.formatSMSMessage(input.TriggerSource, code, userName)

	// Send SMS via the message sender
	output, err := r.sender.SendMessage(ctx, contract.MessageSendInput{
		MessageType: config.SMSMessageType,
		Topic:       "",
		Message:     message,
		From:        DefaultSMSSender,
		To:          phoneNumber,
		Subject:     "",
		Data:        nil,
	})
	if err != nil {
		r.logger.Error(err.Error(), nil)
		return err
	}

	r.logger.Info("SMS sent successfully", output)

	return nil
}

// formatSMSMessage formats the SMS message based on the trigger type
func (r *MessageSender) formatSMSMessage(triggerSource, code, userName string) string {
	switch triggerSource {
	case SMSSignUp:
		return fmt.Sprintf("Welcome to Yoma Fleet Better! Your verification code is %s. DO NOT share it with anyone.", code)
	case SMSForgotPassword:
		return fmt.Sprintf("Your Yoma Fleet Better password reset code is %s. DO NOT share it with anyone.", code)
	case SMSResendCode:
		return fmt.Sprintf("Your Yoma Fleet Better verification code is %s. DO NOT share it with anyone.", code)
	case SMSAdminCreateUser:
		return fmt.Sprintf("Welcome to Yoma Fleet Better! Your temporary password is %s. Please change it after login.", code)
	case SMSUpdateUserAttribute:
		return fmt.Sprintf("Your Yoma Fleet Better verification code is %s. DO NOT share it with anyone.", code)
	case SMSVerifyUserAttribute:
		return fmt.Sprintf("Your Yoma Fleet Better verification code is %s. DO NOT share it with anyone.", code)
	default:
		return fmt.Sprintf("Your OTP for Yoma Fleet Better is %s. DO NOT share it with anyone", code)
	}
}
