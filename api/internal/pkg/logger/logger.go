package logger

import (
	"context"

	"go.uber.org/zap"

	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
)

func NewLogger(ctx context.Context, config *pkgconfig.Config) *zap.Logger {
	serviceName, ok := ctx.Value(pkgconfig.ServiceName{}).(string)
	if !ok {
		serviceName = ""
	}
	serviceVersion, ok := ctx.Value(pkgconfig.ServiceVersion{}).(string)
	if !ok {
		serviceVersion = ""
	}
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(config.LogLevel)
	loggerConfig.DisableStacktrace = true
	loggerConfig.Sampling = nil
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}
	loggerConfig.InitialFields = map[string]interface{}{
		"service": serviceName,
		"version": serviceVersion,
		"env":     config.Environment.String(),
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		panic("failed to provide logger")
	}
	return logger
}