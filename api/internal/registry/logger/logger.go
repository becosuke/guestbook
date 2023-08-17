package logger

import (
	"log"

	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/registry/config"
)

func NewLogger(cfg *config.Config) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(cfg.LogLevel)
	zapConfig.DisableStacktrace = true
	zapConfig.Sampling = nil
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	zapConfig.InitialFields = map[string]interface{}{
		"service": cfg.serviceName,
		"version": i.version,
		"env":     i.InjectConfig().Environment.String(),
	}
	l, err := zapConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
}
