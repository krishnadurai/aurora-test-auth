// Package signing defines the interface to and implementation of signing
package signing

import (
	"context"
	"crypto"
	"fmt"
)

// KeyManager defines the interface for working with a KMS system that
// is able to sign bytes using PKI.
// KeyManager implementations must be able to return a crypto.Signer.
type KeyManager interface {
	NewSigner(ctx context.Context, keyID string) (crypto.Signer, error)
}

// KeyManagerFor returns the appropriate key manager for the given type.
func KeyManagerFor(ctx context.Context, typ KeyManagerType) (KeyManager, error) {
	switch typ {
	case KeyManagerTypeAWSKMS:
		return NewAWSKMS(ctx)
	case KeyManagerTypeAzureKeyVault:
		return NewAzureKeyVault(ctx)
	case KeyManagerTypeGoogleCloudKMS:
		return NewGoogleCloudKMS(ctx)
	case KeyManagerTypeHashiCorpVault:
		return NewHashiCorpVault(ctx)
	case KeyManagerTypeNoop:
		return NewNoop(ctx)
	}

	return nil, fmt.Errorf("unknown key manager type: %v", typ)
}
