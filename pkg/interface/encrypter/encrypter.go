package encrypter

import (
	"context"
	"encoding/base64"
	mpl "github.com/aws/aws-cryptographic-material-providers-library/releases/go/mpl/awscryptographymaterialproviderssmithygenerated"
	mpltypes "github.com/aws/aws-cryptographic-material-providers-library/releases/go/mpl/awscryptographymaterialproviderssmithygeneratedtypes"
	client "github.com/aws/aws-encryption-sdk/releases/go/encryption-sdk/awscryptographyencryptionsdksmithygenerated"
	esdktypes "github.com/aws/aws-encryption-sdk/releases/go/encryption-sdk/awscryptographyencryptionsdksmithygeneratedtypes"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/yomafleet/better-custom-message-sender/pkg/domain/contract"
)

type Encrypter struct {
	keyID string
}

func NewEncrypter(keyID string) contract.EncrypterContract {

	return &Encrypter{keyID: keyID}
}

func (r *Encrypter) Decrypt(ctx context.Context, input contract.DecryptInput) (contract.DecryptOutput, error) {
	decodeValue, err := base64.StdEncoding.DecodeString(input.EncryptedData)
	if err != nil {
		return contract.DecryptOutput{}, err
	}

	return r.decryptWithAwsKMS(ctx, decodeValue)
}

func (r *Encrypter) decryptWithAwsKMS(ctx context.Context, decodeValue []byte) (contract.DecryptOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return contract.DecryptOutput{}, err
	}
	kmsClient := kms.NewFromConfig(cfg)

	matProv, err := mpl.NewClient(mpltypes.MaterialProvidersConfig{})
	if err != nil {
		return contract.DecryptOutput{}, err
	}
	awsKmsKeyring, err := matProv.CreateAwsKmsKeyring(ctx, mpltypes.CreateAwsKmsKeyringInput{
		KmsClient:   kmsClient,
		KmsKeyId:    r.keyID,
		GrantTokens: nil,
	})
	if err != nil {
		return contract.DecryptOutput{}, err
	}
	commitmentPolicy := mpltypes.ESDKCommitmentPolicyRequireEncryptAllowDecrypt
	encryptionClient, err := client.NewClient(esdktypes.AwsEncryptionSdkConfig{
		CommitmentPolicy:      &commitmentPolicy,
		MaxEncryptedDataKeys:  nil,
		NetV4_0_0_RetryPolicy: nil,
	})
	if err != nil {
		return contract.DecryptOutput{}, err
	}

	output, err := encryptionClient.Decrypt(ctx, esdktypes.DecryptInput{
		EncryptionContext: nil,
		Keyring:           awsKmsKeyring,
		Ciphertext:        decodeValue,
	})
	if err != nil {
		return contract.DecryptOutput{}, err
	}

	return contract.DecryptOutput{
		DecryptedData: string(output.Plaintext),
	}, err
}
