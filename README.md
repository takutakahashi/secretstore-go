# secretstore-go

A Go library for securely storing and managing sensitive information in various secret management services. The library's unique feature is its ability to store any Go struct that implements the SecretValue interface, providing type-safe secret management.

## Features

- **Type-safe Secret Management**: Store any Go struct that implements the SecretValue interface
- **Generic Implementation**: Utilize Go's generics for type-safe secret operations
- **Multiple Backend Support**: Designed to support various secret management services (currently supports AWS Secrets Manager)
- **Simple Interface**: Consistent API across different backend implementations
- **Secure by Design**: Built with security best practices in mind

## Installation

```bash
go get github.com/takutakahashi/secretstore-go
```

## Usage

### AWS Secrets Manager Example

```go
package main

import (
    "context"
    "log"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/takutakahashi/secretstore-go/pkg/manager"
    "github.com/takutakahashi/secretstore-go/pkg/secretvalue"
)

func main() {
    // Load AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatalf("Failed to load AWS config: %v", err)
    }

    // Create AWS Secrets Manager client
    smClient := secretsmanager.NewFromConfig(cfg)
    client := manager.NewAWSSecretManagerClient[*secretvalue.EnvSecret](smClient)

    // Create environment variables to store
    envVars := map[string]string{
        "DATABASE_URL": "postgres://user:pass@localhost:5432/db",
        "API_KEY":     "secret-api-key",
    }
    secret := secretvalue.NewEnvSecret(envVars)

    // Store the secret
    err = client.Create(context.Background(), "my-app-env", secret)
    if err != nil {
        log.Fatalf("Failed to create secret: %v", err)
    }
}
```

### Creating Custom Secret Types

You can create your own secret types by implementing the SecretValue interface:

```go
type SecretValue interface {
    GetData() ([]byte, error)
    SetData([]byte) error
}
```

Example implementation:

```go
type CustomSecret struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (s CustomSecret) GetData() ([]byte, error) {
    return json.Marshal(s)
}

func (s *CustomSecret) SetData(data []byte) error {
    return json.Unmarshal(data, s)
}
```

## Supported Secret Managers

- AWS Secrets Manager
- More coming soon...

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details 