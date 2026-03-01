package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"

	"github.com/becosuke/guestbook/api/internal/domain/entity"
)

type envConfig struct {
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	GrpcHost    string `envconfig:"GRPC_HOST" default:"127.0.0.1"`
	GrpcPort    int    `envconfig:"GRPC_PORT" default:"50051"`
	RestHost    string `envconfig:"REST_HOST" default:"127.0.0.1"`
	RestPort    int    `envconfig:"REST_PORT" default:"50080"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func NewConfig() *entity.Config {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		panic("failed to process env config: " + err.Error())
	}

	environment, err := entity.NewEnvironment(env.Environment)
	if err != nil {
		environment = entity.EnvDevelopment
	}

	logLevel, err := zapcore.ParseLevel(env.LogLevel)
	if err != nil {
		logLevel = zapcore.DebugLevel
	}

	return &entity.Config{
		Environment: environment,
		LogLevel:    logLevel,
		GrpcHost:    env.GrpcHost,
		GrpcPort:    env.GrpcPort,
		RestHost:    env.RestHost,
		RestPort:    env.RestPort,
		DatabaseURL: env.DatabaseURL,
	}
}
