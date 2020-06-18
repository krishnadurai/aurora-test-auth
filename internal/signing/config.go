package signing

// KeyManagerType defines a specific key manager.
type KeyManagerType string

const (
	KeyManagerTypeAWSKMS         KeyManagerType = "AWS_KMS"
	KeyManagerTypeAzureKeyVault  KeyManagerType = "AZURE_KEY_VAULT"
	KeyManagerTypeGoogleCloudKMS KeyManagerType = "GOOGLE_CLOUD_KMS"
	KeyManagerTypeHashiCorpVault KeyManagerType = "HASHICORP_VAULT"
	KeyManagerTypeNoop           KeyManagerType = "NOOP"
)

// Config defines configuration.
type Config struct {
	KeyManagerType KeyManagerType `env:"KEY_MANAGER,default=GOOGLE_CLOUD_KMS"`
}
