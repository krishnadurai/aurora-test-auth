package signing

import (
	"context"
	"crypto"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/lstoll/awskms"
)

// Compile-time check to verify implements interface.
var _ KeyManager = (*AWSKMS)(nil)

// AWSKMS implements the signing.KeyManager interface and can be used to sign
// export files using AWS KMS.
type AWSKMS struct {
	svc *kms.KMS
}

func NewAWSKMS(ctx context.Context) (KeyManager, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	svc := kms.New(sess)

	return &AWSKMS{
		svc: svc,
	}, nil
}

func (s *AWSKMS) NewSigner(ctx context.Context, keyID string) (crypto.Signer, error) {
	return awskms.NewSigner(ctx, s.svc, keyID)
}
