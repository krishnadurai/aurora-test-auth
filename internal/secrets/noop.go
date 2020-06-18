package secrets

import (
	"context"
)

// Compile-time check to verify implements interface.
var _ SecretManager = (*Noop)(nil)

// Noop is a secret manager that does nothing and always returns an error.
type Noop struct{}

func NewNoop(ctx context.Context) (SecretManager, error) {
	return &Noop{}, nil
}

// GetSecretValue implements secrets.
func (n *Noop) GetSecretValue(_ context.Context, _ string) (string, error) {
	return "noop-secret", nil
}
