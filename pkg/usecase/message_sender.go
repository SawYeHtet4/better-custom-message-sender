package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yomafleet/better-custom-message-sender/config"
	"github.com/yomafleet/better-custom-message-sender/pkg/domain/contract"
	usecaseerr "github.com/yomafleet/better-custom-message-sender/pkg/usecase/e"
	"github.com/yomafleet/better-custom-message-sender/pkg/usecase/model"
	"html/template"
	"slices"
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
	SingUp                      = "CustomEmailSender_SignUp"
	ForgetPassword              = "CustomEmailSender_ForgotPassword"
	ResendCode                  = "CustomEmailSender_ResendCode"
	UpdateUserAttribute         = "CustomEmailSender_UpdateUserAttribute"
	VerifyUserAttribute         = "CustomEmailSender_VerifyUserAttribute"
	AdminCreateUser             = "CustomEmailSender_AdminCreateUser"
	AccountTakeOverNotification = "CustomEmailSender_AccountTakeOverNotification"

	SMSSingUp              = "CustomSMSSender_SignUp"
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
	output, err := r.encrypter.Decrypt(ctx, contract.DecryptInput{
		EncryptedData: input.Request.CodeParameter,
	})
	if err != nil {
		return err
	}

	r.logger.Info("decrypted data", output)

	// Determine if this is an SMS or Email trigger
	sources := strings.Split(input.TriggerSource, "_")
	r.logger.Info("trigger source", sources)

	if len(sources) > 1 {
		switch sources[0] {
		case "CustomSMSSender":
			return r.sendSMS(ctx, output.DecryptedData, input)
		case "CustomEmailSender":
			return r.sendEmail(ctx, output.DecryptedData, input)
		}
	}

	return fmt.Errorf("unknown trigger source: %s", input.TriggerSource)
}

func (r *MessageSender) sendEmail(ctx context.Context, code string, input model.CognitoEventUserPoolsCustomMessage) error {
	email, ok := input.Request.UserAttributes["email"].(string)
	if !ok {
		err := usecaseerr.NewApplicationError(config.StatusCode500ExternalIntegrationError, "email not found in user attributes")
		r.logger.Error(err.Error(), nil)
		return err
	}
	if slices.Contains(r.params.App().BlackListEmails(), email) {
		err := usecaseerr.NewApplicationError(config.StatusCode403ResourceForbidden, "email is blacklisted")
		r.logger.Error(err.Error(), map[string]string{"email": email})
		return err
	}
	userName, ok := input.Request.UserAttributes["name"].(string)
	if !ok {
		userName = DefaultUserName
	}

	var output contract.MessageSendResponse
	var err error
	switch input.TriggerSource {
	case SingUp:
		output, err = r.sendEmailForUserSingUp(ctx, sendEmailForSingleEventInput{
			code:          code,
			userName:      userName,
			fullName:      "",
			email:         email,
			emailVerified: "",
		})
	case ForgetPassword:
		output, err = r.sendEmailForUserForgotPassword(ctx, sendEmailForSingleEventInput{
			code:          code,
			userName:      userName,
			fullName:      "",
			email:         email,
			emailVerified: "",
		})
	case ResendCode:
		emailVerified, _ := input.Request.UserAttributes["email_verified"].(string)
		output, err = r.sendEmailForUserResendCode(ctx, sendEmailForSingleEventInput{
			code:          code,
			userName:      userName,
			fullName:      "",
			email:         email,
			emailVerified: emailVerified,
		})
	case AdminCreateUser:
		output, err = r.sendEmailForAdminCreateUser(ctx, sendEmailForSingleEventInput{
			code:          code,
			userName:      userName,
			fullName:      input.UserName,
			email:         email,
			emailVerified: "",
		})
	}

	r.logger.Info("email send successfully", output)

	return err
}

func (r *MessageSender) buildSubject(title string) string {

	return fmt.Sprintf("'Yoma Fleet Better' | %s", title)
}

type htmlTemplateGetInput struct {
	Code                string
	Name                string
	UserName            string
	TemplateContentName string
}

func (r *MessageSender) getHTMLTemplate(input htmlTemplateGetInput) (string, error) {
	tmpl, err := template.New("index.html").
		ParseGlob(
			"./assets/templates/email/*.html",
		)
	if err != nil {
		return "", err
	}

	// Execute the template with the data
	var htmlByteContent bytes.Buffer
	err = tmpl.ExecuteTemplate(&htmlByteContent, "base", input)
	if err != nil {
		return "", err
	}

	return string(htmlByteContent.Bytes()), nil
}
func (r *MessageSender) withHTML() bool {

	return r.params.App().MailTransport() != config.SQSTransport
}

