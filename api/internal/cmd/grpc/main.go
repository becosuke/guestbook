package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/becosuke/guestbook/api/internal/adapter/presentation"
	repository_postgres "github.com/becosuke/guestbook/api/internal/adapter/repository/postgres"
	"github.com/becosuke/guestbook/api/internal/pkg/config"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
	"github.com/becosuke/guestbook/api/internal/usecase"
)

const (
	exitOK int = iota
	exitError
)

var (
	serviceName string
	version     string
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
	authFunc := func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
			grpc_auth.StreamServerInterceptor(authFunc),
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
			grpc_auth.UnaryServerInterceptor(authFunc),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pool, err := repository_postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	querier := repository_postgres.NewQuerier(cfg, zapLogger, pool)
	commander := repository_postgres.NewCommander(cfg, zapLogger, pool)
	uc := usecase.NewUsecase(cfg, zapLogger, querier, commander)
	ctrl := presentation.NewGuestbookServiceServer(cfg, zapLogger, uc)
	return &App{
		Config:     cfg,
		Logger:     zapLogger,
		GrpcServer: server,
		Controller: ctrl,
	}
}

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, config.ServiceName{}, serviceName)
	ctx = context.WithValue(ctx, config.ServiceVersion{}, version)

	app := InitializeApp(ctx)
	cfg := app.Config
	zapLogger := app.Logger
	defer func() {
		_ = zapLogger.Sync()
	}()

	grpcServer := app.GrpcServer
	ctrl := app.Controller
	pb.RegisterGuestbookServiceServer(grpcServer, ctrl)
	reflection.Register(grpcServer)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort))
	if err != nil {
		zapLogger.Error("grpc server: failed to listen", zap.Error(err))
		return exitError
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return grpcServer.Serve(listener)
	})
	zapLogger.Info("grpc server: serving", zap.String("host", cfg.GrpcHost), zap.Int("port", cfg.GrpcPort))

	select {
	case <-interrupt:
		zapLogger.Info("received shutdown signal")
	case <-ctx.Done():
		zapLogger.Info("context canceled")
	}
	cancel()

	zapLogger.Info("grpc server: going gracefully shutdown")
	grpcServer.GracefulStop()
	zapLogger.Info("grpc server: completed gracefully shutdown")

	if err := eg.Wait(); err != nil {
		zapLogger.Error("received error", zap.Error(err))
		return exitError
	}

	return exitOK
}
