// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package rest

import (
	"context"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/registry/config"
)

// Injectors from wire.go:

func InitializeApp(ctx context.Context) *App {
	configConfig := config.NewConfig(ctx)
	zapLogger := logger.NewLogger(ctx, configConfig)
	app := &App{
		Config: configConfig,
		Logger: zapLogger,
	}
	return app
}
