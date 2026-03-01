package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(serviceName, serviceVersion, environment string, logLevel zapcore.Level) *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)
	loggerConfig.DisableStacktrace = true
	loggerConfig.Sampling = nil
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.ErrorOutputPaths = []string{"stderr"}
	loggerConfig.InitialFields = map[string]interface{}{
		"service": serviceName,
		"version": serviceVersion,
		"env":     environment,
	}
	logger, err := loggerConfig.Build()
	if err != nil {
		panic("failed to build logger")
	}
	return logger
}
