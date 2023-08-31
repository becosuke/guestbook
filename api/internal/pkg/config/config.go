package config

import (
	"context"
	"os"
	"strconv"

	"go.uber.org/zap/zapcore"
)

type Config struct {
	constConfig
	envConfig
}

type constConfig struct{}

type envConfig struct {
	Environment Environment
	LogLevel    zapcore.Level
	GrpcHost    string
	GrpcPort    int
	RestHost    string
	RestPort    int
}

func NewConfig(ctx context.Context) *Config {
	return &Config{
		constConfig: newConstConfig(ctx),
		envConfig:   newEnvConfig(ctx),
	}
}

type (
	ServiceName    struct{}
	ServiceVersion struct{}
)

func newConstConfig(ctx context.Context) constConfig {
	return constConfig{}
}

func newEnvConfig(ctx context.Context) envConfig {
	environmentString, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environmentString = "development"
	}
	environment, err := NewEnvironment(environmentString)
	if err != nil {
		environment = EnvDevelopment
	}

	logLevelString, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelString = "info"
	}
	logLevel, err := zapcore.ParseLevel(logLevelString)
	if err != nil {
		logLevel = zapcore.InfoLevel
	}

	grpcHost, ok := os.LookupEnv("GRPC_HOST")
	if !ok {
		grpcHost = ""
	}

	grpcPortString, ok := os.LookupEnv("GRPC_PORT")
	if !ok {
		grpcPortString = "50051"
	}
	grpcPort, err := strconv.Atoi(grpcPortString)
	if err != nil {
		grpcPort = 50051
	}

	restHost, ok := os.LookupEnv("REST_HOST")
	if !ok {
		restHost = ""
	}

	restPortString, ok := os.LookupEnv("REST_PORT")
	if !ok {
		restPortString = "50080"
	}
	restPort, err := strconv.Atoi(restPortString)
	if err != nil {
		restPort = 50080
	}

	return envConfig{
		Environment: environment,
		LogLevel:    logLevel,
		GrpcHost:    grpcHost,
		GrpcPort:    grpcPort,
		RestHost:    restHost,
		RestPort:    restPort,
	}
}
