// Package server provides an opinionated http server.
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/krishnadurai/aurora-test-auth/internal/logging"
)

// ServeMetricsIfPrometheus serves the opencensus metrics at /metrics when OBSERVABILITY_EXPORTER set to "prometheus"
func ServeMetricsIfPrometheus(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	exporter := os.Getenv("OBSERVABILITY_EXPORTER")
	metricsPort := os.Getenv("METRICS_PORT")
	if strings.EqualFold(exporter, "prometheus") {
		if metricsPort == "" {
			return fmt.Errorf("OBSERVABILITY_EXPORTER set to 'prometheus' but no METRICS_PORT set")
		}

		exporter, err := prometheus.NewExporter(prometheus.Options{})
		if err != nil {
			return fmt.Errorf("failed to create prometheus exporter: %v", err)
		}

		go func() {
			mux := http.NewServeMux()
			mux.Handle("/metrics", exporter)

			logger.Debugf("Metrics endpoint listening on :%s", metricsPort)
			if err := http.ListenAndServe(":"+metricsPort, mux); err != nil {
				logger.Debugf("error while serving metrics endpoint: %w", err)
			}
		}()
	}
	return nil
}
