package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	infraconfig "github.com/becosuke/guestbook/api/internal/adapter/infrastructure/config"
	"github.com/becosuke/guestbook/api/internal/domain"
	"github.com/becosuke/guestbook/api/internal/pkg/logger"
	"github.com/becosuke/guestbook/api/internal/pkg/pb"
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
	Config *domain.Config
	Logger *zap.Logger
}

func InitializeApp(ctx context.Context) *App {
	cfg := infraconfig.NewConfig()
	zapLogger := logger.NewLogger(serviceName, version, cfg.Environment.String(), cfg.LogLevel)
	return &App{
		Config: cfg,
		Logger: zapLogger,
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

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(interrupt)

	err := pb.RegisterGuestbookServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort), opts)
	if err != nil {
		zapLogger.Error("rest server: failed to register handler", zap.Error(err))
		return exitError
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.RestHost, cfg.RestPort),
		Handler: mux,
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			zapLogger.Error("rest server failed to listen", zap.Error(err))
			return err
		}
		return nil
	})
	zapLogger.Info("rest server: listening", zap.String("host", cfg.RestHost), zap.Int("port", cfg.RestPort))

	select {
	case <-interrupt:
		zapLogger.Info("received shutdown signal")
	case <-ctx.Done():
		zapLogger.Info("context canceled")
	}
	cancel()

	zapLogger.Info("rest server: going gracefully shutdown")
	if err := httpServer.Shutdown(ctx); err != context.Canceled {
		zapLogger.Error("received error on gracefully shutdown", zap.Error(err))
	}
	zapLogger.Info("rest server: completed gracefully shutdown")

	if err := eg.Wait(); err != nil {
		zapLogger.Error("received error", zap.Error(err))
		return exitError
	}

	return exitOK
}
