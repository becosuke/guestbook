//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"

	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
)

type App struct {
	Config *pkgconfig.Config
	Logger *zap.Logger
}

func InitializeApp(ctx context.Context) *App {
	wire.Build(
		pkgconfig.NewConfig,
		logger.NewLogger,

		wire.Struct(new(App), "*"),
	)
	return nil
}
