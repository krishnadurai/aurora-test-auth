package setup

import (
	"context"
	"fmt"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
	"github.com/krishnadurai/aurora-test-auth/internal/metrics"
	"github.com/krishnadurai/aurora-test-auth/internal/secrets"
	"github.com/krishnadurai/aurora-test-auth/internal/serverenv"
	"github.com/krishnadurai/aurora-test-auth/internal/signing"
	"github.com/sethvargo/go-envconfig/pkg/envconfig"
)

// CacheConfigProvider ensures that the environment config can provide a cache config.
// All binaries in this application connect to the cache via the same method.
type CacheConfigProvider interface {
	CacheConfig() *cache.Config
}

// KeyManagerConfigProvider is a marker interface indicating the key manager
// should be installed.
type KeyManagerConfigProvider interface {
	KeyManagerConfig() *signing.Config
}

// SecretManagerConfigProvider signals that the config knows how to configure a
// secret manager.
type SecretManagerConfigProvider interface {
	SecretManagerConfig() *secrets.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	// Build a list of mutators. This list will grow as we initialize more of the
	// configuration, such as the secret manager.
	var mutatorFuncs []envconfig.MutatorFunc

	// Build a list of options to pass to the server env.
	var serverEnvOpts []serverenv.Option

	// TODO: support customizable metrics
	serverEnvOpts = append(serverEnvOpts, serverenv.WithMetricsExporter(metrics.NewLogsBasedFromContext))

	// Load the secret manager - this needs to be loaded first because other
	// processors may require access to secrets.
	var sm secrets.SecretManager
	if provider, ok := config.(SecretManagerConfigProvider); ok {
		logger.Info("configuring secret manager")

		// The environment configuration defines which secret manager to use, so we
		// need to process the envconfig in here. Once we know which secret manager
		// to use, we can load the secret manager and then re-process (later) any
		// secret:// references.
		//
		// NOTE: it is not possible to specify which secret manager to use via a
		// secret:// reference. This configuration option must always be the
		// plaintext string.
		smConfig := provider.SecretManagerConfig()
		if err := envconfig.ProcessWith(ctx, smConfig, l, mutatorFuncs...); err != nil {
			return nil, fmt.Errorf("unable to process secret manager env: %w", err)
		}

		var err error
		sm, err = secrets.SecretManagerFor(ctx, smConfig.SecretManagerType)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to secret manager: %w", err)
		}

		// Enable caching, if a TTL was provided.
		if ttl := smConfig.SecretCacheTTL; ttl > 0 {
			sm, err = secrets.WrapCacher(ctx, sm, ttl)
			if err != nil {
				return nil, fmt.Errorf("unable to create secret manager cache: %w", err)
			}
		}

		// Update the mutators to process secrets.
		mutatorFuncs = append(mutatorFuncs, secrets.Resolver(sm, smConfig))

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithSecretManager(sm))

		logger.Infow("secret manager", "config", smConfig)
	}

	// Load the key manager.
	var km signing.KeyManager
	if provider, ok := config.(KeyManagerConfigProvider); ok {
		logger.Info("configuring key manager")

		kmConfig := provider.KeyManagerConfig()
		if err := envconfig.ProcessWith(ctx, kmConfig, l, mutatorFuncs...); err != nil {
			return nil, fmt.Errorf("unable to process key manager env: %w", err)
		}

		var err error
		km, err = signing.KeyManagerFor(ctx, kmConfig.KeyManagerType)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to key manager: %w", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithKeyManager(km))

		logger.Infow("key manager", "config", kmConfig)
	}

	// Process first round of environment variables.
	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}
	logger.Infow("provided", "config", config)

	// Setup the database connection.
	if provider, ok := config.(CacheConfigProvider); ok {
		logger.Info("configuring cache")

		cacheConfig := provider.CacheConfig()
		cache, err := cache.NewRedisCache(ctx, cacheConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to cache: %v", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithCache(cache))

		logger.Infow("cache", "config", cacheConfig)

	}

	return serverenv.New(ctx, serverEnvOpts...), nil
}
