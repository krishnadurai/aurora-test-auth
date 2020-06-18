// Package metrics contains utilities for exporting metrics.
package metrics

import (
	"context"

	"github.com/krishnadurai/aurora-test-auth/internal/logging"
	"go.uber.org/zap"
)

const logString = "!METRIC! Type = %v cumulative = %v value = %v"

// ExporterFromContext defines a function to create a new exporter based on the current context.
type ExporterFromContext func(context.Context) Exporter

// Exporter defines a generic metric exporter interface used in this application.
type Exporter interface {
	WriteBool(name string, value bool)
	WriteInt(name string, cumulative bool, value int)
	WriteInt64(name string, cumulative bool, value int64)
	WriteIntDistribution(name string, cumulative bool, values []int)
	WriteFloat64(name string, cumulative bool, value float64)
	WriteFloat64Distribution(name string, cumulative bool, values []float64)
}

type exporterImpl struct {
	logger *zap.SugaredLogger
}

// NewLogsBasedFromContext creates a logs based metrics exporter.
func NewLogsBasedFromContext(ctx context.Context) Exporter {
	return NewLogsBasedExporter(logging.FromContext(ctx))
}

// NewLogsBasedExporter creates a new logs based importer from a logger.
func NewLogsBasedExporter(log *zap.SugaredLogger) Exporter {
	return &exporterImpl{
		logger: log,
	}
}

func (e *exporterImpl) WriteBool(name string, value bool) {
	e.logger.Infof(logString, name, false, value)
}

func (e *exporterImpl) WriteInt(name string, cumulative bool, value int) {
	e.logger.Infof(logString, name, cumulative, value)
}

func (e *exporterImpl) WriteInt64(name string, cumulative bool, value int64) {
	e.logger.Infof(logString, name, cumulative, value)
}

func (e *exporterImpl) WriteIntDistribution(name string, cumulative bool, values []int) {
	e.logger.Infof(logString, name, cumulative, values)
}

func (e *exporterImpl) WriteFloat64(name string, cumulative bool, value float64) {
	e.logger.Infof(logString, name, cumulative, value)
}

func (e *exporterImpl) WriteFloat64Distribution(name string, cumulative bool, values []float64) {
	e.logger.Infof(logString, name, cumulative, values)
}
