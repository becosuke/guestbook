package main

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/becosuke/guestbook/api/internal/adapter/controller"
	repository_postgres "github.com/becosuke/guestbook/api/internal/adapter/repository/postgres"
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
	pool, err := repository_postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	querier := repository_postgres.NewQuerier(cfg, zapLogger, pool)
	commander := repository_postgres.NewCommander(cfg, zapLogger, pool)
	uc := usecase.NewUsecase(cfg, zapLogger, querier, commander)
	ctrl := controller.NewGuestbookServiceServer(cfg, zapLogger, uc)
	return &App{
		Config:     cfg,
		Logger:     zapLogger,
		GrpcServer: server,
		Controller: ctrl,
	}
}
