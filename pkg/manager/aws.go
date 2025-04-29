package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerAPI interface {
	CreateSecret(ctx context.Context, input *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	DeleteSecret(ctx context.Context, input *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)
	UpdateSecret(ctx context.Context, input *secretsmanager.UpdateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error)
}

type AWSSecretManagerClient[T SecretValue] struct {
	client SecretsManagerAPI
}

func NewAWSSecretManagerClient[T SecretValue](client SecretsManagerAPI) *AWSSecretManagerClient[T] {
	return &AWSSecretManagerClient[T]{
		client: client,
	}
}

func (c AWSSecretManagerClient[T]) Create(ctx context.Context, name string, data T) error {
	binary, err := data.GetData()
	if err != nil {
		return err
	}
	_, err = c.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		SecretBinary: binary,
	})
	return err
}

func (c AWSSecretManagerClient[T]) Get(ctx context.Context, name string) (T, error) {
	var zero T
	secret, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return zero, err
	}
	if secret.SecretBinary == nil {
		return zero, fmt.Errorf("secret binary is nil")
	}
	fmt.Printf("get secret: %+v\n", string(secret.SecretBinary))
	return FromBinary[T](secret.SecretBinary)
}

func FromBinary[T SecretValue](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func (c AWSSecretManagerClient[T]) Delete(ctx context.Context, name string) error {
	_, err := c.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(name),
	})
	return err
}

func (c AWSSecretManagerClient[T]) Update(ctx context.Context, name string, data T) error {
	binary, err := data.GetData()
	if err != nil {
		return err
	}
	_, err = c.client.UpdateSecret(ctx, &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(name),
		SecretBinary: binary,
	})
	return err
}
