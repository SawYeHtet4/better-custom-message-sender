package contract

import "context"

type EncrypterContract interface {
	Decrypt(ctx context.Context, input DecryptInput) (DecryptOutput, error)
}

type DecryptInput struct {
	EncryptedData string
}

type DecryptOutput struct {
	DecryptedData string
}
