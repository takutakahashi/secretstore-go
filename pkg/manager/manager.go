package manager

import "context"

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
