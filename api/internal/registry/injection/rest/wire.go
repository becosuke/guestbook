//go:build wireinject
// +build wireinject

package rest

import (
	"context"

	"github.com/google/wire"

	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
)

func InitializeApp(ctx context.Context) *App {
	wire.Build(
		pkgconfig.NewConfig,
		logger.NewLogger,

		wire.Struct(new(App), "*"),
	)
	return nil
}
