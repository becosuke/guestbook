package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
}

func InitializeApp(ctx context.Context) *App {
	cfg := config.NewConfig(ctx)
	zapLogger := logger.NewLogger(ctx, cfg)
	return &App{
		Config: cfg,
		Logger: zapLogger,
	}
}
