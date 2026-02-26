package logger

import (
	"github.com/yomafleet/better-custom-message-sender/pkg/domain/contract"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger() contract.Logger {

	return &Logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (r *Logger) Info(message string, data any) {
	if data == nil {
		r.logger.Info(message)

	} else {
		r.logger.Info(message, "data", data)
	}

}

func (r *Logger) Error(message string, data any) {
	if data == nil {
		r.logger.Error(message)

	} else {
		r.logger.Error(message, "data", data)
	}
}
