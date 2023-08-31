//go:build wireinject
// +build wireinject

package rest

import (
	"context"

	"github.com/google/wire"

	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/api/internal/registry/injection"
)

func InitializeApp(ctx context.Context) *App {
	wire.Build(
		pkgconfig.NewConfig,
		injection.ProvideLogger,

		wire.Struct(new(App), "*"),
	)
	return nil
}
