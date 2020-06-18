package signing

import (
	"context"
	"crypto"
	"fmt"
)

// Compile-time check to verify implements interface.
var _ KeyManager = (*Noop)(nil)

// Noop is a key manager that does nothing and always returns an error.
type Noop struct{}

func NewNoop(ctx context.Context) (KeyManager, error) {
	return &Noop{}, nil
}

func (n *Noop) NewSigner(ctx context.Context, keyID string) (crypto.Signer, error) {
	return nil, fmt.Errorf("noop cannot sign")
}
