package config

import (
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
)

const (
	KeyString   = "key"
	ValueString = "value"
	Endpoint    = "/kvstore"
)

type Config struct {
	constConfig
	envConfig
}

func NewConfig() *Config {
	return &Config{
		constConfig: newConstConfig(),
		envConfig:   newEnvConfig(),
	}
}

type constConfig struct {
	KeyString   string
	ValueString string
	Endpoint    string
}

type envConfig struct {
	Environment Environment
	LogLevel    zapcore.Level
	GrpcPort    int
	HttpPort    int
}

func newConstConfig() constConfig {
	return constConfig{
		KeyString:   KeyString,
		ValueString: ValueString,
		Endpoint:    Endpoint,
	}
}

func newEnvConfig() envConfig {
	environmentString, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environmentString = "development"
	}
	environment := NewEnvironment(environmentString)
	if environment == EnvUnknown {
		environment = EnvDevelopment
	}

	logLevelString, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelString = "info"
	}
	logLevel := zapcore.InfoLevel
	_ = logLevel.Set(logLevelString)

	grpcPortString, ok := os.LookupEnv("GRPC_PORT")
	if !ok {
		grpcPortString = "50051"
	}
	grpcPort, err := strconv.Atoi(grpcPortString)
	if err != nil {
		grpcPort = 50051
	}

	httpPortString, ok := os.LookupEnv("HTTP_PORT")
	if !ok {
		httpPortString = "50080"
	}
	httpPort, err := strconv.Atoi(httpPortString)
	if err != nil {
		httpPort = 50080
	}

	return envConfig{
		Environment: environment,
		LogLevel:    logLevel,
		GrpcPort:    grpcPort,
		HttpPort:    httpPort,
	}
}
