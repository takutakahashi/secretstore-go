package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockSecretValue struct {
	Value string `json:"value"`
	Fail  bool   `json:"fail"`
}

func (m mockSecretValue) Data() ([]byte, error) {
	if m.Fail {
		return nil, errors.New("fail")
	}
	return []byte(m.Value), nil
}

type mockSecretsManagerClient struct {
	CreateSecretFunc   func(ctx context.Context, input *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	GetSecretValueFunc func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	DeleteSecretFunc   func(ctx context.Context, input *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)
	UpdateSecretFunc   func(ctx context.Context, input *secretsmanager.UpdateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error)
}

func (m *mockSecretsManagerClient) CreateSecret(ctx context.Context, input *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
	return m.CreateSecretFunc(ctx, input, optFns...)
}
func (m *mockSecretsManagerClient) GetSecretValue(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFunc(ctx, input, optFns...)
}
func (m *mockSecretsManagerClient) DeleteSecret(ctx context.Context, input *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error) {
	if m.DeleteSecretFunc != nil {
		return m.DeleteSecretFunc(ctx, input, optFns...)
	}
	return &secretsmanager.DeleteSecretOutput{}, nil
}
func (m *mockSecretsManagerClient) UpdateSecret(ctx context.Context, input *secretsmanager.UpdateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error) {
	if m.UpdateSecretFunc != nil {
		return m.UpdateSecretFunc(ctx, input, optFns...)
	}
	return &secretsmanager.UpdateSecretOutput{}, nil
}

func TestAWSSecretManagerClient_Create(t *testing.T) {
	want := `{"foo":"bar"}`
	client := &AWSSecretManagerClient[mockSecretValue]{
		client: &mockSecretsManagerClient{
			CreateSecretFunc: func(ctx context.Context, input *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error) {
				if !bytes.Equal(input.SecretBinary, []byte(want)) {
					t.Errorf("unexpected SecretBinary: got %v, want %v", input.SecretBinary, []byte(want))
				}
				return &secretsmanager.CreateSecretOutput{}, nil
			},
		},
	}
	val := mockSecretValue{Value: want}
	if err := client.Create(context.Background(), "test", val); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAWSSecretManagerClient_Get(t *testing.T) {
	want := mockSecretValue{Value: `{"foo":"bar"}`}
	b, _ := json.Marshal(want)
	client := &AWSSecretManagerClient[mockSecretValue]{
		client: &mockSecretsManagerClient{
			GetSecretValueFunc: func(ctx context.Context, input *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
				return &secretsmanager.GetSecretValueOutput{SecretBinary: b}, nil
			},
		},
	}
	got, err := client.Get(context.Background(), "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Value != want.Value || got.Fail != want.Fail {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestFromBinary(t *testing.T) {
	want := mockSecretValue{Value: "abc"}
	b, _ := json.Marshal(want)
	got, err := FromBinary[mockSecretValue](b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Value != want.Value || got.Fail != want.Fail {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestMarshalUnmarshalMockSecretValue(t *testing.T) {
	orig := mockSecretValue{Value: "abc", Fail: false}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var v2 mockSecretValue
	err = json.Unmarshal(b, &v2)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if v2.Value != orig.Value || v2.Fail != orig.Fail {
		t.Errorf("marshal/unmarshal mismatch: got %+v, want %+v", v2, orig)
	}
}
