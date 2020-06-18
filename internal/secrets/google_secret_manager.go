package secrets

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Compile-time check to verify implements interface.
var _ SecretManager = (*GoogleSecretManager)(nil)

// GoogleSecretManager implements SecretManager.
type GoogleSecretManager struct {
	client *secretmanager.Client
}

// NewGoogleSecretManager creates a new secret manager for GCP.
func NewGoogleSecretManager(ctx context.Context) (SecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretmanager.NewClient: %w", err)
	}

	sm := &GoogleSecretManager{
		client: client,
	}

	return sm, nil
}

// GetSecretValue implements the SecretManager interface. Secret names should be
// of the format:
//
//     projects/my-project/secrets/my-secret/versions/123
func (sm *GoogleSecretManager) GetSecretValue(ctx context.Context, name string) (string, error) {
	result, err := sm.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access secret %v: %w", name, err)
	}
	return string(result.Payload.Data), nil
}
