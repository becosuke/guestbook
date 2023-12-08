//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/becosuke/syncmap"
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	"github.com/becosuke/guestbook/api/internal/adapter/repository/generator"
	repository_syncmap "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/application/interactor"
	"github.com/becosuke/guestbook/api/internal/driver/grpcserver"
	"github.com/becosuke/guestbook/api/internal/driver/interceptor"
	pkgconfig "github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
)

type App struct {
	Config     *pkgconfig.Config
	Logger     *zap.Logger
	GrpcServer *grpc.Server
	Controller pb.GuestbookServiceServer
}

var controllerSet = wire.NewSet(
	controller.NewGuestbookServiceServer,
)

var usecaseSet = wire.NewSet(
	interactor.NewUsecase,
)

var repositorySet = wire.NewSet(
	generator.NewGenerator,
	repository_syncmap.NewQuerier,
	repository_syncmap.NewCommander,
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

		syncmap.NewSyncmap,

		wire.Struct(new(App), "*"),
	)
	return nil
}
