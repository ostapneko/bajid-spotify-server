package gcp

import (
	"context"
	"fmt"
	"log"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	client *secretmanager.Client
}

func NewSecretManager(gcpProjectId string) (*SecretManager, error) {
	log.Println("instantiating secrets manager")

	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secrets manager client: %s", err)
	}

	return &SecretManager{client: client}, nil
}

func (s *SecretManager) GetSecret(ctx context.Context, key string) (string, error) {
	req := &secretmanagerpb.GetSecretRequest{Name: key}
	secret, err := s.client.GetSecret(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error accessing secret with key %s: %s", key, err)
	}
	secretValue := secret.String()
	if secretValue == "" {
		return "", fmt.Errorf("empty value for secret with key %s", key)
	}

	obfuscated := strings.Builder{}
	// abc123 -> a****3
	for i, _ := range secretValue {
		if i == 0 || i == len(secretValue)-1 {
			obfuscated.WriteByte(secretValue[i])
		} else {
			obfuscated.WriteByte('*')
		}
	}
	log.Printf("secret %s: %s\n", key, obfuscated.String())

	return secretValue, nil
}
