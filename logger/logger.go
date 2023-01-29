package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level, serviceName, version, environment string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.DisableStacktrace = true
	config.Sampling = nil
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.InitialFields = map[string]interface{}{
		"service": serviceName,
		"version": version,
		"env":     environment,
	}
	return config.Build()
}

func NewFakeLogger() *zap.Logger {
	return zap.NewNop()
}
