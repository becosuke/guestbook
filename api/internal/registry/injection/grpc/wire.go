//go:build wireinject
// +build wireinject

package grpc

import (
	"context"

	"github.com/google/wire"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	syncmap_repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/driver/grpcserver"
	syncmap_driver "github.com/becosuke/guestbook/api/internal/driver/syncmap"
	pkgconfig "github.com/becosuke/guestbook/api/internal/registry/config"
	"github.com/becosuke/guestbook/api/internal/registry/injection"
)

var controllerSet = wire.NewSet(
	controller.NewGuestbookServiceServer,
	controller.NewBoundary,
)

var usecaseSet = wire.NewSet(
	usecase.NewUsecase,
)

var repositorySet = wire.NewSet(
	syncmap_repository.NewGenerator,
	syncmap_repository.NewQuerier,
	syncmap_repository.NewCommander,
	syncmap_repository.NewBoundary,
)

var driverSet = wire.NewSet(
	syncmap_driver.NewSyncmap,
)

func InitializeApp(ctx context.Context) *App {
	wire.Build(
		pkgconfig.NewConfig,
		injection.ProvideLogger,
		injection.ProvideAuthFunc,
		grpcserver.NewGrpcServer,
		controllerSet,
		usecaseSet,
		repositorySet,
		driverSet,

		wire.Struct(new(App), "*"),
	)
	return nil
}
