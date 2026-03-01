package domain

import (
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Environment Environment
	LogLevel    zapcore.Level
	GrpcHost    string
	GrpcPort    int
	RestHost    string
	RestPort    int
	DatabaseURL string
}
