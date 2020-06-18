package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/krishnadurai/aurora-test-auth/internal/authcodes"
	"github.com/krishnadurai/aurora-test-auth/internal/interrupt"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
	"github.com/krishnadurai/aurora-test-auth/internal/server"
	"github.com/krishnadurai/aurora-test-auth/internal/setup"
)

func main() {
	ctx, done := interrupt.Context()
	defer done()

	if err := realMain(ctx); err != nil {
		logger := logging.FromContext(ctx)
		logger.Fatal(err)
	}
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config authcodes.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	handler, err := authcodes.NewVerifyCodeHandler(&config, env)
	if err != nil {
		return fmt.Errorf("cleanup.NewExposureHandler: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv, err := server.New(config.Port)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Infof("listening on :%s", config.Port)

	return srv.ServeHTTPHandler(ctx, mux)
}
