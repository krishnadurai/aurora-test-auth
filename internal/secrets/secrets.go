// Package secrets defines a minimum abstract interface for a secret manager.
// Allows for a different implementation to be bound within the servernv.ServeEnv
package secrets

import (
	"context"
	"fmt"
)

// SecretManager defines the minimum shared functionality for a secret manager
// used by this application.
type SecretManager interface {
	GetSecretValue(ctx context.Context, name string) (string, error)
}

// SecretManagerFunc is a func that returns a secret manager or error.
type SecretManagerFunc func(ctx context.Context) (SecretManager, error)

// SecretManagerFor returns the secret manager for the given type, or an error
// if one does not exist.
func SecretManagerFor(ctx context.Context, typ SecretManagerType) (SecretManager, error) {
	switch typ {
	case SecretManagerTypeAWSSecretsManager:
		return NewAWSSecretsManager(ctx)
	case SecretManagerTypeAzureKeyVault:
		return NewAzureKeyVault(ctx)
	case SecretManagerTypeGoogleSecretManager:
		return NewGoogleSecretManager(ctx)
	case SecretManagerTypeGoogleHashiCorpVault:
		return NewHashiCorpVault(ctx)
	case SecretManagerTypeNoop:
		return NewNoop(ctx)
	}

	return nil, fmt.Errorf("unknown secret manager type: %v", typ)
}
