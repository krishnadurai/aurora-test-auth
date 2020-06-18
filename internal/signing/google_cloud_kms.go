package signing

import (
	"context"
	"crypto"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/sethvargo/go-gcpkms/pkg/gcpkms"
)

// Compile-time check to verify implements interface.
var _ KeyManager = (*GoogleCloudKMS)(nil)

// GoogleCloudKMS implements the signing.KeyManager interface and can be used to sign
// export files.
type GoogleCloudKMS struct {
	client *kms.KeyManagementClient
}

func NewGoogleCloudKMS(ctx context.Context) (KeyManager, error) {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GoogleCloudKMS{client}, nil
}

func (kms *GoogleCloudKMS) NewSigner(ctx context.Context, keyID string) (crypto.Signer, error) {
	return gcpkms.NewSigner(ctx, kms.client, keyID)
}
