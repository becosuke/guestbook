package main

import (
	"context"

	"github.com/becosuke/syncmap"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	repository_syncmap "github.com/becosuke/guestbook/api/internal/adapter/repository/syncmap"
	"github.com/becosuke/guestbook/api/internal/driver/grpcserver"
	"github.com/becosuke/guestbook/api/internal/driver/interceptor"
	"github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
	"github.com/becosuke/guestbook/api/internal/usecase"
)

type App struct {
	Config     *config.Config
	Logger     *zap.Logger
	GrpcServer *grpc.Server
	Controller pb.GuestbookServiceServer
}

func InitializeApp(ctx context.Context) *App {
	cfg := config.NewConfig(ctx)
	zapLogger := logger.NewLogger(ctx, cfg)
	authFunc := interceptor.NewAuthFunc(ctx)
	server := grpcserver.NewGrpcServer(ctx, zapLogger, authFunc)
	store := syncmap.NewSyncmap()
	querier := repository_syncmap.NewQuerier(cfg, zapLogger, store)
	commander := repository_syncmap.NewCommander(cfg, zapLogger, store)
	uc := usecase.NewUsecase(cfg, zapLogger, querier, commander)
	ctrl := controller.NewGuestbookServiceServer(cfg, zapLogger, uc)
	return &App{
		Config:     cfg,
		Logger:     zapLogger,
		GrpcServer: server,
		Controller: ctrl,
	}
}
