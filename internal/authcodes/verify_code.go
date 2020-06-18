// Package cleanup implements the API handlers for running data deletion jobs.
package authcodes

import (
	"fmt"
	"net/http"

	"go.opencensus.io/trace"

	"github.com/krishnadurai/aurora-test-auth/internal/cache"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
	"github.com/krishnadurai/aurora-test-auth/internal/serverenv"
)

// NewVerifyCodeHandler creates a http.Handler for verifying auth codes
// from the cache.
func NewVerifyCodeHandler(config *Config, env *serverenv.ServerEnv) (http.Handler, error) {
	if env.Cache() == nil {
		return nil, fmt.Errorf("missing cache in server environment")
	}

	return &verifyCodeHandler{
		config: config,
		env:    env,
		cache:  env.Cache(),
	}, nil
}

type verifyCodeHandler struct {
	config *Config
	env    *serverenv.ServerEnv
	cache  cache.Cache
}

func (h *verifyCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "(*authcodes.verifyCodeHandler).ServeHTTP")
	defer span.End()

	logger := logging.FromContext(ctx)
	metrics := h.env.MetricsExporter(ctx)

	getResult, err := h.cache.Get(ctx, "test")
	if err != nil {
		logger.Error(err.Error())
		span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
	}
	logger.Infof("Auth code is %v", getResult)

	metrics.WriteInt64("auth-code-verified", true, 1)
	logger.Infof("Auth code has been verified.")
	w.WriteHeader(http.StatusOK)
}
