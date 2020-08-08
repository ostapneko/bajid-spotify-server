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
	client    *secretmanager.Client
	projectId string
}

func NewSecretManager(gcpProjectId string) (*SecretManager, error) {
	log.Println("instantiating secrets manager")

	ctx := context.Background()

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secrets manager client: %s", err)
	}

	return &SecretManager{client: client, projectId: gcpProjectId}, nil
}

func (s *SecretManager) GetSecret(ctx context.Context, key string) (string, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.projectId, key)
	versionReq := &secretmanagerpb.GetSecretVersionRequest{
		Name: name,
	}
	_, err := s.client.GetSecretVersion(ctx, versionReq)

	if err != nil {
		return "", fmt.Errorf("error accessing secret version with key %s: %s", key, err)
	}

	name = fmt.Sprintf("projects/%s/secrets/%s/versions/1", s.projectId, key)
	secretReq := &secretmanagerpb.AccessSecretVersionRequest{Name: name}

	secret, err := s.client.AccessSecretVersion(ctx, secretReq)

	if err != nil {
		return "", fmt.Errorf("error accessing secret with key %s: %s", key, err)
	}

	secretValue := string(secret.Payload.Data)
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
