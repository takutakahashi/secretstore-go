package manager

import (
	"context"
	"encoding/json"
)

type SecretValue interface {
	GetData() ([]byte, error)
	SetData([]byte) error
}

type Client[T SecretValue] interface {
	Get(ctx context.Context, name string) (T, error)
	Create(ctx context.Context, name string, data T) error
	Update(ctx context.Context, name string, data T) error
	Delete(ctx context.Context, name string) error
}

func FromBinary[T SecretValue](data []byte) (T, error) {
	var result T
	var zero T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return zero, err
	}
	if err := result.SetData(data); err != nil {
		return zero, err
	}
	return result, nil
}
