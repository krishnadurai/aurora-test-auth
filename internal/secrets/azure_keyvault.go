package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/krishnadurai/aurora-test-auth/internal/azurekeyvault"
)

// Compile-time check to verify implements interface.
var _ SecretManager = (*AzureKeyVault)(nil)

// AzureKeyVault implements SecretManager.
type AzureKeyVault struct {
	client *keyvault.BaseClient
}

// NewAzureKeyVault creates a new KeyVault that can interact fetch secrets.
func NewAzureKeyVault(ctx context.Context) (SecretManager, error) {
	authorizer, err := azurekeyvault.GetKeyVaultAuthorizer()
	if err != nil {
		return nil, fmt.Errorf("secrets.NewAzureKeyVault: auth: %w", err)
	}

	client := keyvault.New()
	client.Authorizer = authorizer

	sm := &AzureKeyVault{
		client: &client,
	}

	return sm, nil
}

// GetSecretValue implements the SecretManager interface. Secrets are specified
// in the format:
//
//     AZURE_KEY_VAULT_NAME/SECRET_NAME/SECRET_VERSION
//
// For example:
//
//     my-company-vault/api-key/1
//
// If the secret version is omitted, the latest version is used.
func (kv *AzureKeyVault) GetSecretValue(ctx context.Context, name string) (string, error) {
	// Extract vault, secret, and version.
	var vaultName, secretName, version string
	parts := strings.SplitN(name, "/", 3)
	switch len(parts) {
	case 0, 1:
		return "", fmt.Errorf("%v is not a valid secret ref", name)
	case 2:
		vaultName, secretName, version = parts[0], parts[1], ""
	case 3:
		vaultName, secretName, version = parts[0], parts[1], parts[2]
	}

	// Lookup in KeyVault
	vaultURL := fmt.Sprintf("https://%s.vault.azure.net", vaultName)
	result, err := kv.client.GetSecret(ctx, vaultURL, secretName, version)
	if err != nil {
		return "", fmt.Errorf("failed to access secret %v: %w", name, err)
	}
	if result.Value == nil {
		return "", fmt.Errorf("found secret %v, but value was nil", name)
	}
	return *result.Value, nil
}
