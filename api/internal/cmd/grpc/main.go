package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"buf.build/go/protovalidate"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	infraconfig "github.com/becosuke/guestbook/api/internal/adapter/infrastructure/config"
	"github.com/becosuke/guestbook/api/internal/adapter/presentation/interceptor"
	"github.com/becosuke/guestbook/api/internal/adapter/presentation"
	"github.com/becosuke/guestbook/api/internal/adapter/repository"
	"github.com/becosuke/guestbook/api/internal/domain"
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
	Config     *domain.Config
	Logger     *zap.Logger
	GrpcServer *grpc.Server
	Controller pb.GuestbookServiceServer
}

func InitializeApp(ctx context.Context) *App {
	cfg := infraconfig.NewConfig()
	zapLogger := logger.NewLogger(serviceName, version, cfg.Environment.String(), cfg.LogLevel)
	authFunc := func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
	validator, err := protovalidate.New()
	if err != nil {
		zapLogger.Fatal("failed to create protovalidate validator", zap.Error(err))
	}
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zapLogger),
			grpc_auth.StreamServerInterceptor(authFunc),
			interceptor.ProtovalidateStreamServerInterceptor(validator),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zapLogger),
			grpc_auth.UnaryServerInterceptor(authFunc),
			interceptor.ProtovalidateUnaryServerInterceptor(validator),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		zapLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	repos := repository.NewRepositories(ctx, cfg, zapLogger, pool)
	uc := usecase.NewUsecase(cfg, zapLogger, repos)
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
