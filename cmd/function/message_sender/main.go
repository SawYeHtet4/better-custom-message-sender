package main

import (
	"github.com/yomafleet/better-custom-message-sender/pkg/interface/encrypter"
	"github.com/yomafleet/better-custom-message-sender/pkg/interface/logger"
	"github.com/yomafleet/better-custom-message-sender/pkg/interface/message_sender"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yomafleet/better-custom-message-sender/config"
	"github.com/yomafleet/better-custom-message-sender/pkg/usecase"
)

func main() {
	params := config.NewEnvProvider()
	messageSender := message_sender.NewMessageSender(params)
	encrypterAgent := encrypter.NewEncrypter(params.Aws().KeyID())
	loggerAgent := logger.NewLogger()

	handler := usecase.NewMessageSender(params, messageSender, encrypterAgent, loggerAgent)
	lambda.Start(handler.Handle)
}