type sendEmailForSingleEventInput struct {
	code          string
	userName      string
	fullName      string
	email         string
	emailVerified string
}

type sendEmailForAllEventInput struct {
	sendEmailForSingleEventInput
	topic        string
	subject      string
	templateName string
	data         map[string]any
}

func (r *MessageSender) sendEmailForUserSingUp(ctx context.Context, input sendEmailForSingleEventInput) (contract.MessageSendResponse, error) {
	return r.sendEmailForEvent(ctx, sendEmailForAllEventInput{
		sendEmailForSingleEventInput: input,
		topic:                        "otp_verify",
		subject:                      "Verify your email",
		templateName:                 "verify",
		data: map[string]any{
			"user": map[string]any{"name": input.userName},
			"type": "register",
			"code": input.code,
		},
	})
}

func (r *MessageSender) sendEmailForUserForgotPassword(ctx context.Context, input sendEmailForSingleEventInput) (contract.MessageSendResponse, error) {
	return r.sendEmailForEvent(ctx, sendEmailForAllEventInput{
		sendEmailForSingleEventInput: input,
		topic:                        "otp_verify",
		subject:                      "Password Reset",
		templateName:                 "forgot_password",
		data: map[string]any{
			"user": map[string]any{"name": input.userName},
			"type": "forgot_password",
			"code": input.code,
		},
	})
}

func (r *MessageSender) sendEmailForUserResendCode(ctx context.Context, input sendEmailForSingleEventInput) (contract.MessageSendResponse, error) {
	if input.emailVerified == "false" {
		return r.sendEmailForUserSingUp(ctx, input)
	} else {

		return r.sendEmailForEvent(ctx, sendEmailForAllEventInput{
			sendEmailForSingleEventInput: input,
			topic:                        "otp_verify",
			subject:                      "Password Reset",
			templateName:                 "forgot_password",
			data: map[string]any{
				"user": map[string]any{"name": input.userName},
				"type": "resend",
				"code": input.code,
			},
		})
	}
}

func (r *MessageSender) sendEmailForAdminCreateUser(ctx context.Context, input sendEmailForSingleEventInput) (contract.MessageSendResponse, error) {
	return r.sendEmailForEvent(ctx, sendEmailForAllEventInput{
		sendEmailForSingleEventInput: input,
		topic:                        "better_auth_admin_create_user",
		subject:                      "Your Temporary Password",
		templateName:                 "forgot_password",
		data: map[string]any{
			"user": map[string]any{
				"name":           input.userName,
				"username":       input.fullName,
				"attribute_name": "Email",
			},
			"temporary_password": input.code,
		},
	})
}

func (r *MessageSender) sendEmailForEvent(
	ctx context.Context,
	input sendEmailForAllEventInput,
) (contract.MessageSendResponse, error) {
	from := ""
	message := ""
	var err error
	if r.withHTML() {
		input.topic = ""
		from = r.params.App().EmailFrom()
		input.data = nil

		message, err = r.getHTMLTemplate(htmlTemplateGetInput{
			Code:                input.code,
			Name:                input.userName,
			UserName:            input.userName,
			TemplateContentName: input.templateName,
		})
		if err != nil {
			return contract.MessageSendResponse{}, err
		}
	}
	return r.sender.SendMessage(ctx, contract.MessageSendInput{
		MessageType: config.EmailMessageType,
		Topic:       input.topic,
		Message:     message,
		From:        from,
		To:          input.email,
		Subject:     r.buildSubject(input.subject),
		Data:        input.data,
	})
}

func (r *MessageSender) sendSMS(ctx context.Context, code string, input model.CognitoEventUserPoolsCustomMessage) error {
	// Get phone number from user attributes
	phoneNumber, ok := input.Request.UserAttributes["phone_number"].(string)
	if !ok {
		err := usecaseerr.NewApplicationError(config.StatusCode500ExternalIntegrationError, "phone number not found in user attributes")
		r.logger.Error(err.Error(), nil)
		return err
	}

	// Format the SMS message
	message := fmt.Sprintf("Your OTP for Yoma Fleet Better is %s. DO NOT share it with anyone", code)

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
