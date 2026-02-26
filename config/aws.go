package config

import (
	"os"
)

type AwsConfig struct {
	keyID       string
	sqsQueueURL string
}

func (a AwsConfig) SqsQueueURL() string {
	return a.sqsQueueURL
}

func (a AwsConfig) KeyID() string {
	return a.keyID
}

func newAWSConfig() *AwsConfig {
	return &AwsConfig{
		keyID:       os.Getenv("KEY_ID"),
		sqsQueueURL: os.Getenv("SQS_QUEUE_URL"),
	}
}
