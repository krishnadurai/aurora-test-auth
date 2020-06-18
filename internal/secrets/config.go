package secrets

import (
	"time"
)

// SecretManagerType represents a type of secret manager.
type SecretManagerType string

const (
	SecretManagerTypeAWSSecretsManager    SecretManagerType = "AWS_SECRETS_MANAGER"
	SecretManagerTypeAzureKeyVault        SecretManagerType = "AZURE_KEY_VAULT"
	SecretManagerTypeGoogleSecretManager  SecretManagerType = "GOOGLE_SECRET_MANAGER"
	SecretManagerTypeGoogleHashiCorpVault SecretManagerType = "HASHICORP_VAULT"
	SecretManagerTypeNoop                 SecretManagerType = "NOOP"
)

// Config represents the config for a secret manager.
type Config struct {
	SecretManagerType SecretManagerType `env:"SECRET_MANAGER, default=GOOGLE_SECRET_MANAGER"`
	SecretsDir        string            `env:"SECRETS_DIR, default=/var/run/secrets"`
	SecretCacheTTL    time.Duration     `env:"SECRET_CACHE_TTL, default=5m"`
}
