//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	syncmap_repository "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/usecase"
	"github.com/becosuke/guestbook/api/internal/driver/grpcserver"
	"github.com/becosuke/guestbook/api/internal/driver/interceptor"
	syncmap_driver "github.com/becosuke/guestbook/api/internal/driver/syncmap"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/pbgo"
)

type App struct {
	Config     *pkgconfig.Config
	Logger     *zap.Logger
	GrpcServer *grpc.Server
	Controller pbgo.GuestbookServiceServer
}

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
		logger.NewLogger,
		interceptor.NewAuthFunc,
		grpcserver.NewGrpcServer,
		controllerSet,
		usecaseSet,
		repositorySet,
		driverSet,

		wire.Struct(new(App), "*"),
	)
	return nil
}
