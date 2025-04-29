package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/takutakahashi/secretstore-go/pkg/manager"
	"github.com/takutakahashi/secretstore-go/pkg/secretvalue"
)

func main() {
	// AWS Secrets Manager クライアントの作成
	// 環境変数から認証情報と設定を読み込む
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("AWS設定の読み込みに失敗しました: %v", err)
	}
	smClient := secretsmanager.NewFromConfig(cfg)
	client := manager.NewAWSSecretManagerClient[*secretvalue.EnvSecret](smClient)

	// 環境変数の設定
	envVars := map[string]string{
		"DATABASE_URL":     "postgres://user:pass@localhost:5432/db",
		"API_KEY":          "secret-api-key",
		"REDIS_CONNECTION": "redis://localhost:6379",
	}
	secret := secretvalue.NewEnvSecret(envVars)
	fmt.Printf("secret: %+v\n", secret)

	ctx := context.Background()
	secretName := fmt.Sprintf("my-app-env-%d", time.Now().UnixNano())

	// シークレットの作成
	if err := client.Create(ctx, secretName, secret); err != nil {
		log.Fatalf("Failed to create secret: %v", err)
	}
	fmt.Printf("Created secret: %s\n", secretName)

	// シークレットの取得
	retrievedSecret, err := client.Get(ctx, secretName)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}
	fmt.Printf("Retrieved secret: %+v\n", retrievedSecret)

	// シークレットの更新
	envVars["NEW_KEY"] = "new-value"
	updatedSecret := secretvalue.NewEnvSecret(envVars)
	fmt.Printf("updatedSecret: %+v\n", updatedSecret)
	if err := client.Update(ctx, secretName, updatedSecret); err != nil {
		log.Fatalf("Failed to update secret: %v", err)
	}
	fmt.Printf("Updated secret: %s\n", secretName)

	// シークレットの削除
	if err := client.Delete(ctx, secretName); err != nil {
		log.Fatalf("Failed to delete secret: %v", err)
	}
	fmt.Printf("Deleted secret: %s\n", secretName)
}
